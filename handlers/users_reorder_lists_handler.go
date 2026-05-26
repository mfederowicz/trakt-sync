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

// UsersReorderListsHandler struct for handler
type UsersReorderListsHandler struct{ common CommonLogic }

// Handle to handle users: reorder_lists action
func (m UsersReorderListsHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	toReorder := m.common.CreateItemsToReorder(items)
	result, err := m.usersReorderLists(client, options, &toReorder)
	if err != nil {
		return fmt.Errorf("reorder lists error:%w", err)
	}
	options.Output = "users_reorder_lists_results.json"
	printer.Println("write result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (UsersReorderListsHandler) usersReorderLists(client *internal.Client, options *str.Options, items *str.ItemsToReorder) (*str.ReorderResults, error) {
	user := options.UserName
	result, _, err := client.Users.ReorderLists(
		client.BuildCtxFromOptions(options),
		&user,
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
