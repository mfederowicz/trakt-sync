// Package cmds used for commands modules
package cmds

import (
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestCheckinRequests checks which items checkin looks up and how it reports an active checkin (409).
func TestCheckinRequests(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/chk/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/chk/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/chk/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/chk/token.json", "/chk/user_settings.json"

	showEpisode := []string{"-a", "show_episode", "-trakt_id", "136121", "-episode_code", "1x5"}
	tests := []struct {
		name       string
		args       []string
		status     int
		body       string
		wantLookup string
		wantErr    string
	}{
		{name: "movie", args: []string{"-a", "movie", "-trakt_id", "28"}, status: http.StatusCreated, body: `{"id":1}`, wantLookup: "/movies/28"},
		{name: "show_episode", args: showEpisode, status: http.StatusCreated, body: `{"id":1}`, wantLookup: "/shows/136121"},
		{name: "show_episode active checkin", args: showEpisode, status: http.StatusConflict, body: `{"expires_at":"2026-01-15T10:00:00.000Z"}`, wantLookup: "/shows/136121", wantErr: "exists, expires:"},
		{name: "show_episode active checkin without expires_at", args: showEpisode, status: http.StatusConflict, body: `{}`, wantLookup: "/shows/136121", wantErr: "exists"},
		{name: "show_episode abs active checkin without expires_at", args: []string{"-a", "show_episode", "-trakt_id", "37696", "-episode_abs", "6"}, status: http.StatusConflict, body: `{}`, wantLookup: "/shows/37696", wantErr: "exists"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			setup := internal.Setup()
			defer setup.Teardown()
			var gotLookup string
			setup.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(`{"connections":{}}`))
			})
			for _, path := range []string{"/movies/", "/shows/"} {
				setup.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
					gotLookup = r.URL.Path
					_, _ = w.Write([]byte(`{"title":"Test","ids":{"trakt":1}}`))
				})
			}
			setup.Mux.HandleFunc("/checkin", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			})

			err := CheckinCmd.Exec(fs, setup.Client, config, tt.args)
			assert.Equal(t, tt.wantLookup, gotLookup)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}
