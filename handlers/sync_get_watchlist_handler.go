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

// SyncGetWatchlistHandler struct for handler
type SyncGetWatchlistHandler struct{ common CommonLogic }

// Handle to handle sync: get_watchlist action
func (m SyncGetWatchlistHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}
	

	printer.Println("Returns all items in a user's watchlist filtered by type:", options.Type)

	items, err := m.syncGetWatchlist(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get ratings error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m SyncGetWatchlistHandler) syncGetWatchlist(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	watchlist, err := m.common.FetchWatchlist(client, options, consts.DefaultPage)
	if err != nil {
		return nil, fmt.Errorf("fetch watchlist error:%w", err)
	}

	if len(watchlist) == consts.ZeroValue {
		return nil, errors.New("empty watchlist")
	}

	return watchlist, nil
}
