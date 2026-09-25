// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestExecTypeFromConfigFileAndFlag(t *testing.T) {
	tests := []struct {
		name     string
		fileType string
		args     []string
		want     string
	}{
		{name: "no config type, no -t", want: cfg.DefaultConfig().Type},
		{name: "config type, no -t", fileType: "shows", want: "shows"},
		{name: "config type and -t", fileType: "shows", args: []string{"-t", "episodes"}, want: "episodes"},
		{name: "no config type, -t", args: []string{"-t", "shows"}, want: "shows"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// flags share their values with the global flags, so reset them before and after
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			fs := afero.NewMemMapFs()
			tmpPath := "/tmp-type/"
			assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
			assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
			assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

			fileConfig := cfg.DefaultConfig()
			fileConfig.ClientID = "a"
			fileConfig.ClientSecret = "b"
			fileConfig.TokenPath = tmpPath + "token.json"
			fileConfig.SettingsPath = tmpPath + "user_settings.json"
			fileConfig.Type = tt.fileType

			var got string
			command := &Command{Name: "watchlist", Run: func(c *Command, _ ...string) error {
				got = c.Options.Type
				return nil
			}}
			assert.NoError(t, command.Exec(fs, internal.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.want, got)
		})
	}
}
