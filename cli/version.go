// Package cli for basic cli functions
package cli

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"

	// readBuildInfo is replaced in tests
	readBuildInfo = debug.ReadBuildInfo
)

// GenAppVersion gen app verrsion string
func GenAppVersion() error {
	return fmt.Errorf("Version:\t%s\n%s", AppVersion(), genBuildInfo())
}

// AppVersion returns the app version from ldflags, or from build info for go install
func AppVersion() string {
	if version == "dev" {
		version = genDev(genBuildInfo())
	}

	return version
}

// UserAgent returns the User-Agent sent to the Trakt API, e.g. trakt-sync/1.19.1
func UserAgent() string {
	return consts.AppName + "/" + AppVersion()
}

func genBuildInfo() string {
	var buildInfo string
	if date != "unknown" && builtBy != "unknown" {
		buildInfo = fmt.Sprintf("Built\t\t%s by %s", localBuildDate(date, time.Local), builtBy)
	}

	if commit != "none" {
		buildInfo = fmt.Sprintf("Commit:\t\t%s\n%s", commit, buildInfo)
	}

	return buildInfo
}

// localBuildDate converts the build date from ldflags (Makefile or GoReleaser layout) to loc;
// a date in another layout is returned as it is
func localBuildDate(built string, loc *time.Location) string {
	for _, layout := range []string{consts.BuildDateFormat, time.RFC3339} {
		t, err := time.Parse(layout, built)
		if err == nil {
			return t.In(loc).Format(consts.BuildDateFormat)
		}
	}

	return built
}

func genDev(info string) string {
	ver := version
	bi, ok := readBuildInfo()
	// test binaries built with older Go versions have no main module version
	if ok && bi.Main.Version != consts.EmptyString {
		var version = bi.Main.Version
		var versionNoPrefix = bi.Main.Version[1:]

		if strings.HasPrefix(version, "v") {
			ver = versionNoPrefix
		}

		if len(info) == consts.EmptyBuildInfoLen {
			ver = version
		}
	}
	return ver
}
