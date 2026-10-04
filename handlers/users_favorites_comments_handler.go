// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersFavoritesCommentsHandler struct for handler
type UsersFavoritesCommentsHandler struct{ common CommonLogic }

// Handle to handle users: favorites_comments action
func (m UsersFavoritesCommentsHandler) Handle(options *str.Options, client *trakt.Client) error {
	err := m.common.CheckSortAndTypes(options)
	if err != nil {
		return err
	}
	printer.Println("Returns all top level comments for the favorites. sorted by:", options.Sort)
	items, err := m.usersFavoritesComments(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get favorites comments error:%w", err)
	}
	jsonData, err := json.MarshalIndent(items, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m UsersFavoritesCommentsHandler) usersFavoritesComments(client *trakt.Client, options *str.Options, page int) ([]*str.ListComment, error) {
	items, err := m.common.FetchUsersFavoritesComments(client, options, page)

	if err != nil {
		return nil, err
	}

	return items, nil
}
