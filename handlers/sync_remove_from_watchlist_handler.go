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

// SyncRemoveFromWatchlistHandler struct for handler
type SyncRemoveFromWatchlistHandler struct{ common CommonLogic }

// Handle to handle sync: remove_from_watchlist action
func (m SyncRemoveFromWatchlistHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("clean watchlist")
	toRemove := m.common.CreateItemsToRemove(items)
	result, err := m.syncRemoveFromWatchlist(client, options, &toRemove)
	if err != nil {
		return fmt.Errorf("clean watchlist error:%w", err)
	}

	options.Output = "sync_remove_from_watchlist_results.json"

	printer.Println("write cleanup result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SyncRemoveFromWatchlistHandler) syncRemoveFromWatchlist(client *internal.Client, options *str.Options, items *str.ItemsToRemove) (*str.RemoveResult, error) {
	result, err := client.Sync.RemoveItemsFromWatchlist(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
