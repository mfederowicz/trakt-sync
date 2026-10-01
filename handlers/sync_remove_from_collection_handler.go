// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncRemoveFromCollectionHandler struct for handler
type SyncRemoveFromCollectionHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_collection action
func (m SyncRemoveFromCollectionHandler) Handle(options *str.Options, client *trakt.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("Remove from collection")
	result, err := m.syncRemoveFromCollection(client, options, items)
	if err != nil {
		return fmt.Errorf("remove from collection error:%w", err)
	}

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write result to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (SyncRemoveFromCollectionHandler) syncRemoveFromCollection(client *trakt.Client, options *str.Options, items *str.ItemsList) (*str.CollectionRemoveResult, error) {
	result, _, err := client.Sync.RemoveItemsFromCollection(
		cli.ContextFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
