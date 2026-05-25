// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersFavoritesHandler struct for handler
type UsersFavoritesHandler struct{ common CommonLogic }

// Handle to handle users: favorites action
func (m UsersFavoritesHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Returns the top 100 shows and movies a user has favorited filtered by type:", options.Type)

	items, err := m.usersFavorites(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get favorites error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m UsersFavoritesHandler) usersFavorites(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	items, err := m.common.FetchUsersFavorites(client, options, page)

	if err != nil {
		return nil, err
	}

	return items, nil
}
