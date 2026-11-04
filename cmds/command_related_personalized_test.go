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

func TestRelatedPersonalizedFlag(t *testing.T) {
	modules := []struct {
		name string
		cmd  *Command
	}{
		{name: consts.Movies, cmd: MoviesCmd},
		{name: consts.Shows, cmd: ShowsCmd},
	}
	cases := []struct {
		name string
		args []string
		want bool
	}{
		{name: "flag set", args: []string{"-a", "related", "-i", "tron-legacy-2010", "-personalized"}, want: true},
		{name: "flag not set", args: []string{"-a", "related", "-i", "tron-legacy-2010"}, want: false},
	}
	for _, m := range modules {
		for _, tc := range cases {
			m, tc := m, tc
			t.Run(m.name+" "+tc.name, func(t *testing.T) {
				resetAllFlags()
				t.Cleanup(resetAllFlags)

				fs := afero.NewMemMapFs()
				tmpPath := "/tmp-related-personalized/"
				assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
				assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
				assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

				fileConfig := cfg.DefaultConfig()
				fileConfig.ClientID = "a"
				fileConfig.ClientSecret = "b"
				fileConfig.TokenPath = tmpPath + "token.json"
				fileConfig.SettingsPath = tmpPath + "user_settings.json"

				var got *str.Options
				command := &Command{Name: m.name, Flag: m.cmd.Flag, Run: func(c *Command, _ ...string) error {
					got = c.UpdateOptionsWithCommandFlags(c.Options)
					return nil
				}}
				assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, tc.args))
				if assert.NotNil(t, got) {
					assert.Equal(t, "related", got.Action)
					assert.Equal(t, tc.want, got.Personalized)
				}
			})
		}
	}
}
