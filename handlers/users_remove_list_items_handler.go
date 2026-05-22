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

// UsersRemoveListItemsHandler struct for handler
type UsersRemoveListItemsHandler struct{ common CommonLogic }

// Handle to handle users: remove list items action
func (u UsersRemoveListItemsHandler) Handle(options *str.Options, client *internal.Client) error {
	input, err := u.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("remove list items")
	toList := u.common.CreateItemsToAdd(input)
	removeResult, resp, err := u.usersRemoveListItems(client, options, &toList)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("list:%s not found", options.ID)
	}

	if err != nil {
		return fmt.Errorf("remove list items error:%w", err)
	}
	options.Output = "users_remove_list_items_results.json"
	print("write result to:" + options.Output)
	jsonDataResult, _ := json.MarshalIndent(removeResult, "", "  ")
	writer.WriteJSON(options, jsonDataResult)
	return nil
}

func (UsersRemoveListItemsHandler) usersRemoveListItems(client *internal.Client, options *str.Options, items *str.HistoryItems) (*str.RemoveResult, *str.Response, error) {
	user := options.UserName
	listID := options.ID
	result, resp, err := client.Users.RemoveListItems(
		client.BuildCtxFromOptions(options),
		&user,
		&listID,
		items,
	)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
