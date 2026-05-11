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

// SyncAddToFavoritesHandler struct for handler
type SyncAddToFavoritesHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_favorites action
func (m SyncAddToFavoritesHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("add to favorites")
	toFavorites := m.common.CreateItemsToAdd(items)
	addResult, err := m.syncAddToFavorites(client, options, &toFavorites)
	if err != nil {
		return fmt.Errorf("add to favorites error:%w", err)
	}

	options.Output = "sync_add_to_favorites_results.json"

	print("write result to:" + options.Output)
	jsonDataResult, _ := json.MarshalIndent(addResult, "", "  ")
	writer.WriteJSON(options, jsonDataResult)
	return nil
}

func (SyncAddToFavoritesHandler) syncAddToFavorites(client *internal.Client, options *str.Options, items *str.HistoryItems) (*str.AddResult, error) {
	result, err := client.Sync.AddItemsToFavorites(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
