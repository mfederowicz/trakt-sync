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

// users -hide reaches the options; without the flag nothing is set.
func TestUsersWatchlistHideFlag(t *testing.T) {
	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-users-hide/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "no flag", args: []string{"-a", "watchlist"}, want: ""},
		{name: "-hide value", args: []string{"-a", "watchlist", "-t", "shows", "-hide", "ended"}, want: "ended"},
		{name: "-hide=value", args: []string{"-a", "watchlist", "-hide=rated"}, want: "rated"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			var got *str.Options
			command := &Command{Name: consts.Users, Flag: UsersCmd.Flag, Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				return nil
			}}
			assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.want, got.HideItems)
			assert.False(t, got.Hide, "the recommendations hide switch stays off")
		})
	}
}
