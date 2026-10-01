// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersRemoveHiddenItemsHandler struct for handler
type UsersRemoveHiddenItemsHandler struct{ common CommonLogic }

// Handle to handle users: remove_hidden_items action
func (u UsersRemoveHiddenItemsHandler) Handle(options *str.Options, client *trakt.Client) error {
	items, err := u.common.ReadInput(*options)
	if err != nil {
		return err
	}
	toHidden := u.common.CreateItemsToHidden(options.Section, items)
	addResult, err := u.common.UsersRemoveHiddenItems(client, options, &toHidden)
	if err != nil {
		return fmt.Errorf("remove hidden items error:%w", err)
	}

	jsonDataResult, err := json.MarshalIndent(addResult, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write result to:" + options.Output)
	writer.WriteJSON(options, jsonDataResult)
	return nil
}
