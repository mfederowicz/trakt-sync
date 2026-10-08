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

func TestShowsLastActivityFlag(t *testing.T) {
	resetAllFlags()
	t.Cleanup(resetAllFlags)

	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-shows-last-activity/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "not set", args: []string{"-a", "collection_progress", "-i", "the-sopranos"}},
		{name: "collected", args: []string{"-a", "collection_progress", "-i", "the-sopranos", "-last_activity", "collected"}, want: "collected"},
		{name: "watched", args: []string{"-a", "watched_progress", "-i", "the-sopranos", "-last_activity", "watched"}, want: "watched"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			fileConfig := cfg.DefaultConfig()
			fileConfig.ClientID = "a"
			fileConfig.ClientSecret = "b"
			fileConfig.TokenPath = tmpPath + "token.json"
			fileConfig.SettingsPath = tmpPath + "user_settings.json"

			var got *str.Options
			command := &Command{Name: consts.Shows, Flag: ShowsCmd.Flag, Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				return nil
			}}
			assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.want, got.LastActivity)
		})
	}
}
