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

// SyncAddToCollectionHandler struct for handler
type SyncAddToCollectionHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_collection action
func (m SyncAddToCollectionHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(options.CollectionItems)
	if err != nil {
		return err
	}
	printer.Println("Add collection")
	result, err := m.syncAddToCollection(client, options, items)
	if err != nil {
		return fmt.Errorf("add to collection error:%w", err)
	}

	print("write result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (SyncAddToCollectionHandler) syncAddToCollection(client *internal.Client, options *str.Options, items *str.ItemsList) (*str.CollectionAddResult, error) {
	result, err := client.Sync.AddItemsToCollection(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
