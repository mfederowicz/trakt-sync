// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersFollowersHandler struct for handler
type UsersFollowersHandler struct{ common CommonLogic }

// Handle to handle users: followers action
func (u UsersFollowersHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all followers including when the relationship began.")

	items, err := u.fetchFollowers(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get followers error:%w", err)
	}
	jsonData, err := json.MarshalIndent(items, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersFollowersHandler) fetchFollowers(client *trakt.Client, options *str.Options, page int) ([]*str.Follower, error) {
	items, err := u.common.FetchFollowers(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch followers error:%w", err)
	}

	if len(items) == consts.ZeroValue {
		return nil, errors.New("empty list")
	}

	return items, nil
}
