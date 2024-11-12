// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersWatchedHandler struct for handler
type UsersWatchedHandler struct{}

// Handle to handle users: watched action
func (UsersWatchedHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("fetch watched for:" + options.UserName + " and type:" + options.Type)
	watched, resp, err := fetchUsersWatched(client, options)
	if err != nil {
		return fmt.Errorf("fetch user watched error:%w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found user watched for:%s, type:%s", options.UserName, options.Type)
	}

	printer.Printf("Found %s user watched type:%s\n", options.UserName, options.Type)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(watched, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func fetchUsersWatched(client *internal.Client, options *str.Options) ([]*str.UserWatched, *str.Response, error) {
	username := options.UserName
	watchType := options.Type
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	watched, resp, err := client.Users.GetWatched(
		context.Background(),
		&username,
		&watchType,
		&opts,
	)

	return watched, resp, err
}
