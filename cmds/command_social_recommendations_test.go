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

func TestSocialRecommendationsFlags(t *testing.T) {
	t.Cleanup(resetAllFlags)

	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-social/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	var got *str.Options
	command := &Command{Name: consts.SocialRecommendations, Flag: SocialRecommendationsCmd.Flag, Run: func(c *Command, _ ...string) error {
		got = c.UpdateOptionsWithCommandFlags(c.Options)
		return nil
	}}
	args := []string{"-a", "shows", "-ignore_collected", "true", "-ignore_watched", "true", "-ignore_watchlisted", "false", "-watch_window", "30"}
	assert.NoError(t, command.Exec(fs, internal.NewClient(nil), fileConfig, args))
	assert.Equal(t, "shows", got.Action)
	assert.Equal(t, "true", got.IgnoreCollected)
	assert.Equal(t, "true", got.IgnoreWatched)
	assert.Equal(t, "false", got.IgnoreWatchlisted)
	assert.Equal(t, 30, got.WatchWindow)
	assert.Equal(t, "export_social_recommendations_shows.json", got.Output)
}
