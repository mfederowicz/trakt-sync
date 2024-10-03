// Package cli for basic cli functions
package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

func fail(err string) {
	printer.Fprintln(os.Stderr, err)
}

// open browser for https://trakt.tv/activate code activation
func openBrowser(url string) error {
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

// check if user accept device code or not
func deviceCodeVerification(deviceToken *str.NewDeviceToken, oauth *internal.OauthService, config *cfg.Config) bool {
	token, resp, err := oauth.PoolForTheAccessToken(context.Background(), deviceToken)

	if err != nil {
		printer.Println("Error:", err)
		return false
	}

	if resp.StatusCode == http.StatusTeapot {
		fail("Error: Your device is not connected")
		return false
	}

	if resp.StatusCode == http.StatusOK {
		tokenjson, _ := json.Marshal(token)
		if err := os.WriteFile(config.TokenPath, tokenjson, consts.X644); err != nil {
			printer.Println(err.Error())
		}
	}

	return resp.StatusCode == http.StatusOK
}

// fetch new device code for client
func fetchNewDeviceCodeForClient(config *cfg.Config, oauth *internal.OauthService) (*str.DeviceCode, error) {
	code, resp, err := oauth.GenerateNewDeviceCodes(
		context.Background(),
		&str.NewDeviceCode{ClientID: &config.ClientID})

	if err != nil {
		return nil, printer.Errorf("Error generate new device code:" + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		return code, nil
	}

	return nil, nil
}

// PoolNewDeviceCode pool new device code (open browser and wait for correct code activation)
func PoolNewDeviceCode(config *cfg.Config, oauth *internal.OauthService) error {
	printer.Println("Polling for new device code...")

	device, err := fetchNewDeviceCodeForClient(config, oauth)
	if err != nil {
		return printer.Errorf("Error generate new device code:" + err.Error())
	}

	showCodeAndOpenBrowser(device)

	verifyCode(device, config, oauth)

	return nil
}

// show new device code to stdout and open browser
func showCodeAndOpenBrowser(device *str.DeviceCode) {
	printer.Println("Go to:" + device.VerificationURL)
	printer.Println("Enter code: " + device.UserCode)

	browserErr := openBrowser(device.VerificationURL)
	if browserErr != nil {
		fail("Error opening browser:" + browserErr.Error())
	}
}

// verify device code in loop with intervals
func verifyCode(device *str.DeviceCode, config *cfg.Config, oauth *internal.OauthService) {
	const (
		counterNoSeconds = 0
	)

	count := device.ExpiresIn
	for {
		token := &str.NewDeviceToken{
			Code:         &device.DeviceCode,
			ClientID:     &config.ClientID,
			ClientSecret: &config.ClientSecret,
		}
		if verified := deviceCodeVerification(token, oauth, config); verified {
			printer.Println("Device code verified!")
			break
		}
		count -= device.Interval
		if count == counterNoSeconds {
			printer.Println("Time out!")
			break
		}
		time.Sleep(time.Duration(device.Interval) * time.Second)
	}
}
