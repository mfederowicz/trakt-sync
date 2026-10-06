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

// TestExecStopsOnInvalidDateFlags checks that a date flag with a value that is not a date stops Exec
// before the command runs, and that a date or no value at all lets it run.
func TestExecStopsOnInvalidDateFlags(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/dates/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/dates/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/dates/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/dates/token.json", "/dates/user_settings.json"

	tests := []struct {
		name    string
		module  string
		args    []string
		wantErr string
	}{
		{name: "movies start_date typo", module: consts.Movies, args: []string{"-start_date", "2026-13-45"}, wantErr: `movies: invalid -start_date "2026-13-45", want YYYY-MM-DD`},
		{name: "people start_date word", module: consts.People, args: []string{"-start_date", "yesterday"}, wantErr: `people: invalid -start_date "yesterday", want YYYY-MM-DD`},
		{name: "shows start_date other order", module: consts.Shows, args: []string{"-start_date", "01-10-2026"}, wantErr: `shows: invalid -start_date "01-10-2026", want YYYY-MM-DD`},
		{name: "shows reset_at with a time", module: consts.Shows, args: []string{"-reset_at", "2026-10-01T10:00:00Z"}, wantErr: `shows: invalid -reset_at "2026-10-01T10:00:00Z", want YYYY-MM-DD`},
		{name: "sync start_at typo", module: consts.Sync, args: []string{"-start_at", "2026-02-30"}, wantErr: `sync: invalid -start_at "2026-02-30", want YYYY-MM-DD`},
		{name: "sync end_at typo", module: consts.Sync, args: []string{"-start_at", "2026-02-01", "-end_at", "2026-2-7"}, wantErr: `sync: invalid -end_at "2026-2-7", want YYYY-MM-DD`},
		{name: "users start_at=value", module: consts.Users, args: []string{"-start_at=soon"}, wantErr: `users: invalid -start_at "soon", want YYYY-MM-DD`},
		{name: "users end_at typo", module: consts.Users, args: []string{"-end_at", "2026/02/07"}, wantErr: `users: invalid -end_at "2026/02/07", want YYYY-MM-DD`},
		{name: "movies start_date", module: consts.Movies, args: []string{"-start_date", "2026-01-15"}},
		{name: "people start_date", module: consts.People, args: []string{"-start_date", "2026-01-15"}},
		{name: "shows start_date and reset_at", module: consts.Shows, args: []string{"-start_date", "2026-01-15", "-reset_at", "2026-01-20"}},
		{name: "sync start_at and end_at", module: consts.Sync, args: []string{"-start_at", "2026-01-15", "-end_at", "2026-01-20"}},
		{name: "users start_at and end_at", module: consts.Users, args: []string{"-start_at", "2026-01-15", "-end_at", "2026-01-20"}},
		{name: "movies without dates", module: consts.Movies, args: []string{}},
		{name: "sync without dates", module: consts.Sync, args: []string{}},
		{name: "empty value", module: consts.Shows, args: []string{"-reset_at", ""}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)
			var module *Command
			for _, c := range Commands {
				if c.Name == tt.module {
					module = c
				}
			}
			if module == nil {
				t.Fatalf("no command %q", tt.module)
			}
			ran := false
			cmd := &Command{Name: module.Name, Flag: module.Flag, Run: func(*Command, ...string) error { ran = true; return nil }}
			err := cmd.Exec(fs, trakt.NewClient(nil), config, tt.args)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.True(t, ran, "the command did not run")
				return
			}
			assert.EqualError(t, err, tt.wantErr)
			assert.False(t, ran, "the command ran with an invalid date")
		})
	}
}
