// Package cli for basic cli functions
package cli

import (
	"fmt"
	"os"
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
func GenAppVersion() {

	var buildInfo string

	if date != "unknown" && builtBy != "unknown" {
		buildInfo = fmt.Sprintf("Built\t\t%s by %s\n", date, builtBy)
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
				fmt.Printf("version %s\n", version)
				os.Exit(0)
			}
		}
	}

	fmt.Printf("Version:\t%s\n%s", version, buildInfo)
	os.Exit(0)

}
