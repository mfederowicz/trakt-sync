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
	var buildInfo string
	if date != "unknown" && builtBy != "unknown" {
		buildInfo = fmt.Sprintf("Built\t\t%s by %s", date, builtBy)
	}

	if commit != "none" {
		buildInfo = fmt.Sprintf("Commit:\t\t%s\n%s", commit, buildInfo)
	}

	if version == "dev" {
		version = genDev(buildInfo)
	}

	return fmt.Errorf("Version:\t%s\n%s", version, buildInfo)
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
