// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
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
