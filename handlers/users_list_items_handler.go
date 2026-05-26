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

// UsersListItemsHandler struct for handler
type UsersListItemsHandler struct{ common CommonLogic }

// Handle to handle users: list_items action
func (u UsersListItemsHandler) Handle(options *str.Options, client *internal.Client) error {
	err := u.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Get all items on a personal list.")
	items, err := u.fetchListItems(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get list items error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersListItemsHandler) fetchListItems(client *internal.Client, options *str.Options, page int) ([]*str.UserListItem, error) {
	items, err := u.common.FetchUsersListItems(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch list items error:%w", err)
	}

	if len(items) == consts.ZeroValue {
		return nil, errors.New("empty list items")
	}

	return items, nil
}
