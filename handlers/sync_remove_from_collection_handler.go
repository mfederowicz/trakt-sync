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

// SyncRemoveFromCollectionHandler struct for handler
type SyncRemoveFromCollectionHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_collection action
func (m SyncRemoveFromCollectionHandler) Handle(options *str.Options, client *internal.Client) error {
	items, err := m.common.ReadInput(options.CollectionItems)
	if err != nil {
		return err
	}
	printer.Println("Remove from collection")
	result, err := m.syncRemoveFromCollection(client, options, items)
	if err != nil {
		return fmt.Errorf("remove from collection error:%w", err)
	}

	print("write result to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (SyncRemoveFromCollectionHandler) syncRemoveFromCollection(client *internal.Client, options *str.Options, items *str.ItemsList) (*str.CollectionRemoveResult, error) {
	result, err := client.Sync.RemoveItemsFromCollection(
		client.BuildCtxFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
