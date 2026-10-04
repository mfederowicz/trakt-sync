// Package cli for basic cli functions
package cli

import (
	"runtime/debug"
	"testing"
	"time"
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

// TestLocalBuildDate checks the build date is shown in the given zone for both ldflags layouts.
func TestLocalBuildDate(t *testing.T) {
	cest := time.FixedZone("CEST", 2*60*60)
	est := time.FixedZone("EST", -5*60*60)

	cases := []struct {
		name string
		in   string
		loc  *time.Location
		want string
	}{
		{name: "makefile layout", in: "2026-10-04 12:30 UTC", loc: cest, want: "2026-10-04 14:30 CEST"},
		{name: "goreleaser layout", in: "2026-10-04T12:30:45Z", loc: cest, want: "2026-10-04 14:30 CEST"},
		{name: "goreleaser layout with offset", in: "2026-10-04T14:30:45+02:00", loc: time.UTC, want: "2026-10-04 12:30 UTC"},
		{name: "day changes", in: "2026-10-04 02:00 UTC", loc: est, want: "2026-10-03 21:00 EST"},
		{name: "utc stays utc", in: "2026-10-04 12:30 UTC", loc: time.UTC, want: "2026-10-04 12:30 UTC"},
		{name: "unknown layout is kept", in: "2026-10-01", loc: cest, want: "2026-10-01"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := localBuildDate(tc.in, tc.loc); got != tc.want {
				t.Errorf("localBuildDate(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
