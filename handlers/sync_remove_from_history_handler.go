// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncRemoveFromHistoryHandler struct for handler
type SyncRemoveFromHistoryHandler struct{ common CommonLogic }

// Handle to handle sync: remove_from_history action
func (m SyncRemoveFromHistoryHandler) Handle(options *str.Options, client *trakt.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("clean history")
	toRemove := m.common.CreateItemsToRemove(items)

	result, err := m.syncRemoveFromHistory(client, options, &toRemove)
	if err != nil {
		return fmt.Errorf("clean history error:%w", err)
	}

	printer.Println("write cleanup result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SyncRemoveFromHistoryHandler) syncRemoveFromHistory(client *trakt.Client, options *str.Options, items *str.ItemsToRemove) (*str.RemoveResult, error) {
	result, _, err := client.Sync.RemoveItemsFromHistory(
		cli.ContextFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
