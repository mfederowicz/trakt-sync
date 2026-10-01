// Package cli for basic cli functions
package cli

import (
	"encoding/json"
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
		assert.ErrorContains(t, err, "Error generate new device code:")
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
		assert.True(t, deviceCodeVerification(deviceToken(config), s.Client, config, options))
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

	for name, status := range map[string]int{
		"pending":       http.StatusBadRequest,
		"not connected": http.StatusTeapot,
		"not found":     http.StatusNotFound,
		"already used":  http.StatusConflict,
		"expired":       http.StatusGone,
		"slow down":     http.StatusTooManyRequests,
		"server error":  http.StatusInternalServerError,
	} {
		name, status := name, status
		t.Run(name, func(t *testing.T) {
			config := credentialsConfig(t)
			s := trakttest.Setup()
			defer s.Teardown()
			s.Mux.HandleFunc("/oauth/device/token", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
			})

			options := &str.Options{}
			assert.False(t, deviceCodeVerification(deviceToken(config), s.Client, config, options))
			assert.Empty(t, options.Token.AccessToken)
			_, err := os.Stat(config.TokenPath)
			assert.True(t, os.IsNotExist(err), "no token file is written")
		})
	}
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
		verifyCode(device, config, s.Client, options)
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
		verifyCode(device, config, s.Client, options)
		assert.Equal(t, 1, attempts)
		assert.Empty(t, options.Token.AccessToken)
	})
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
