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

// SyncReorderWatchlistHandler struct for handler
type SyncReorderWatchlistHandler struct{ common CommonLogic }

// Handle to handle sync: reorder_watchlist action
func (m SyncReorderWatchlistHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("reorder watchlist")
	toReorder := m.common.CreateItemsToReorder(items)
	result, err := m.syncReorderWatchlist(client, options, &toReorder)
	if err != nil {
		return fmt.Errorf("reorder watchlist error:%w", err)
	}
	options.Output = "sync_reorder_watchlist_results.json"
	printer.Println("write result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SyncReorderWatchlistHandler) syncReorderWatchlist(client *internal.Client, options *str.Options, items *str.ItemsToReorder) (*str.ReorderResults, error) {
	result, err := client.Sync.ReorderWatchlistItems(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
