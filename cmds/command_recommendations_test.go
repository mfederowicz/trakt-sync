// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRecommendationsFlags(t *testing.T) {
	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-recommendations/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	tests := []struct {
		name string
		args []string
		want str.Options
	}{
		{name: "no ignore flags", args: []string{"-a", "movies"}, want: str.Options{}},
		{
			name: "all ignore flags",
			args: []string{"-a", "shows", "-ignore_collected", "true", "-ignore_watched", "true", "-ignore_watchlisted", "false", "-watch_window", "30"},
			want: str.Options{IgnoreCollected: "true", IgnoreWatched: "true", IgnoreWatchlisted: "false", WatchWindow: 30},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			var got *str.Options
			command := &Command{Name: consts.Recommendations, Flag: RecommendationsCmd.Flag, Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				return nil
			}}
			assert.NoError(t, command.Exec(fs, internal.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.want.IgnoreCollected, got.IgnoreCollected)
			assert.Equal(t, tt.want.IgnoreWatched, got.IgnoreWatched)
			assert.Equal(t, tt.want.IgnoreWatchlisted, got.IgnoreWatchlisted)
			assert.Equal(t, tt.want.WatchWindow, got.WatchWindow)
		})
	}
}
