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

// SyncRemoveFromRatingsHandler struct for handler
type SyncRemoveFromRatingsHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_ratings action
func (m SyncRemoveFromRatingsHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("clean ratings")
	toRemove := m.common.CreateItemsToRemoveRatings(items)
	result, err := m.syncRemoveFromRatings(client, options, &toRemove)
	if err != nil {
		return fmt.Errorf("clean ratings error:%w", err)
	}

	options.Output = "sync_remove_from_ratings_results.json"

	printer.Println("write cleanup result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SyncRemoveFromRatingsHandler) syncRemoveFromRatings(client *internal.Client, options *str.Options, items *str.ItemsToRemove) (*str.RemoveResult, error) {
	result, err := client.Sync.RemoveItemsFromRatings(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
