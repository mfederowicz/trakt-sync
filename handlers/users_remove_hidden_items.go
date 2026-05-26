// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersRemoveHiddenItemsHandler struct for handler
type UsersRemoveHiddenItemsHandler struct{ common CommonLogic }

// Handle to handle users: remove_hidden_items action
func (u UsersRemoveHiddenItemsHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := u.common.ReadInput(*options)
	if err != nil {
		return err
	}
	toHidden := u.common.CreateItemsToHidden(options.Section, items)
	addResult, err := u.common.UsersRemoveHiddenItems(client, options, &toHidden)
	if err != nil {
		return fmt.Errorf("remove hidden items error:%w", err)
	}

	options.Output = "users_remove_hidden_items_results.json"

	print("write result to:" + options.Output)
	jsonDataResult, _ := json.MarshalIndent(addResult, "", "  ")
	writer.WriteJSON(options, jsonDataResult)
	return nil
}
