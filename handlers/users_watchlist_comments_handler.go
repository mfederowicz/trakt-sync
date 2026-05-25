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

// UsersWatchlistCommentsHandler struct for handler
type UsersWatchlistCommentsHandler struct{ common CommonLogic }

// Handle to handle users: watchlist_comments action
func (m UsersWatchlistCommentsHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckSortAndTypes(options)
	if err != nil {
		return err
	}
	printer.Println("Returns all top level comments for the watchlist. sorted by:", options.Sort)
	items, err := m.usersWatchlistComments(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get watchlist error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m UsersWatchlistCommentsHandler) usersWatchlistComments(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	items, err := m.common.FetchUsersWatchlistComments(client, options, page)

	if err != nil {
		return nil, err
	}

	return items, nil
}
