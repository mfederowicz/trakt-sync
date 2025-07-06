// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncGetWatchedHandler struct for handler
type SyncGetWatchedHandler struct{ common CommonLogic }

// Handle to handle sync: get_watched action
func (m SyncGetWatchedHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Get watched type:", options.Type)
	items, err := m.syncGetWatchedItems(client, options)
	if err != nil {
		return fmt.Errorf("get watched error:%w", err)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (SyncGetWatchedHandler) syncGetWatchedItems(client *internal.Client, options *str.Options) ([]*str.UserWatched, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	items, _, err := client.Sync.GetWatched(
		client.BuildCtxFromOptions(options),
		&options.Type,
		&opts,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}
