// Package cmds used for commands modules
package cmds

import (
	"flag"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestDeprecatedFlags checks the dead flags still parse, print a deprecation note and leave the command running.
func TestDeprecatedFlags(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/deprecated/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/deprecated/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/deprecated/user_settings.json", []byte("{}"), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/deprecated/token.json", "/deprecated/user_settings.json"

	tests := []struct {
		name string
		cmd  *Command
		args []string
		want string
	}{
		{name: "studio_ids", cmd: MoviesCmd, args: []string{"-a", "trending", "-studio_ids", "1,2"}, want: "flag -studio_ids is no longer supported by the Trakt API and is ignored\n"},
		{name: "languages is a filter again", cmd: ShowsCmd, args: []string{"-a", "trending", "-languages", "en"}, want: ""},
		{name: "query", cmd: SearchCmd, args: []string{"-a", "text_query", "-query", "tron"}, want: "flag -query is deprecated, use -q\n"},
		{name: "notes_id", cmd: NotesCmd, args: []string{"-a", "note", "-notes_id", "5"}, want: "flag -notes_id is deprecated, use -i\n"},
		{name: "scrobble delete", cmd: ScrobbleCmd, args: []string{"-a", "start", "-delete"}, want: "flag -delete is deprecated and is ignored\n"},
		{name: "delete stays valid elsewhere", cmd: NotesCmd, args: []string{"-a", "note", "-i", "5", "-delete"}, want: ""},
		{name: "no deprecated flag", cmd: MoviesCmd, args: []string{"-a", "trending"}, want: ""},
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
			out := captureStdout(t, func() {
				assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), config, tt.args))
			})
			assert.True(t, ran, "the command must still run")
			if tt.want == "" {
				assert.NotContains(t, out, "deprecated")
				assert.NotContains(t, out, "no longer supported")
			} else {
				assert.Contains(t, out, tt.want)
			}
		})
	}
}

// TestFlagSetIn checks the before-the-module-name check on a separate flag set (setting flag.CommandLine would leak).
func TestFlagSetIn(t *testing.T) {
	fs := flag.NewFlagSet("main", flag.ContinueOnError)
	fs.String("studio_ids", "", "")
	fs.String("query", "", "")
	assert.NoError(t, fs.Parse([]string{"-studio_ids", "1"}))
	assert.True(t, flagSetIn(fs, "studio_ids"))
	assert.False(t, flagSetIn(fs, "query"))
}
