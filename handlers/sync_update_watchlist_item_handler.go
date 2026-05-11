// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// SyncUpdateWatchlistItemHandler struct for handler
type SyncUpdateWatchlistItemHandler struct{ common CommonLogic }

// Handle to handle sync: update_watchlist_item action
func (m SyncUpdateWatchlistItemHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Update the watchlist item note.")

	err := m.syncUpdateWatchlistItem(client, options)
	if err != nil {
		return fmt.Errorf("update watchlist item error:%w", err)
	}
	fmt.Println("update notes success for watchlist item:", options.ListItemID)

	return nil
}

func (SyncUpdateWatchlistItemHandler) syncUpdateWatchlistItem(client *internal.Client, options *str.Options) error {
	update := new(str.WatchlistItem)
	update.Notes = &options.Notes

	err := client.Sync.UpdateWatchlistItem(client.BuildCtxFromOptions(options), options.ListItemID, update)
	if err != nil {
		return fmt.Errorf("update watchlist item error:%w", err)
	}

	return nil
}
