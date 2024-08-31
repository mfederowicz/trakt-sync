// Package cli for basic cli functions
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

func fail(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
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
		fmt.Println("Error:", err)
		return false
	}

	if resp.StatusCode == 418 {
		fail("Error: Your device is not connected")
		return false
	}

	if resp.StatusCode == 200 {

		tokenjson, _ := json.Marshal(token)
		if err := os.WriteFile(config.TokenPath, tokenjson, 0644); err != nil {
			fmt.Println(err.Error())
		}

	}

	return resp.StatusCode == 200
}

// fetch new device code for client
func fetchNewDeviceCodeForClient(config *cfg.Config, oauth *internal.OauthService) *str.DeviceCode {

	code, resp, err := oauth.GenerateNewDeviceCodes(context.Background(), &str.NewDeviceCode{ClientID: &config.ClientID})

	if err != nil {
		fail("Error generate new device code:" + err.Error())
	}

	if resp.StatusCode == 200 {
		return code
	}

	return nil

}

// PoolNewDeviceCode pool new device code (open browser and wait for correct code activation)
func PoolNewDeviceCode(config *cfg.Config, oauth *internal.OauthService) {

	fmt.Println("Polling for new device code...")

	device := fetchNewDeviceCodeForClient(config, oauth)

	showCodeAndOpenBrowser(device)

	verifyCode(device, config, oauth)

}

// show new device code to stdout and open browser
func showCodeAndOpenBrowser(device *str.DeviceCode) {

	fmt.Println("Go to:" + device.VerificationURL)
	fmt.Println("Enter code: " + device.UserCode)

	browserErr := openBrowser(device.VerificationURL)
	if browserErr != nil {
		fail("Error opening browser:" + browserErr.Error())
	}

}

// verify device code in loop with intervals
func verifyCode(device *str.DeviceCode, config *cfg.Config, oauth *internal.OauthService) {

	count := device.ExpiresIn
	for {

		token := &str.NewDeviceToken{
			Code:         &device.DeviceCode,
			ClientID:     &config.ClientID,
			ClientSecret: &config.ClientSecret,
		}
		if verified := deviceCodeVerification(token, oauth, config); verified {
			fmt.Println("Device code verified!")
			break
		}
		count -= device.Interval
		if count == 0 {
			fmt.Println("Time out!")
			break
		}
		time.Sleep(time.Duration(device.Interval) * time.Second)
	}

}
