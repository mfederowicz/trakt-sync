// Package cmds used for commands modules
package cmds

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/trakt/trakttest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestScrobbleProgressFlag checks scrobble -progress reaches the request body (it was never sent before).
func TestScrobbleProgressFlag(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/progress/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/progress/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/progress/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/progress/token.json", "/progress/user_settings.json"

	tests := []struct {
		name     string
		progress []string
		want     *float64
	}{
		{name: "with -progress", progress: []string{"-progress", "10.25"}, want: func() *float64 { v := 10.25; return &v }()},
		{name: "without -progress", want: nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			setup := trakttest.Setup()
			defer setup.Teardown()
			var sent struct {
				Progress *float64 `json:"progress"`
			}
			calls := 0
			setup.Mux.HandleFunc("/scrobble/start", func(w http.ResponseWriter, r *http.Request) {
				calls++
				assert.Equal(t, http.MethodPost, r.Method)
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(body, &sent))
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":1,"action":"start"}`))
			})

			args := append([]string{"-o", filepath.Join(t.TempDir(), "out.json"), "-a", "start", "-t", "episode", "-i", "73629"}, tt.progress...)
			captureStdout(t, func() {
				assert.NoError(t, ScrobbleCmd.Exec(fs, setup.Client, config, args))
			})
			assert.Equal(t, 1, calls)
			assert.Equal(t, tt.want, sent.Progress)
		})
	}
}
