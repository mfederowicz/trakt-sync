// Package cli for basic cli functions
package cli

import (
	"testing"
)

// TestUserAgent checks the User-Agent names trakt-sync and the ldflags version.
func TestUserAgent(t *testing.T) {
	orig := version
	t.Cleanup(func() { version = orig })
	version = "1.19.1"

	if got, want := UserAgent(), "trakt-sync/1.19.1"; got != want {
		t.Errorf("UserAgent() = %q, want %q", got, want)
	}
}
