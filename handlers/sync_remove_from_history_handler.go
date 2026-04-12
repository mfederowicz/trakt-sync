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

// SyncRemoveFromHistoryHandler struct for handler
type SyncRemoveFromHistoryHandler struct{ common CommonLogic }

// Handle to handle sync: remove_from_history action
func (m SyncRemoveFromHistoryHandler) Handle(options *str.Options, client *internal.Client) error {
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

	options.Output = "export_sync_remove_from_history_results.json"

	printer.Println("write cleanup result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)	
	return nil
}

func (SyncRemoveFromHistoryHandler) syncRemoveFromHistory(client *internal.Client, options *str.Options, items *str.ItemsToRemove) (*str.HistoryRemoveResult, error) {
	result, err := client.Sync.RemoveItemsFromHistory(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}


