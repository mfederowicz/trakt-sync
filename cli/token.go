// Package cli for basic cli functions
package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/spf13/afero"
)

// ValidAccessToken valid if access_token is expired or not, and refresh if expired
func ValidAccessToken(config *cfg.Config, client *internal.Client, options *str.Options) bool {
	token, err := ReadTokenFromFile(config.TokenPath)
	if err != nil {
		printer.Println("Error reading token:", err)
		return false
	}
	if len(token.AccessToken) == consts.ZeroValue {
		printer.Println("Error: empty access token")
		return false
	}

	if token.Expired() {
		if refreshed := refreshToken(config, client, options); refreshed {
			printer.Println("Token refreshed!")
		}

		// Reload the updated token from the file
		token, err = ReadTokenFromFile(config.TokenPath)
		if err != nil {
			printer.Println("Error reading updated token:", err)
			return false
		}

		if refreshedSettings := RefreshUserSettings(config, client, options); refreshedSettings {
			printer.Println("User settings refreshed!")
		}
	}

	return !token.Expired()
}

// ReadUserSettingsFromFile reads user settings from the specified file
func ReadUserSettingsFromFile(filePath string) (*str.UserSettings, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("ooo")
		return nil, err
	}

	var settings str.UserSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// ReadTokenFromFile reads the token from the specified file
func ReadTokenFromFile(filePath string) (*str.Token, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var token str.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

// refresh access token to new one
func refreshToken(config *cfg.Config, client *internal.Client, options *str.Options) bool {
	token, err := ReadTokenFromFile(config.TokenPath)
	if err != nil {
		printer.Println("Error reading token:", err)
		return false
	}

	grantType := "refresh_token"

	currentToken := &str.CurrentDeviceToken{
		RefreshToken: &token.RefreshToken,
		ClientID:     &config.ClientID,
		ClientSecret: &config.ClientSecret,
		RedirectURI:  &config.RedirectURI,
		GrantType:    &grantType,
	}

	newToken, resp, err := client.Oauth.ExchangeRefreshTokenForAccessToken(
		client.BuildCtxFromOptions(options),
		currentToken,
	)

	if err != nil {
		printer.Println("Error exchange token:", err)
		return false
	}

	if resp.StatusCode == http.StatusOK {
		tokenjson, _ := json.Marshal(newToken)
		if err := os.WriteFile(config.TokenPath, tokenjson, consts.X644); err != nil {
			printer.Println(err.Error())
			return false
		}

		return true
	}

	return false
}

// RefreshUserSettings user settings
func RefreshUserSettings(config *cfg.Config, client *internal.Client, options *str.Options) bool {
	newSettings, resp, err := client.Users.RetrieveSettings(
		client.BuildCtxFromOptions(options),
	)

	if err != nil {
		printer.Println("Error get settings:", err)
		return false
	}

	if resp.StatusCode == http.StatusOK {
		settingsjson, _ := json.Marshal(newSettings)
		if err := os.WriteFile(config.SettingsPath, settingsjson, consts.X644); err != nil {
			printer.Println(err.Error())
			return false
		}

		return true
	}

	return false
}

// HandleToken process token check and refresh
func HandleToken(fs afero.Fs, config *cfg.Config, client *internal.Client, options str.Options) {
	if !ValidAccessToken(config, client, &options) {
		PoolNewDeviceCode(config, client, &options)
	}

	options, _ = cfg.OptionsFromConfig(fs, config)
	if len(options.Token.AccessToken) > consts.ZeroValue && options.UserSettings.User == nil {
		RefreshUserSettings(config, client, &options)
		printer.Println("User settings refreshed!")
	}
}
