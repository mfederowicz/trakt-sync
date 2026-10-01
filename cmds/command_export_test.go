// Package cmds used for commands modules
package cmds

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/trakt/trakttest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestExportEpisodeWithoutTitle checks the episode exports with an episode that has no title: the API allows it, it panicked before.
func TestExportEpisodeWithoutTitle(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/export/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/export/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/export/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/export/token.json", "/export/user_settings.json"

	const item = `[{"id":1,"watched_at":"2026-01-02T03:04:05.000Z","listed_at":"2026-01-02T03:04:05.000Z","type":"episode",
		"episode":{"season":2,"number":5,"title":null,"ids":{"trakt":9,"imdb":"tt9","tmdb":90,"tvdb":900}},
		"show":{"title":"Tron","ids":{"trakt":1}}}]`
	for _, command := range []*Command{CollectionCmd, HistoryCmd, WatchlistCmd} {
		for _, format := range []string{consts.ImdbFormat, consts.TmdbFormat, "tvdb"} {
			command, format := command, format
			t.Run(command.Name+" "+format, func(t *testing.T) {
				resetAllFlags()
				t.Cleanup(resetAllFlags)

				setup := trakttest.Setup()
				defer setup.Teardown()
				setup.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte(item))
				})

				output := filepath.Join(t.TempDir(), "out.json")
				args := []string{"-o", output, "-t", "episodes", "-f", format}
				captureStdout(t, func() {
					assert.NoError(t, command.Exec(fs, setup.Client, config, args))
				})
				written, err := afero.ReadFile(afero.NewOsFs(), output)
				assert.NoError(t, err)
				assert.Contains(t, string(written), `"title": "`+consts.NoEpisodeTitle+`"`)
			})
		}
	}
}
