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
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersWatchedHandler struct for handler
type UsersWatchedHandler struct{}

// Handle to handle users: watched action
func (UsersWatchedHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("fetch watched for:" + options.UserName + " and type:" + options.Type)
	watched, resp, err := fetchUsersWatched(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch user watched error:%w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found user watched for:%s, type:%s", options.UserName, options.Type)
	}

	printer.Printf("Found %s user watched type:%s\n", options.UserName, options.Type)

	jsonData, err := json.MarshalIndent(watched, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func fetchUsersWatched(client *trakt.Client, options *str.Options, page int) ([]*str.UserWatched, *str.Response, error) {
	username := options.UserName
	watchType := options.Type
	opts := uri.ListOptions{Page: page, Limit: consts.PerPage, Extended: options.ExtendedInfo}
	watched, resp, err := client.Users.GetWatched(
		cli.ContextFromOptions(options),
		username,
		watchType,
		&opts,
	)
	if err != nil {
		return nil, resp, fmt.Errorf("fetch user watched error:%w", err)
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		waitPageDelay()

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, resp, err := fetchUsersWatched(client, options, nextPage)
		if err != nil {
			return nil, resp, err
		}

		// Append items from the next page to the current page
		watched = append(watched, nextPageItems...)
	}

	return watched, resp, nil
}
