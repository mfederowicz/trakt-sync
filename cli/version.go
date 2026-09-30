// Package cli for basic cli functions
package cli

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/mfederowicz/trakt-sync/consts"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
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
		buildInfo = fmt.Sprintf("Built\t\t%s by %s", date, builtBy)
	}

	if commit != "none" {
		buildInfo = fmt.Sprintf("Commit:\t\t%s\n%s", commit, buildInfo)
	}

	return buildInfo
}

func genDev(info string) string {
	ver := version
	bi, ok := debug.ReadBuildInfo()
	if ok {
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
