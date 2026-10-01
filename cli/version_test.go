// Package cli for basic cli functions
package cli

import (
	"runtime/debug"
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

// TestAppVersionFromBuildInfo checks the version taken from the build info when no ldflags version is set.
func TestAppVersionFromBuildInfo(t *testing.T) {
	cases := []struct {
		name        string
		mainVersion string
		ok          bool
		commit      string
		want        string
	}{
		{name: "no build info", ok: false, want: "dev"},
		{name: "empty main version", ok: true, want: "dev"},
		{name: "go install keeps the tag", mainVersion: "v1.21.0", ok: true, want: "v1.21.0"},
		{name: "commit set drops the v prefix", mainVersion: "v1.21.0", ok: true, commit: "abc1234", want: "1.21.0"},
		{name: "local build", mainVersion: "(devel)", ok: true, want: "(devel)"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			origVersion, origCommit, origRead := version, commit, readBuildInfo
			t.Cleanup(func() { version, commit, readBuildInfo = origVersion, origCommit, origRead })

			version, commit = "dev", "none"
			if tc.commit != "" {
				commit = tc.commit
			}
			readBuildInfo = func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{Main: debug.Module{Version: tc.mainVersion}}, tc.ok
			}

			if got := AppVersion(); got != tc.want {
				t.Errorf("AppVersion() = %q, want %q", got, tc.want)
			}
		})
	}
}
