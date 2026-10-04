// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersWatchlistHandler struct for handler
type UsersWatchlistHandler struct{ common CommonLogic }

// Handle to handle users: watchlist action
func (m UsersWatchlistHandler) Handle(options *str.Options, client *trakt.Client) error {
	if err := checkMediaFilters(options); err != nil {
		return err
	}
	if len(options.HideItems) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.WatchlistHideFilters, options.HideItems) {
		return fmt.Errorf("hide '%s' is not valid, available values: %v", options.HideItems, cfg.WatchlistHideFilters)
	}
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}
	// -sort switches to the /{type}/{sort} route; check it before any request
	if len(options.SortPath) > consts.ZeroValue {
		if _, err := sortRouteType(consts.Watchlist, options); err != nil {
			return err
		}
	}

	printer.Println("Returns all items in a user's watchlist filtered by type:", options.Type)

	items, err := m.usersWatchlist(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get watchlist error:%w", err)
	}
	if len(items) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}
	jsonData, err := json.MarshalIndent(items, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m UsersWatchlistHandler) usersWatchlist(client *trakt.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	items, err := m.common.FetchUsersWatchlist(client, options, page)

	if err != nil {
		return nil, err
	}

	return items, nil
}
