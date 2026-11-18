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

// TestUsersCommentSortDefault checks that the users comment actions send no sort without -s, so Trakt sorts by
// newest, and pass -s through when it is set.
func TestUsersCommentSortDefault(t *testing.T) {
	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-users-comment-sort/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = tmpPath+"token.json", tmpPath+"user_settings.json"

	tests := []struct {
		name     string
		args     []string
		wantSort string
	}{
		{name: "list_comments default", args: []string{"-a", "list_comments", "-i", "star-wars"}},
		{name: "list_comments -s", args: []string{"-a", "list_comments", "-i", "star-wars", "-s", "likes"}, wantSort: "likes"},
		{name: "watchlist_comments default", args: []string{"-a", "watchlist_comments"}},
		{name: "watchlist_comments -s", args: []string{"-a", "watchlist_comments", "-s", "oldest"}, wantSort: "oldest"},
		{name: "favorites_comments default", args: []string{"-a", "favorites_comments"}},
		{name: "favorites_comments -s", args: []string{"-a", "favorites_comments", "-s", "replies"}, wantSort: "replies"},
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
			assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), config, tt.args))
			assert.Equal(t, tt.wantSort, got.Sort)
		})
	}
}
