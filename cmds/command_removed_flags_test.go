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

// TestRemovedFlags checks the removed dead flags stop the run as unknown flags, and their neighbours still work.
func TestRemovedFlags(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/removed/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/removed/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/removed/user_settings.json", []byte("{}"), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/removed/token.json", "/removed/user_settings.json"

	tests := []struct {
		name    string
		cmd     *Command
		args    []string
		wantErr string
	}{
		{name: "studio_ids", cmd: MoviesCmd, args: []string{"-a", "trending", "-studio_ids", "1,2"}, wantErr: "movies: flag provided but not defined: -studio_ids"},
		{name: "query", cmd: SearchCmd, args: []string{"-a", "text_query", "-query", "tron"}, wantErr: "search: flag provided but not defined: -query"},
		{name: "notes_id", cmd: NotesCmd, args: []string{"-a", "note", "-notes_id", "5"}, wantErr: "notes: flag provided but not defined: -notes_id"},
		{name: "scrobble delete", cmd: ScrobbleCmd, args: []string{"-a", "start", "-delete"}, wantErr: "scrobble: flag provided but not defined: -delete"},
		{name: "search q stays", cmd: SearchCmd, args: []string{"-a", "text_query", "-q", "tron"}},
		{name: "languages stays", cmd: ShowsCmd, args: []string{"-a", "trending", "-languages", "en"}},
		{name: "notes i and delete stay", cmd: NotesCmd, args: []string{"-a", "note", "-i", "5", "-delete"}},
		{name: "no removed flag", cmd: MoviesCmd, args: []string{"-a", "trending"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)
			ran := false
			command := &Command{Name: tt.cmd.Name, Flag: tt.cmd.Flag, Run: func(*Command, ...string) error {
				ran = true
				return nil
			}}
			err := command.Exec(fs, trakt.NewClient(nil), config, tt.args)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.True(t, ran, "the command must run")
				return
			}
			assert.EqualError(t, err, tt.wantErr)
			assert.False(t, ran, "the command must not run")
		})
	}
}
