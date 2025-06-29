// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncGetCollectionHandler struct for handler
type SyncGetCollectionHandler struct{ common CommonLogic }

// Handle to handle sync: get_collection action
func (m SyncGetCollectionHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get collection type:", options.Type)
	items, err := m.syncGetCollectionItems(client, options)
	if err != nil {
		return fmt.Errorf("get collection error:%w", err)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (SyncGetCollectionHandler) syncGetCollectionItems(client *internal.Client, options *str.Options) ([]*str.ExportlistItem, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}

	if options.Type == consts.Seasons {
		items, _, err := client.Sync.GetCollectedSeasons(
			client.BuildCtxFromOptions(options),
			&opts,
		)
		if err != nil {
			return nil, err
		}

		return items, nil
	}

	items, _, err := client.Sync.GetCollection(
		client.BuildCtxFromOptions(options),
		&options.Type,
		&opts,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}
