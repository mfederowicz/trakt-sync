// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestRemovedActionNames checks the old hyphenated action names of calendars and search are unknown actions now.
func TestRemovedActionNames(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/actions/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/actions/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/actions/user_settings.json", []byte("{}"), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/actions/token.json", "/actions/user_settings.json"

	tests := []struct {
		name    string
		cmd     *Command
		action  string
		wantErr string
	}{
		{name: "calendars my-shows", cmd: CalendarsCmd, action: "my-shows", wantErr: `calendars: unknown action "my-shows"`},
		{name: "calendars all-new-shows", cmd: CalendarsCmd, action: "all-new-shows", wantErr: `calendars: unknown action "all-new-shows"`},
		{name: "calendars hot-releases", cmd: CalendarsCmd, action: "hot-releases", wantErr: `calendars: unknown action "hot-releases"`},
		{name: "search text-query", cmd: SearchCmd, action: "text-query", wantErr: `search: unknown action "text-query"`},
		{name: "search id-lookup", cmd: SearchCmd, action: "id-lookup", wantErr: `search: unknown action "id-lookup"`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)
			command := &Command{Name: tt.cmd.Name, Flag: tt.cmd.Flag, Run: tt.cmd.Run}
			var err error
			captureStdout(t, func() {
				err = command.Exec(fs, trakt.NewClient(nil), config, []string{"-a", tt.action})
			})
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}
