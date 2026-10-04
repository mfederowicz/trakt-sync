// Package cmds used for commands modules
package cmds

import (
	"flag"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// ignore_* and watch_window of the config file reach the options of the modules that send them; a flag wins.
func TestConfigFileIgnoreOptions(t *testing.T) {
	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-config-ignore/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))

	fromFile := str.Options{IgnoreCollected: "true", IgnoreWatched: "true", IgnoreWatchlisted: "true", WatchWindow: 14}
	modules := []struct {
		name   string
		flags  *flag.FlagSet
		action []string
	}{
		{name: consts.Recommendations, flags: &RecommendationsCmd.Flag, action: []string{"-a", "movies"}},
		{name: consts.SocialRecommendations, flags: &SocialRecommendationsCmd.Flag, action: []string{"-a", "movies"}},
		{name: consts.SmartLists, flags: &SmartListsCmd.Flag, action: []string{"-a", "items", "-i", "top-sci-fi"}},
		{name: consts.Users, flags: &UsersCmd.Flag, action: []string{"-a", "activities", "-t", "friends"}},
	}
	cases := []struct {
		name string
		args []string
		want str.Options
	}{
		{name: "config file only", want: fromFile},
		{name: "flag wins", args: []string{"-ignore_watched", "false"},
			want: str.Options{IgnoreCollected: "true", IgnoreWatched: "false", IgnoreWatchlisted: "true", WatchWindow: 14}},
	}

	for _, m := range modules {
		for _, tt := range cases {
			m, tt := m, tt
			t.Run(m.name+" "+tt.name, func(t *testing.T) {
				resetAllFlags()
				t.Cleanup(resetAllFlags)

				fileConfig := cfg.DefaultConfig()
				fileConfig.ClientID = "a"
				fileConfig.ClientSecret = "b"
				fileConfig.TokenPath = tmpPath + "token.json"
				fileConfig.SettingsPath = tmpPath + "user_settings.json"
				fileConfig.IgnoreCollected = cfg.BoolText(fromFile.IgnoreCollected)
				fileConfig.IgnoreWatched = cfg.BoolText(fromFile.IgnoreWatched)
				fileConfig.IgnoreWatchlisted = cfg.BoolText(fromFile.IgnoreWatchlisted)
				fileConfig.WatchWindow = fromFile.WatchWindow

				var got *str.Options
				command := &Command{Name: m.name, Flag: *m.flags, Run: func(c *Command, _ ...string) error {
					got = c.UpdateOptionsWithCommandFlags(c.Options)
					return nil
				}}
				args := append(append([]string{}, m.action...), tt.args...)
				assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, args))
				assert.Equal(t, tt.want.IgnoreCollected, got.IgnoreCollected)
				assert.Equal(t, tt.want.IgnoreWatched, got.IgnoreWatched)
				assert.Equal(t, tt.want.IgnoreWatchlisted, got.IgnoreWatchlisted)
				assert.Equal(t, tt.want.WatchWindow, got.WatchWindow)
			})
		}
	}
}
