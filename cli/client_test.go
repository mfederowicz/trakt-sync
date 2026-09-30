// Package cli for basic cli functions
package cli

import (
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/stretchr/testify/assert"
)

// TestNewClient checks the CLI client sends the config client id, the trakt-sync User-Agent and the token.
func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		wantAuth string
	}{
		{name: "with token", token: "file-token", wantAuth: "Bearer file-token"},
		{name: "without token", token: "", wantAuth: ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config := &cfg.Config{ClientID: "client-id"}
			req, err := NewClient(config, str.Token{AccessToken: tt.token}).NewRequest(http.MethodGet, "users/settings", nil)
			assert.NoError(t, err)
			assert.Equal(t, "client-id", req.Header.Get("trakt-api-key"))
			assert.Equal(t, UserAgent(), req.Header.Get("User-Agent"))
			assert.Equal(t, tt.wantAuth, req.Header.Get("Authorization"))
		})
	}
}
