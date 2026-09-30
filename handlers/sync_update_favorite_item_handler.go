// Package handlers used to handle module actions
package handlers

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// SyncUpdateFavoriteItemHandler struct for handler
type SyncUpdateFavoriteItemHandler struct{ common CommonLogic }

// Handle to handle sync: update_favorite_item action
func (m SyncUpdateFavoriteItemHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Update the favorite item note.")

	err := m.syncUpdateFavoriteItem(client, options)
	if err != nil {
		return fmt.Errorf("update favorite item error:%w", err)
	}
	fmt.Println("update notes success for favorite item:", options.ListItemID)

	return nil
}

func (SyncUpdateFavoriteItemHandler) syncUpdateFavoriteItem(client *trakt.Client, options *str.Options) error {
	update := new(str.FavoriteItem)
	update.Notes = &options.Notes

	err := client.Sync.UpdateFavoriteItem(client.BuildCtxFromOptions(options), options.ListItemID, update)
	if err != nil {
		return fmt.Errorf("update favorite item error:%w", err)
	}

	return nil
}
