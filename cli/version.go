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
				version = bi.Main.Version[1:]
			}
			if len(buildInfo) == 0 {
				return fmt.Errorf("version %s", version)
			}
		}
	}

	fmt.Printf("Version:\t%s\n%s", version, buildInfo)
	return nil
}
