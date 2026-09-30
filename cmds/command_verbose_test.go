// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestMaskSecret checks only the last chars of a secret are kept.
func TestMaskSecret(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		want   string
	}{
		{name: "empty", secret: "", want: "<not set>"},
		{name: "short", secret: "abcd", want: "****"},
		{name: "long", secret: "0123456789abcdef", want: "****cdef"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, maskSecret(tt.secret))
		})
	}
}

// TestProcessVerboseMasksCredentials checks -v never prints the access token or client id.
func TestProcessVerboseMasksCredentials(t *testing.T) {
	const token = "secret-access-token-1234"
	const clientID = "secret-client-id-5678"
	options := &str.Options{
		Verbose: true,
		Headers: map[string]any{
			"Authorization": "Bearer " + token,
			"trakt-api-key": clientID,
			"User-Agent":    "trakt-sync/1.19.1",
		},
	}

	out := captureStdout(t, func() { processVerbose(options) })

	assert.NotContains(t, out, token)
	assert.NotContains(t, out, clientID)
	assert.Contains(t, out, "Authorization header:****1234")
	assert.Contains(t, out, "trakt-api-key header:****5678")
	assert.Contains(t, out, "User-Agent header:trakt-sync/1.19.1")
}

// TestDebugLoggerOnlyWithVerbose checks the client debug output is printer.Println with -v and off without it.
func TestDebugLoggerOnlyWithVerbose(t *testing.T) {
	assert.Nil(t, debugLogger(false))
	out := captureStdout(t, func() { debugLogger(true)("fetch countries url:", "countries/movies") })
	assert.Equal(t, "fetch countries url: countries/movies\n", out)
}

// TestExecSetsDebugLogger checks Exec turns the client debug output on only with -v.
func TestExecSetsDebugLogger(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/verbose/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/verbose/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/verbose/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))

	tests := []struct {
		name    string
		verbose bool
	}{
		{name: "verbose", verbose: true},
		{name: "quiet", verbose: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)
			config := cfg.DefaultConfig()
			config.ClientID, config.ClientSecret = "a", "b"
			config.TokenPath, config.SettingsPath = "/verbose/token.json", "/verbose/user_settings.json"
			config.Verbose = tt.verbose
			client := internal.NewClient(nil)
			// a logger left from an earlier run must not survive a run without -v
			client.DebugLogger = func(...any) {}

			captureStdout(t, func() {
				assert.NoError(t, MoviesCmd.Exec(fs, client, config, []string{"-a", "no_such_action"}))
			})

			assert.Equal(t, tt.verbose, client.DebugLogger != nil)
		})
	}
}
