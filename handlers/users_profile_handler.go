// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersProfileHandler struct for handler
type UsersProfileHandler struct{}

// Handle to handle users: profile action
func (UsersProfileHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("fetch profile for:" + options.UserName)
	stats, resp, err := fetchUserProfile(client, options)
	if err != nil {
		return fmt.Errorf("fetch user profile error:%w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found user profile for:%s", options.UserName)
	}

	printer.Printf("Found %s user profile\n", options.UserName)

	jsonData, err := json.MarshalIndent(stats, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func fetchUserProfile(client *trakt.Client, options *str.Options) (*str.UserProfile, *str.Response, error) {
	username := options.UserName
	profile, resp, err := client.Users.GetProfile(
		cli.ContextFromOptions(options),
		username,
	)

	return profile, resp, err
}
