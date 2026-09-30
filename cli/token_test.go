// Package cli for basic cli functions
package cli

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestValidAccessTokenRefreshesSettingsWithNewToken checks the user settings request after a token refresh
// sends the refreshed token, not the expired one the client was built with.
func TestValidAccessTokenRefreshesSettingsWithNewToken(t *testing.T) {
	dir := t.TempDir()
	config := &cfg.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		TokenPath:    filepath.Join(dir, "token.json"),
		SettingsPath: filepath.Join(dir, "user_settings.json"),
	}
	expired := `{"access_token":"old-token","refresh_token":"refresh","expires_in":1,"created_at":1}`
	assert.NoError(t, os.WriteFile(config.TokenPath, []byte(expired), consts.X600))

	s := internal.Setup()
	defer s.Teardown()

	s.Mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		created := strconv.FormatInt(time.Now().Unix(), consts.BaseInt)
		test.SafeFprint(w, `{"access_token":"new-token","refresh_token":"refresh2","expires_in":7776000,"created_at":`+created+`}`)
	})
	settingsAuth := ""
	s.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		settingsAuth = r.Header.Get("Authorization")
		test.SafeFprint(w, `{"user":{"username":"sean"}}`)
	})

	client := s.Client.WithClientID(config.ClientID).WithAuthToken("old-token")

	assert.True(t, ValidAccessToken(config, client, &str.Options{}))
	assert.Equal(t, "Bearer new-token", settingsAuth)
}

// TestHandleTokenReturnsClientWithToken checks HandleToken hands back a client that sends the token from token.json.
func TestHandleTokenReturnsClientWithToken(t *testing.T) {
	dir := t.TempDir()
	config := cfg.DefaultConfig()
	config.ClientID = "client-id"
	config.TokenPath = filepath.Join(dir, "token.json")
	config.SettingsPath = filepath.Join(dir, "user_settings.json")
	created := strconv.FormatInt(time.Now().Unix(), consts.BaseInt)
	valid := `{"access_token":"file-token","refresh_token":"refresh","expires_in":7776000,"created_at":` + created + `}`
	assert.NoError(t, os.WriteFile(config.TokenPath, []byte(valid), consts.X600))
	assert.NoError(t, os.WriteFile(config.SettingsPath, []byte(`{"user":{"username":"sean"}}`), consts.X600))

	s := internal.Setup()
	defer s.Teardown()

	client := HandleToken(afero.NewOsFs(), config, s.Client.WithClientID(config.ClientID), str.Options{})

	req, err := client.NewRequest(http.MethodGet, "users/settings", nil)
	assert.NoError(t, err)
	assert.Equal(t, "Bearer file-token", req.Header.Get("Authorization"))
}
