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

// UsersFollowingHandler struct for handler
type UsersFollowingHandler struct{ common CommonLogic }

// Handle to handle users: following action
func (u UsersFollowingHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all user's they follow including when the relationship began.")

	items, err := u.fetchFollowing(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get followers error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersFollowingHandler) fetchFollowing(client *internal.Client, options *str.Options, page int) ([]*str.Follower, error) {
	items, err := u.common.FetchFollowing(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch following error:%w", err)
	}

	if len(items) == consts.ZeroValue {
		return nil, errors.New("empty list")
	}

	return items, nil
}
