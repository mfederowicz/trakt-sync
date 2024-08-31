// Package cli for basic cli functions
package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// ValidAccessToken valid if access_token is expired or not, and refresh if expired
func ValidAccessToken(config *cfg.Config, oauth *internal.OauthService) bool {

	token, err := ReadTokenFromFile(config.TokenPath)
	if err != nil {
		fmt.Println("Error reading token:", err)
		return false
	}

	if token.Expired() {

		if refreshed := refreshToken(config, oauth); refreshed {
			fmt.Println("Token refresed!")
		}

		// Reload the updated token from the file
		token, err = ReadTokenFromFile(config.TokenPath)
		if err != nil {
			fmt.Println("Error reading updated token:", err)
			return false
		}
	}

	return !token.Expired()

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
func refreshToken(config *cfg.Config, oauth *internal.OauthService) bool {

	token, err := ReadTokenFromFile(config.TokenPath)
	if err != nil {
		fmt.Println("Error reading token:", err)
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

	newToken, resp, err := oauth.ExchangeRefreshTokenForAccessToken(
		context.Background(),
		currentToken,
	)

	if err != nil {
		fmt.Println("Error exchange token:", err)
		return false
	}

	if resp.StatusCode == 200 {

		tokenjson, _ := json.Marshal(newToken)
		if err := os.WriteFile(config.TokenPath, tokenjson, 0644); err != nil {
			fmt.Println(err.Error())
			return false
		}

		return true

	}

	return false

}
