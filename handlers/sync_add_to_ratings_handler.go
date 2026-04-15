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

// SyncAddToRatingsHandler struct for handler
type SyncAddToRatingsHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_ratings action
func (m SyncAddToRatingsHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("add to ratings")
	toRatings := m.common.CreateItemsToAddRatings(items)
	addResult, err := m.syncAddToRatings(client, options, &toRatings)
	if err != nil {
		return fmt.Errorf("add to ratings error:%w", err)
	}
	
	options.Output = "export_sync_add_to_ratings_results.json"
	
	print("write result to:" + options.Output)
	jsonDataResult, _ := json.MarshalIndent(addResult, "", "  ")
	writer.WriteJSON(options, jsonDataResult)
	return nil
}

func (SyncAddToRatingsHandler) syncAddToRatings(client *internal.Client, options *str.Options, items *str.RatingItems) (*str.AddResult, error) {
	result, err := client.Sync.AddItemsToRatings(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
