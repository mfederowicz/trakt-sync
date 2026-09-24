// Package cli for basic cli functions
package cli

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// openBrowser is replaced in tests so no real browser is started.
var openBrowser = OpenBrowser

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

// HandleVIPResponse handles the Trakt VIP responses (developer guide "VIP Methods"):
// 426 (VIP only) opens X-Upgrade-URL; 420 (account limit exceeded) opens it for a non-VIP user
// and only reports the limit for a VIP user. It returns nil when resp and err are not VIP related.
func HandleVIPResponse(resp *str.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusUpgradeRequired {
		return HandleUpgrade(resp)
	}

	var limits *internal.UpgradeUserLimitsError
	if errors.As(err, &limits) && limits.Response != nil {
		return handleAccountLimit(limits.Response.Header)
	}

	return nil
}

// HandleUpgrade opens X-Upgrade-URL (or the default VIP page) so the user can sign up for Trakt VIP.
func HandleUpgrade(r *str.Response) error {
	printer.Println("user account upgrade required")
	return openUpgradeURL(r.Header, "trakt vip required")
}

func handleAccountLimit(h http.Header) error {
	msg := "account limit exceeded"
	if limit := h.Get(internal.HeaderAccountLimit); limit != "" {
		msg = fmt.Sprintf("%s (limit: %s)", msg, limit)
	}
	// a VIP user already has the higher limits, so there is nothing to upgrade
	if h.Get(internal.HeaderVIPUser) == "true" {
		return errors.New(msg)
	}

	return openUpgradeURL(h, msg)
}

func openUpgradeURL(h http.Header, reason string) error {
	upgradeURL := h.Get(internal.HeaderUpgradeURL)
	if upgradeURL == "" {
		upgradeURL = internal.DefaultUpgradeURL
	}
	if err := openBrowser(upgradeURL); err != nil {
		return fmt.Errorf("%s, open %s to upgrade (browser error: %w)", reason, upgradeURL, err)
	}

	return fmt.Errorf("%s, browser opened: %s", reason, upgradeURL)
}
