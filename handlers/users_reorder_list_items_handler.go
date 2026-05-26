// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersReorderListItemsHandler struct for handler
type UsersReorderListItemsHandler struct{ common CommonLogic }

// Handle to handle users: reorder_list_items action
func (m UsersReorderListItemsHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	toReorder := m.common.CreateItemsToReorder(items)
	result, err := m.usersReorderListItems(client, options, &toReorder)
	if err != nil {
		return fmt.Errorf("reorder list items error:%w", err)
	}
	options.Output = "users_reorder_list_items_results.json"
	printer.Println("write result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (UsersReorderListItemsHandler) usersReorderListItems(client *internal.Client, options *str.Options, items *str.ItemsToReorder) (*str.ReorderResults, error) {
	user := options.UserName
	listID := options.ID
	result, _, err := client.Users.ReorderListItems(
		client.BuildCtxFromOptions(options),
		&user,
		&listID,
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
