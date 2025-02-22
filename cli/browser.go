// Package cli for basic cli functions
package cli

import (
	"errors"
	"os"
	"os/exec"
	"runtime"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// OpenBrowser open browser for url
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return os.ErrNotExist
	}

	return cmd.Start()
}

// HandleUpgrade parses related headers, and returns upgradeUrl.
func HandleUpgrade(r *str.Response) error {
	upgradeURL := r.Header.Get("X-Upgrade-URL")
	printer.Println("user account upgrade required")
	browserErr := OpenBrowser(upgradeURL)
	if browserErr != nil {
		return errors.New("Error opening browser:" + browserErr.Error())
	}

	return errors.New("browser opened:" + upgradeURL)
}
