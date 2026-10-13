// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// movies, shows, media and recommendations copy -theme to the options; without the flag the theme is empty.
func TestModuleThemeFlag(t *testing.T) {
	modules := []struct {
		command *Command
		action  string
	}{
		{command: MoviesCmd, action: "trending"},
		{command: ShowsCmd, action: "popular"},
		{command: MediaCmd, action: "anticipated"},
		{command: RecommendationsCmd, action: "movies"},
	}
	cases := []struct {
		name string
		args []string
		want string
	}{
		{name: "theme", args: []string{"-theme", "halloween"}, want: "halloween"},
		{name: "no flag"},
	}
	for _, m := range modules {
		for _, tt := range cases {
			m, tt := m, tt
			t.Run(m.command.Name+" "+tt.name, func(t *testing.T) {
				resetAllFlags()
				t.Cleanup(resetAllFlags)

				fs := afero.NewMemMapFs()
				tmpPath := "/tmp-" + m.command.Name + "-theme/"
				assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
				assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
				assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

				fileConfig := cfg.DefaultConfig()
				fileConfig.ClientID = "a"
				fileConfig.ClientSecret = "b"
				fileConfig.TokenPath = tmpPath + "token.json"
				fileConfig.SettingsPath = tmpPath + "user_settings.json"

				var got *str.Options
				command := &Command{Name: m.command.Name, Flag: m.command.Flag, Run: func(c *Command, _ ...string) error {
					got = c.UpdateOptionsWithCommandFlags(c.Options)
					return nil
				}}
				args := append([]string{"-a", m.action}, tt.args...)
				assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, args))
				assert.Equal(t, m.action, got.Action)
				assert.Equal(t, tt.want, got.Theme)
			})
		}
	}
}
