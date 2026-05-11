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

// SyncReorderFavoritesHandler struct for handler
type SyncReorderFavoritesHandler struct{ common CommonLogic }

// Handle to handle sync: reorder_favorites action
func (m SyncReorderFavoritesHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("reorder favorites")
	toReorder := m.common.CreateItemsToReorder(items)
	result, err := m.syncReorderFavorites(client, options, &toReorder)
	if err != nil {
		return fmt.Errorf("reorder favorites error:%w", err)
	}
	options.Output = "sync_reorder_favorites_results.json"
	printer.Println("write result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SyncReorderFavoritesHandler) syncReorderFavorites(client *internal.Client, options *str.Options, items *str.ItemsToReorder) (*str.ReorderResults, error) {
	result, err := client.Sync.ReorderFavoritesItems(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
