// Package cli for basic cli functions
package cli

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

func fail(err string) {
	printer.Fprintln(os.Stderr, err)
}

// answers that end the polling: the device code can no longer be approved
var deviceCodeFinalErrors = map[int]error{
	http.StatusNotFound: errors.New("invalid device code"),
	http.StatusConflict: errors.New("device code already used"),
	http.StatusGone:     errors.New("device code expired"),
	http.StatusTeapot:   errors.New("device code denied, your device is not connected"),
}

// check if user accept device code or not, the error says that polling again is pointless
func deviceCodeVerification(deviceToken *str.NewDeviceToken, client *trakt.Client, config *cfg.Config, options *str.Options) (bool, error) {
	token, resp, err := client.Oauth.PollForAccessToken(ContextFromOptions(options), deviceToken)
	// no response at all, such as a network error
	if resp == nil {
		printer.Println("Error:", err)
		return false, nil
	}

	if finalErr, found := deviceCodeFinalErrors[resp.StatusCode]; found {
		return false, finalErr
	}

	if resp.StatusCode != http.StatusBadRequest && err != nil {
		printer.Println("Error:", err)
		return false, nil
	}

	if resp.StatusCode == http.StatusBadRequest {
		printer.Println("Wait for code")
		return false, nil
	}

	if resp.StatusCode == http.StatusOK {
		if err := writer.WritePrivateJSON(config.TokenPath, token); err != nil {
			printer.Println(err.Error())
		}

		options.Token = *token.ToToken()
		if len(options.Token.AccessToken) > consts.ZeroValue {
			client = client.WithAuthToken(options.Token.AccessToken)
		}
		RefreshUserSettings(config, client, options)
		printer.Println("User settings refreshed!")
	}

	return resp.StatusCode == http.StatusOK, nil
}

// fetch new device code for client
func fetchNewDeviceCodeForClient(config *cfg.Config, client *trakt.Client, options *str.Options) (*str.DeviceCode, error) {
	code, resp, err := client.Oauth.GenerateNewDeviceCodes(
		ContextFromOptions(options),
		&str.NewDeviceCode{ClientID: &config.ClientID})

	if err != nil {
		return nil, fmt.Errorf("generate new device code: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		return code, nil
	}

	return nil, fmt.Errorf("generate new device code: unexpected status %d", resp.StatusCode)
}

// PoolNewDeviceCode pool new device code (open browser and wait for correct code activation)
func PoolNewDeviceCode(config *cfg.Config, client *trakt.Client, options *str.Options) error {
	printer.Println("Polling for new device code...")

	device, err := fetchNewDeviceCodeForClient(config, client, options)
	if err != nil {
		return err
	}

	showCodeAndOpenBrowser(device)

	return verifyCode(device, config, client, options)
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
func verifyCode(device *str.DeviceCode, config *cfg.Config, client *trakt.Client, options *str.Options) error {
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
		verified, err := deviceCodeVerification(token, client, config, options)
		if err != nil {
			return err
		}
		if verified {
			printer.Println("Device code verified!")
			break
		}
		count -= device.Interval
		if count <= counterNoSeconds {
			printer.Println("Time out!")
			break
		}
		time.Sleep(time.Duration(device.Interval) * time.Second)
	}

	return nil
}
