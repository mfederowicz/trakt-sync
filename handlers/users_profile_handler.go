// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersProfileHandler struct for handler
type UsersProfileHandler struct{}

// Handle to handle users: profile action
func (UsersProfileHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("fetch profile for:" + options.UserName)
	stats, resp, err := fetchUserProfile(client, options)
	if err != nil {
		return fmt.Errorf("fetch user profile error:%w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found user profile for:%s", options.UserName)
	}

	printer.Printf("Found %s user profile\n", options.UserName)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(stats, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func fetchUserProfile(client *internal.Client, options *str.Options) (*str.UserProfile, *str.Response, error) {
	username := options.UserName
	profile, resp, err := client.Users.GetProfile(
		client.BuildCtxFromOptions(options),
		&username,
	)

	return profile, resp, err
}
