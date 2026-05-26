// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersFollowersHandler struct for handler
type UsersFollowersHandler struct{ common CommonLogic }

// Handle to handle users: followers action
func (u UsersFollowersHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all followers including when the relationship began.")

	items, err := u.fetchFollowers(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get followers error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersFollowersHandler) fetchFollowers(client *internal.Client, options *str.Options, page int) ([]*str.Follower, error) {
	items, err := u.common.FetchFollowers(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch followers error:%w", err)
	}

	if len(items) == consts.ZeroValue {
		return nil, errors.New("empty list")
	}

	return items, nil
}
