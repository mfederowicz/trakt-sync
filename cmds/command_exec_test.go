// Package cmds used for commands modules
package cmds

import (
	"errors"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestExecReturnsRecoveredPanic checks that a panic in a command's Run comes back from Exec as an error.
func TestExecReturnsRecoveredPanic(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/exec/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/exec/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/exec/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/exec/token.json", "/exec/user_settings.json"

	tests := []struct {
		name    string
		run     func(*Command, ...string) error
		wantErr string
	}{
		{name: "panic", run: func(*Command, ...string) error { panic(errors.New("boom")) }, wantErr: "panic error:boom"},
		{name: "fatal", run: func(c *Command, _ ...string) error { c.Fatalf("stop"); return nil }, wantErr: "fatal error"},
		{name: "no panic", run: func(*Command, ...string) error { return nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &Command{Name: "exec_test", Run: tt.run}
			err := cmd.Exec(fs, trakt.NewClient(nil), config, []string{})
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

// TestExecStopsOnFlagErrors checks that an unknown flag, a leftover argument or -h stops Exec before the command runs.
func TestExecStopsOnFlagErrors(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/exec/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/exec/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/exec/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/exec/token.json", "/exec/user_settings.json"

	tests := []struct {
		name    string
		args    []string
		wantRun bool
		wantErr string
	}{
		{name: "unknown flag", args: []string{"-no_such_flag", "x"}, wantErr: "exec_test: flag provided but not defined: -no_such_flag"},
		{name: "unknown flag=value", args: []string{"-no_such_flag=x"}, wantErr: "exec_test: flag provided but not defined: -no_such_flag"},
		{name: "argument before the flags", args: []string{"foo", "-o", "out.json"}, wantErr: `exec_test: unexpected argument "foo"`},
		{name: "argument after the flags", args: []string{"-o", "out.json", "foo"}, wantErr: `exec_test: unexpected argument "foo"`},
		{name: "argument that is a module name", args: []string{"shows", "-o", "out.json"}, wantErr: `exec_test: unexpected argument "shows"`},
		{name: "argument that is a flag name", args: []string{"o", "out.json"}, wantErr: `exec_test: unexpected argument "o"`},
		{name: "help", args: []string{"-h"}},
		{name: "known flag", args: []string{"-o", "out.json"}, wantRun: true},
		{name: "known flag=value", args: []string{"-o=out.json"}, wantRun: true},
		{name: "empty value", args: []string{"-o", ""}, wantRun: true},
		{name: "value starting with a dash", args: []string{"-o", "-5"}, wantRun: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags) // "known flag" sets the global -o; don't leak it into later tests
			ran := false
			cmd := &Command{Name: "exec_test", Run: func(*Command, ...string) error { ran = true; return nil }}
			err := cmd.Exec(fs, trakt.NewClient(nil), config, tt.args)
			assert.Equal(t, tt.wantRun, ran)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

// TestModulesRuntimeReturnsErrors checks that ModulesRuntime returns an error (so main exits with 1)
// for an unknown command, an ambiguous prefix and a failing command, and nil on success.
func TestModulesRuntimeReturnsErrors(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/exec/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/exec/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/exec/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/exec/token.json", "/exec/user_settings.json"

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "unknown command", args: []string{"no_such_command"}, wantErr: `unknown command "no_such_command"`},
		{name: "ambiguous prefix", args: []string{"s"}, wantErr: `non-unique command prefix "s"`},
		{name: "failing command", args: []string{"checkin", "-no_such_flag"}, wantErr: "checkin: flag provided but not defined: -no_such_flag"},
		{name: "argument after a module", args: []string{"checkin", "foo"}, wantErr: `checkin: unexpected argument "foo"`},
		{name: "success", args: []string{"help"}},
		{name: "help takes a command name", args: []string{"help", "movies"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			err := ModulesRuntime(tt.args, fs, config, trakt.NewClient(nil))
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}
