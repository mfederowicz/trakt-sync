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

// UsersWatchlistHandler struct for handler
type UsersWatchlistHandler struct{ common CommonLogic }

// Handle to handle users: watchlist action
func (m UsersWatchlistHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Returns all items in a user's watchlist filtered by type:", options.Type)

	items, err := m.usersWatchlist(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get watchlist error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m UsersWatchlistHandler) usersWatchlist(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	items, err := m.common.FetchUsersWatchlist(client, options, page)

	if err != nil {
		return nil, err
	}

	return items, nil
}
