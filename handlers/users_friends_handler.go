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

// UsersFriendsHandler struct for handler
type UsersFriendsHandler struct{ common CommonLogic }

// Handle to handle users: friends action
func (u UsersFriendsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all friends for a user including when the relationship began.")

	items, err := u.fetchFriends(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get friends error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersFriendsHandler) fetchFriends(client *internal.Client, options *str.Options, page int) ([]*str.Friend, error) {
	items, err := u.common.FetchFriends(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch friends error:%w", err)
	}

	if len(items) == consts.ZeroValue {
		return nil, errors.New("empty list")
	}

	return items, nil
}
