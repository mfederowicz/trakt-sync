// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncGetFavoritesHandler struct for handler
type SyncGetFavoritesHandler struct{ common CommonLogic }

// Handle to handle sync: get_favorites action
func (m SyncGetFavoritesHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Returns all items in a user's favorites filtered by type:", options.Type)

	items, err := m.syncGetFavorites(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get favorites error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m SyncGetFavoritesHandler) syncGetFavorites(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	favorites, err := m.common.FetchFavorites(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch favorites error:%w", err)
	}

	if len(favorites) == consts.ZeroValue {
		return nil, errors.New("empty favorites")
	}

	return favorites, nil
}
