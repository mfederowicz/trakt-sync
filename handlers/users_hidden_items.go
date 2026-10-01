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

// UsersHiddenItemsHandler struct for handler
type UsersHiddenItemsHandler struct{ common CommonLogic }

// Handle to handle users: hidden_items action
func (u UsersHiddenItemsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Get hidden items")
	items, err := u.common.FetchUsersHiddenItems(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get hidden items error:%w", err)
	}

	jsonData, err := json.MarshalIndent(items, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}
