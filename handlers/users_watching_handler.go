// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
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

// UsersWatchingHandler struct for handler
type UsersWatchingHandler struct{ common CommonLogic }

// Handle to handle users: watching action
func (m UsersWatchingHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.UserName) == consts.ZeroValue {
		return errors.New(consts.EmptyUserNameMsg)
	}

	result, resp, err := m.usersWatching(client, options)

	if resp != nil && resp.StatusCode == http.StatusNoContent {
		return fmt.Errorf("watching error:%s", consts.UserNotWatchingAnything)
	}

	if err != nil {
		return fmt.Errorf("watching error:%w", err)
	}
	printer.Println("write result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (UsersWatchingHandler) usersWatching(client *trakt.Client, options *str.Options) (*str.WatchingResult, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	user := options.UserName
	result, resp, err := client.Users.GetWatching(
		cli.ContextFromOptions(options),
		&user,
		&opts,
	)

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return nil, resp, fmt.Errorf("user not found:%s", options.UserName)
	}
	if err != nil {
		return nil, resp, fmt.Errorf("block error:%w", err)
	}

	return result, resp, nil
}
