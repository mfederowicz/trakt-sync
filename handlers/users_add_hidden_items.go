// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersAddHiddenItemsHandler struct for handler
type UsersAddHiddenItemsHandler struct{ common CommonLogic }

// Handle to handle users: add_hidden_items action
func (u UsersAddHiddenItemsHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := u.common.ReadInput(*options)
	if err != nil {
		return err
	}
	toHidden := u.common.CreateItemsToHidden(options.Section, items)
	addResult, err := u.common.UsersAddToHiddenItems(client, options, &toHidden)
	if err != nil {
		return fmt.Errorf("add hidden items error:%w", err)
	}

	options.Output = "users_add_hidden_items_results.json"

	print("write result to:" + options.Output)
	jsonDataResult, _ := json.MarshalIndent(addResult, "", "  ")
	writer.WriteJSON(options, jsonDataResult)
	return nil
}
