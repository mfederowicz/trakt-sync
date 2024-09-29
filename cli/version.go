// Package cli for basic cli functions
package cli

import (
	"fmt"
	"runtime/debug"
	"strings"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

const (
	First             = 1
	EmptyBuildInfoLen = 0
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
		genDev(buildInfo)
	}

	return fmt.Errorf("Version:\t%s\n%s", version, buildInfo)
}

func genDev(info string) error {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		var version string = bi.Main.Version
		var versionNoPrefix string = bi.Main.Version[1:]

		if strings.HasPrefix(version, "v") {
			version = versionNoPrefix
		}
		if len(info) == EmptyBuildInfoLen {
			return fmt.Errorf("version %s", version)
		}
	}
	return nil
}
