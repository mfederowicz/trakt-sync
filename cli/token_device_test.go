// Package cli for basic cli functions
package cli

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/trakt/trakttest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const testSettings = `{"user":{"username":"sean"},"account":{"timezone":"Europe/Warsaw"}}`

// credentialsConfig returns a config whose token and settings files live in a fresh temp directory.
func credentialsConfig(t *testing.T) *cfg.Config {
	t.Helper()
	dir := t.TempDir()
	config := cfg.DefaultConfig()
	config.ClientID = "client-id"
	config.ClientSecret = "client-secret"
	config.TokenPath = filepath.Join(dir, "token.json")
	config.SettingsPath = filepath.Join(dir, "user_settings.json")
	return config
}

// tokenJSON builds a token file body that was created now and lives for expiresIn seconds.
func tokenJSON(accessToken string, expiresIn int) string {
	created := strconv.FormatInt(time.Now().Unix(), consts.BaseInt)
	return `{"access_token":"` + accessToken + `","token_type":"bearer","refresh_token":"refresh","scope":"public","expires_in":` +
		strconv.Itoa(expiresIn) + `,"created_at":` + created + `}`
}

const expiredToken = `{"access_token":"old-token","refresh_token":"refresh","expires_in":1,"created_at":1}`

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	assert.NoError(t, os.WriteFile(path, []byte(content), consts.X600))
}

// decodeBody decodes the JSON request body into v.
func decodeBody(t *testing.T, r *http.Request, v any) {
	t.Helper()
	assert.NoError(t, json.NewDecoder(r.Body).Decode(v))
}

func TestReadCredentialFiles(t *testing.T) {
	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token.json")
	settingsPath := filepath.Join(dir, "user_settings.json")
	broken := filepath.Join(dir, "broken.json")
	missing := filepath.Join(dir, "missing.json")
	writeFile(t, tokenPath, `{"access_token":"access","refresh_token":"refresh","expires_in":7200,"created_at":1790863634}`)
	writeFile(t, settingsPath, testSettings)
	writeFile(t, broken, "not json")

	token, err := ReadTokenFromFile(tokenPath)
	if assert.NoError(t, err) {
		assert.Equal(t, str.Token{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: 7200, CreatedAt: 1790863634}, *token)
	}
	settings, err := ReadUserSettingsFromFile(settingsPath)
	if assert.NoError(t, err) {
		assert.Equal(t, "sean", *settings.User.Username)
		assert.Equal(t, "Europe/Warsaw", *settings.Account.Timezone)
	}

	for _, path := range []string{broken, missing} {
		token, err := ReadTokenFromFile(path)
		assert.Nil(t, token, path)
		assert.Error(t, err, path)

		settings, err := ReadUserSettingsFromFile(path)
		assert.Nil(t, settings, path)
		assert.Error(t, err, path)
	}
}

func TestValidAccessToken(t *testing.T) {
	t.Run("valid token needs no request", func(t *testing.T) {
		config := credentialsConfig(t)
		writeFile(t, config.TokenPath, tokenJSON("access", 7200))
		s := trakttest.Setup()
		defer s.Teardown()
		s.Mux.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		})

		assert.True(t, ValidAccessToken(config, s.Client, &str.Options{}))
	})

	t.Run("no token file", func(t *testing.T) {
		s := trakttest.Setup()
		defer s.Teardown()
		assert.False(t, ValidAccessToken(credentialsConfig(t), s.Client, &str.Options{}))
	})

	t.Run("token file without an access token", func(t *testing.T) {
		config := credentialsConfig(t)
		writeFile(t, config.TokenPath, `{"created_at":1}`)
		s := trakttest.Setup()
		defer s.Teardown()
		assert.False(t, ValidAccessToken(config, s.Client, &str.Options{}))
	})

	t.Run("expired token that cannot be refreshed", func(t *testing.T) {
		config := credentialsConfig(t)
		writeFile(t, config.TokenPath, expiredToken)
		s := trakttest.Setup()
		defer s.Teardown()
		s.Mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
			test.AssertMethod(t, r, http.MethodPost)
			w.WriteHeader(http.StatusUnauthorized)
		})
		s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})

		assert.False(t, ValidAccessToken(config, s.Client, &str.Options{}))
		token, err := ReadTokenFromFile(config.TokenPath)
		if assert.NoError(t, err) {
			assert.Equal(t, "old-token", token.AccessToken, "the token file stays as it was")
		}
	})
}

func TestRefreshToken(t *testing.T) {
	t.Run("sends the refresh token and stores the new one", func(t *testing.T) {
		config := credentialsConfig(t)
		config.RedirectURI = "urn:ietf:wg:oauth:2.0:oob"
		writeFile(t, config.TokenPath, expiredToken)
		s := trakttest.Setup()
		defer s.Teardown()
		var sent str.CurrentDeviceToken
		s.Mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
			test.AssertMethod(t, r, http.MethodPost)
			decodeBody(t, r, &sent)
			test.SafeFprint(w, tokenJSON("new-token", 7776000))
		})

		assert.True(t, refreshToken(config, s.Client, &str.Options{}))
		assert.Equal(t, "refresh", *sent.RefreshToken)
		assert.Equal(t, "client-id", *sent.ClientID)
		assert.Equal(t, "client-secret", *sent.ClientSecret)
		assert.Equal(t, "urn:ietf:wg:oauth:2.0:oob", *sent.RedirectURI)
		assert.Equal(t, "refresh_token", *sent.GrantType)

		token, err := ReadTokenFromFile(config.TokenPath)
		if assert.NoError(t, err) {
			assert.Equal(t, "new-token", token.AccessToken)
			assert.False(t, token.Expired())
		}
	})

	t.Run("no token file", func(t *testing.T) {
		s := trakttest.Setup()
		defer s.Teardown()
		assert.False(t, refreshToken(credentialsConfig(t), s.Client, &str.Options{}))
	})

	t.Run("token file cannot be written", func(t *testing.T) {
		config := credentialsConfig(t)
		writeFile(t, config.TokenPath, expiredToken)
		assert.NoError(t, os.Chmod(config.TokenPath, 0o400))
		if file, err := os.OpenFile(config.TokenPath, os.O_WRONLY, 0); err == nil {
			_ = file.Close()
			t.Skip("read-only files are writable for this user")
		}
		s := trakttest.Setup()
		defer s.Teardown()
		s.Mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, _ *http.Request) {
			test.SafeFprint(w, tokenJSON("new-token", 7776000))
		})

		assert.False(t, refreshToken(config, s.Client, &str.Options{}))
	})
}

func TestRefreshUserSettings(t *testing.T) {
	t.Run("stores the settings", func(t *testing.T) {
		config := credentialsConfig(t)
		s := trakttest.Setup()
		defer s.Teardown()
		s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, r *http.Request) {
			test.AssertMethod(t, r, http.MethodGet)
			test.SafeFprint(w, testSettings)
		})

		assert.True(t, RefreshUserSettings(config, s.Client, &str.Options{}))
		settings, err := ReadUserSettingsFromFile(config.SettingsPath)
		if assert.NoError(t, err) {
			assert.Equal(t, "sean", *settings.User.Username)
		}
		info, err := os.Stat(config.SettingsPath)
		if assert.NoError(t, err) {
			assert.Equal(t, os.FileMode(consts.X600), info.Mode().Perm())
		}
	})

	t.Run("request fails", func(t *testing.T) {
		config := credentialsConfig(t)
		s := trakttest.Setup()
		defer s.Teardown()
		s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})

		assert.False(t, RefreshUserSettings(config, s.Client, &str.Options{}))
		_, err := os.Stat(config.SettingsPath)
		assert.True(t, os.IsNotExist(err), "no settings file is written")
	})

	t.Run("settings file cannot be written", func(t *testing.T) {
		config := credentialsConfig(t)
		config.SettingsPath = filepath.Join(filepath.Dir(config.SettingsPath), "no-such-dir", "user_settings.json")
		s := trakttest.Setup()
		defer s.Teardown()
		s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, _ *http.Request) {
			test.SafeFprint(w, testSettings)
		})

		assert.False(t, RefreshUserSettings(config, s.Client, &str.Options{}))
	})
}

// TestHandleTokenFetchesMissingUserSettings checks a valid token with settings that hold no user triggers a settings refresh.
func TestHandleTokenFetchesMissingUserSettings(t *testing.T) {
	config := credentialsConfig(t)
	writeFile(t, config.TokenPath, tokenJSON("file-token", 7776000))
	writeFile(t, config.SettingsPath, `{"account":{"timezone":"UTC"}}`)
	s := trakttest.Setup()
	defer s.Teardown()
	settingsAuth := ""
	s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, r *http.Request) {
		settingsAuth = r.Header.Get("Authorization")
		test.SafeFprint(w, testSettings)
	})

	HandleToken(afero.NewOsFs(), config, s.Client.WithClientID(config.ClientID), str.Options{})

	assert.Equal(t, "Bearer file-token", settingsAuth)
	settings, err := ReadUserSettingsFromFile(config.SettingsPath)
	if assert.NoError(t, err) && assert.NotNil(t, settings.User) {
		assert.Equal(t, "sean", *settings.User.Username)
	}
}

// A device login that fails is reported, the error used to be dropped without a word.
func TestHandleTokenReportsFailedDeviceLogin(t *testing.T) {
	openBrowser = func(string) error { return nil }
	t.Cleanup(func() { openBrowser = OpenBrowser })

	config := credentialsConfig(t)
	writeFile(t, config.TokenPath, `{}`)
	writeFile(t, config.SettingsPath, `{}`)
	s := trakttest.Setup()
	defer s.Teardown()
	s.Mux.HandleFunc("/oauth/device/code", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	stderr := os.Stderr
	r, w, err := os.Pipe()
	if !assert.NoError(t, err) {
		return
	}
	os.Stderr = w
	HandleToken(afero.NewOsFs(), config, s.Client.WithClientID(config.ClientID), str.Options{})
	os.Stderr = stderr
	assert.NoError(t, w.Close())

	out, err := io.ReadAll(r)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "generate new device code:")
}

func TestFetchNewDeviceCodeForClient(t *testing.T) {
	config := credentialsConfig(t)

	t.Run("new code", func(t *testing.T) {
		s := trakttest.Setup()
		defer s.Teardown()
		var sent str.NewDeviceCode
		s.Mux.HandleFunc("/oauth/device/code", func(w http.ResponseWriter, r *http.Request) {
			test.AssertMethod(t, r, http.MethodPost)
			decodeBody(t, r, &sent)
			test.SafeFprint(w, `{"device_code":"d9c126a7","user_code":"5055CC52","verification_url":"https://trakt.tv/activate","expires_in":600,"interval":5}`)
		})

		code, err := fetchNewDeviceCodeForClient(config, s.Client, &str.Options{})
		if assert.NoError(t, err) {
			assert.Equal(t, str.DeviceCode{DeviceCode: "d9c126a7", UserCode: "5055CC52", VerificationURL: "https://trakt.tv/activate", ExpiresIn: 600, Interval: 5}, *code)
		}
		assert.Equal(t, "client-id", *sent.ClientID)
	})

	t.Run("request fails", func(t *testing.T) {
		s := trakttest.Setup()
		defer s.Teardown()
		s.Mux.HandleFunc("/oauth/device/code", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		})

		code, err := fetchNewDeviceCodeForClient(config, s.Client, &str.Options{})
		assert.Nil(t, code)
		assert.ErrorContains(t, err, "generate new device code:")
	})
}

func TestDeviceCodeVerification(t *testing.T) {
	deviceToken := func(config *cfg.Config) *str.NewDeviceToken {
		code := "d9c126a7"
		return &str.NewDeviceToken{Code: &code, ClientID: &config.ClientID, ClientSecret: &config.ClientSecret}
	}

	t.Run("approved code stores the token and the user settings", func(t *testing.T) {
		config := credentialsConfig(t)
		s := trakttest.Setup()
		defer s.Teardown()
		var sent str.NewDeviceToken
		s.Mux.HandleFunc("/oauth/device/token", func(w http.ResponseWriter, r *http.Request) {
			test.AssertMethod(t, r, http.MethodPost)
			decodeBody(t, r, &sent)
			test.SafeFprint(w, tokenJSON("device-token", 7776000))
		})
		settingsAuth := ""
		s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, r *http.Request) {
			settingsAuth = r.Header.Get("Authorization")
			test.SafeFprint(w, testSettings)
		})

		options := &str.Options{}
		verified, err := deviceCodeVerification(deviceToken(config), s.Client, config, options)
		assert.NoError(t, err)
		assert.True(t, verified)
		assert.Equal(t, "d9c126a7", *sent.Code)
		assert.Equal(t, "client-secret", *sent.ClientSecret)
		assert.Equal(t, "device-token", options.Token.AccessToken)
		assert.Equal(t, "Bearer device-token", settingsAuth, "the settings request uses the new token")

		token, err := ReadTokenFromFile(config.TokenPath)
		if assert.NoError(t, err) {
			assert.Equal(t, "device-token", token.AccessToken)
		}
		settings, err := ReadUserSettingsFromFile(config.SettingsPath)
		if assert.NoError(t, err) {
			assert.Equal(t, "sean", *settings.User.Username)
		}
	})

	// final says whether the answer ends the polling
	for name, tc := range map[string]struct {
		status int
		final  string
	}{
		"pending":      {status: http.StatusBadRequest},
		"slow down":    {status: http.StatusTooManyRequests},
		"server error": {status: http.StatusInternalServerError},
		"denied":       {status: http.StatusTeapot, final: "device code denied"},
		"not found":    {status: http.StatusNotFound, final: "invalid device code"},
		"already used": {status: http.StatusConflict, final: "device code already used"},
		"expired":      {status: http.StatusGone, final: "device code expired"},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			config := credentialsConfig(t)
			s := trakttest.Setup()
			defer s.Teardown()
			s.Mux.HandleFunc("/oauth/device/token", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
			})

			options := &str.Options{}
			verified, err := deviceCodeVerification(deviceToken(config), s.Client, config, options)
			assert.False(t, verified)
			if tc.final == consts.EmptyString {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.final)
			}
			assert.Empty(t, options.Token.AccessToken)
			_, err = os.Stat(config.TokenPath)
			assert.True(t, os.IsNotExist(err), "no token file is written")
		})
	}
}

// A request that gets no response (network down) counts as a failed attempt, it used to panic on the nil response.
func TestDeviceCodeVerificationWithoutResponse(t *testing.T) {
	config := credentialsConfig(t)
	s := trakttest.Setup()
	s.Teardown()

	code := "d9c126a7"
	token := &str.NewDeviceToken{Code: &code, ClientID: &config.ClientID, ClientSecret: &config.ClientSecret}
	options := &str.Options{}
	assert.NotPanics(t, func() {
		verified, err := deviceCodeVerification(token, s.Client, config, options)
		assert.NoError(t, err)
		assert.False(t, verified)
	})
	assert.Empty(t, options.Token.AccessToken)
}

func TestVerifyCode(t *testing.T) {
	// interval 0 keeps the loop from sleeping; expires_in 0 ends it after one attempt
	device := &str.DeviceCode{DeviceCode: "d9c126a7"}

	t.Run("approved on the first attempt", func(t *testing.T) {
		config := credentialsConfig(t)
		s := trakttest.Setup()
		defer s.Teardown()
		attempts := 0
		s.Mux.HandleFunc("/oauth/device/token", func(w http.ResponseWriter, _ *http.Request) {
			attempts++
			test.SafeFprint(w, tokenJSON("device-token", 7776000))
		})
		s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, _ *http.Request) {
			test.SafeFprint(w, testSettings)
		})

		options := &str.Options{}
		assert.NoError(t, verifyCode(device, config, s.Client, options))
		assert.Equal(t, 1, attempts)
		assert.Equal(t, "device-token", options.Token.AccessToken)
	})

	t.Run("gives up when the code expires", func(t *testing.T) {
		config := credentialsConfig(t)
		s := trakttest.Setup()
		defer s.Teardown()
		attempts := 0
		s.Mux.HandleFunc("/oauth/device/token", func(w http.ResponseWriter, _ *http.Request) {
			attempts++
			w.WriteHeader(http.StatusBadRequest)
		})

		options := &str.Options{}
		assert.NoError(t, verifyCode(device, config, s.Client, options))
		assert.Equal(t, 1, attempts)
		assert.Empty(t, options.Token.AccessToken)
	})
}

// An interval that does not divide expires_in still ends the polling, it used to skip zero and poll forever.
func TestVerifyCodeStopsPastExpiry(t *testing.T) {
	config := credentialsConfig(t)
	s := trakttest.Setup()
	defer s.Teardown()
	attempts := 0
	s.Mux.HandleFunc("/oauth/device/token", func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		assert.NoError(t, verifyCode(&str.DeviceCode{DeviceCode: "d9c126a7", ExpiresIn: 1, Interval: 2}, config, s.Client, &str.Options{}))
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("verifyCode keeps polling after the code expired")
	}
	assert.Equal(t, 1, attempts)
}

// A denied, used, expired or unknown code ends the polling at once, it used to go on until the code's lifetime ran out.
func TestPoolNewDeviceCodeStopsOnFinalAnswer(t *testing.T) {
	openBrowser = func(string) error { return nil }
	t.Cleanup(func() { openBrowser = OpenBrowser })

	for name, status := range map[string]int{
		"denied":       http.StatusTeapot,
		"not found":    http.StatusNotFound,
		"already used": http.StatusConflict,
		"expired":      http.StatusGone,
	} {
		name, status := name, status
		t.Run(name, func(t *testing.T) {
			config := credentialsConfig(t)
			s := trakttest.Setup()
			defer s.Teardown()
			s.Mux.HandleFunc("/oauth/device/code", func(w http.ResponseWriter, _ *http.Request) {
				test.SafeFprint(w, `{"device_code":"d9c126a7","user_code":"5055CC52","verification_url":"https://trakt.tv/activate","expires_in":600,"interval":0}`)
			})
			attempts := 0
			s.Mux.HandleFunc("/oauth/device/token", func(w http.ResponseWriter, _ *http.Request) {
				attempts++
				if attempts == 1 {
					w.WriteHeader(status)
					return
				}
				// a later attempt would be approved, so the test ends even without the fix
				test.SafeFprint(w, tokenJSON("device-token", 7776000))
			})
			s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, _ *http.Request) {
				test.SafeFprint(w, testSettings)
			})

			options := &str.Options{}
			assert.Error(t, PoolNewDeviceCode(config, s.Client, options))
			assert.Equal(t, 1, attempts)
			assert.Empty(t, options.Token.AccessToken)
		})
	}
}

// A device code answer without a code is an error, it used to be returned as nil without an error.
func TestFetchNewDeviceCodeForClientWithoutCode(t *testing.T) {
	config := credentialsConfig(t)
	s := trakttest.Setup()
	defer s.Teardown()
	s.Mux.HandleFunc("/oauth/device/code", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	code, err := fetchNewDeviceCodeForClient(config, s.Client, &str.Options{})
	assert.Error(t, err)
	assert.Nil(t, code)
}

// The verification page is opened through the openBrowser hook, so the whole flow can run without a real browser.
func TestPoolNewDeviceCode(t *testing.T) {
	opened := ""
	openBrowser = func(url string) error {
		opened = url
		return nil
	}
	t.Cleanup(func() { openBrowser = OpenBrowser })

	config := credentialsConfig(t)
	s := trakttest.Setup()
	defer s.Teardown()
	s.Mux.HandleFunc("/oauth/device/code", func(w http.ResponseWriter, _ *http.Request) {
		test.SafeFprint(w, `{"device_code":"d9c126a7","user_code":"5055CC52","verification_url":"https://trakt.tv/activate","expires_in":0,"interval":0}`)
	})
	s.Mux.HandleFunc("/oauth/device/token", func(w http.ResponseWriter, _ *http.Request) {
		test.SafeFprint(w, tokenJSON("device-token", 7776000))
	})
	s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, _ *http.Request) {
		test.SafeFprint(w, testSettings)
	})

	options := &str.Options{}
	assert.NoError(t, PoolNewDeviceCode(config, s.Client, options))
	assert.Equal(t, "https://trakt.tv/activate", opened)
	assert.Equal(t, "device-token", options.Token.AccessToken)
}

func TestGenAppVersion(t *testing.T) {
	origVersion, origCommit, origDate, origBuiltBy := version, commit, date, builtBy
	t.Cleanup(func() { version, commit, date, builtBy = origVersion, origCommit, origDate, origBuiltBy })

	version, commit, date, builtBy = "1.22.0", "abc1234", "2026-10-01", "goreleaser"
	assert.EqualError(t, GenAppVersion(), "Version:\t1.22.0\nCommit:\t\tabc1234\nBuilt\t\t2026-10-01 by goreleaser")

	commit, date, builtBy = "none", "unknown", "unknown"
	assert.EqualError(t, GenAppVersion(), "Version:\t1.22.0\n")
	assert.Equal(t, "1.22.0", AppVersion())
}
