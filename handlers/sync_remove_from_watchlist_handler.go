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

// SyncRemoveFromWatchlistHandler struct for handler
type SyncRemoveFromWatchlistHandler struct{ common CommonLogic }

// Handle to handle sync: remove_from_watchlist action
func (m SyncRemoveFromWatchlistHandler) Handle(options *str.Options, client *trakt.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("clean watchlist")
	toRemove := m.common.CreateItemsToRemove(items)
	result, err := m.syncRemoveFromWatchlist(client, options, &toRemove)
	if err != nil {
		return fmt.Errorf("clean watchlist error:%w", err)
	}

	printer.Println("write cleanup result to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SyncRemoveFromWatchlistHandler) syncRemoveFromWatchlist(client *trakt.Client, options *str.Options, items *str.ItemsToRemove) (*str.RemoveResult, error) {
	result, _, err := client.Sync.RemoveItemsFromWatchlist(
		cli.ContextFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
