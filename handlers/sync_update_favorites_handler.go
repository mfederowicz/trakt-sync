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

// SyncUpdateFavoritesHandler struct for handler
type SyncUpdateFavoritesHandler struct{ common CommonLogic }

// Handle to handle sync: update_favorites action
func (m SyncUpdateFavoritesHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Update the favorites by sending 1 or more parameters.")

	result, err := m.syncUpdateFavorites(client, options)
	if err != nil {
		return fmt.Errorf("update favorites error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (SyncUpdateFavoritesHandler) syncUpdateFavorites(client *internal.Client, options *str.Options) (*str.PersonalList, error) {
	update := new(str.PersonalList)
	update.Description = &options.Description
	update.SortBy = &options.SortBy
	update.SortHow = &options.SortHow

	result, err := client.Sync.UpdateFavorites(client.BuildCtxFromOptions(options), update)
	if err != nil {
		return nil, fmt.Errorf("update favorites error:%w", err)
	}

	return result, nil
}
