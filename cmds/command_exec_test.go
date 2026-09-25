// Package cmds used for commands modules
package cmds

import (
	"errors"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
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
			err := cmd.Exec(fs, internal.NewClient(nil), config, []string{})
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}
