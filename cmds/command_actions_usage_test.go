// Package cmds used for commands modules
package cmds

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// captureStdout returns what fn prints to os.Stdout (printer writes there directly).
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r) // a read error only shortens the captured output, which the asserts catch
		done <- buf.String()
	}()
	fn()
	assert.NoError(t, w.Close())
	return <-done
}

// TestActionsUsageListsRegisteredActions checks that the help for an unknown action lists real actions only.
func TestActionsUsageListsRegisteredActions(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/usage/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/usage/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/usage/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/usage/token.json", "/usage/user_settings.json"

	tests := []struct {
		name    string
		cmd     *Command
		want    []string
		notWant []string
	}{
		{name: "movies", cmd: MoviesCmd, want: []string{consts.Updates, consts.Related}, notWant: []string{"updated", "releated"}},
		{name: "shows", cmd: ShowsCmd, want: []string{consts.Related, consts.ResetShowProgress}, notWant: []string{"releated", consts.Boxoffice, consts.Releases}},
		{name: "users", cmd: UsersCmd, want: []string{consts.FollowerRequests}, notWant: []string{"follow_request"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)
			out := captureStdout(t, func() {
				assert.EqualError(t, tt.cmd.Exec(fs, trakt.NewClient(nil), config, []string{"-a", "no_such_action"}),
					tt.cmd.Name+`: unknown action "no_such_action"`)
			})
			assert.Contains(t, out, "Available actions:")
			for _, a := range tt.want {
				assert.Contains(t, out, "  - "+a+"\n")
			}
			for _, a := range tt.notWant {
				assert.NotContains(t, out, "  - "+a+"\n")
			}
		})
	}
}

// TestUnknownActionOrTypeFails checks every module returns an error (exit status 1) after the usage
// for a wrong or missing -a, and the type modules for a wrong -t.
func TestUnknownActionOrTypeFails(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/unknown/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/unknown/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/unknown/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/unknown/token.json", "/unknown/user_settings.json"

	typeModules := map[string]bool{consts.Certifications: true, consts.Countries: true, consts.Genres: true, consts.Languages: true}
	for _, c := range Commands {
		c := c
		// help has no actions; collection, history and watchlist have no -a flag (a given -a already fails to parse)
		if c.Name == "help" || c.Name == consts.Collection || c.Name == consts.History || c.Name == consts.Watchlist {
			continue
		}
		t.Run(c.Name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)
			args, want := []string{"-a", "no_such_action"}, c.Name+`: unknown action "no_such_action"`
			if typeModules[c.Name] {
				args, want = []string{"-t", "no_such_type"}, c.Name+`: unknown type "no_such_type"`
			}
			out := captureStdout(t, func() {
				assert.EqualError(t, c.Exec(fs, trakt.NewClient(nil), config, args), want)
			})
			assert.Contains(t, out, "Usage: ./trakt-sync "+c.Name)
		})
	}

	t.Run("missing action", func(t *testing.T) {
		resetAllFlags()
		t.Cleanup(resetAllFlags)
		captureStdout(t, func() {
			assert.EqualError(t, TeamCmd.Exec(fs, trakt.NewClient(nil), config, []string{"-a", ""}), "team: no action given, use -a")
		})
	})
}
