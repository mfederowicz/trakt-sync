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

// GenAppVersion gen app verrsion string
func GenAppVersion() error {
	const (
		First = 1
		EmptyBuildInfoLen = 0
	)
	var buildInfo string

	if date != "unknown" && builtBy != "unknown" {
		buildInfo = fmt.Sprintf("Built\t\t%s by %s", date, builtBy)
	}

	if commit != "none" {
		buildInfo = fmt.Sprintf("Commit:\t\t%s\n%s", commit, buildInfo)
	}

	if version == "dev" {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			version = bi.Main.Version
			if strings.HasPrefix(version, "v") {
				version = bi.Main.Version[First:]
			}
			if len(buildInfo) == EmptyBuildInfoLen {
				return fmt.Errorf("version %s", version)
			}
		}
	}

	return fmt.Errorf("Version:\t%s\n%s", version, buildInfo)
}
