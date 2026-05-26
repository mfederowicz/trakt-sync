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

// UsersAddListItemsHandler struct for handler
type UsersAddListItemsHandler struct{ common CommonLogic }

// Handle to handle users: add list items action
func (u UsersAddListItemsHandler) Handle(options *str.Options, client *internal.Client) error {
	input, err := u.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("add list items")
	toList := u.common.CreateItemsToAdd(input)
	addResult, resp, err := u.usersAddListItems(client, options, &toList)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("list:%s not found", options.ID)
	}

	if err != nil {
		return fmt.Errorf("add list items error:%w", err)
	}
	options.Output = "users_add_list_items_results.json"
	print("write result to:" + options.Output)
	jsonDataResult, _ := json.MarshalIndent(addResult, "", "  ")
	writer.WriteJSON(options, jsonDataResult)
	return nil
}

func (UsersAddListItemsHandler) usersAddListItems(client *internal.Client, options *str.Options, items *str.HistoryItems) (*str.AddResult, *str.Response, error) {
	user := options.UserName
	listID := options.ID
	result, resp, err := client.Users.AddListItems(
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
