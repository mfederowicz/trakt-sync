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

// SyncReorderFavoritesHandler struct for handler
type SyncReorderFavoritesHandler struct{ common CommonLogic }

// Handle to handle sync: reorder_favorites action
func (m SyncReorderFavoritesHandler) Handle(options *str.Options, client *trakt.Client) error {
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
	printer.Println("write result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SyncReorderFavoritesHandler) syncReorderFavorites(client *trakt.Client, options *str.Options, items *str.ItemsToReorder) (*str.ReorderResults, error) {
	result, _, err := client.Sync.ReorderFavoritesItems(
		cli.ContextFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
