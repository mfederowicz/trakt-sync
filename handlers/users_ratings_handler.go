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

// UsersRatingsHandler struct for handler
type UsersRatingsHandler struct{ common CommonLogic }

// Handle to handle users: ratings action
func (m UsersRatingsHandler) Handle(options *str.Options, client *trakt.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Get user's ratings for type:", options.Type)

	items, err := m.usersRatings(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get ratings error:%w", err)
	}
	jsonData, err := json.MarshalIndent(items, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m UsersRatingsHandler) usersRatings(client *trakt.Client, options *str.Options, page int) ([]*str.RatingListItem, error) {
	items, err := m.common.FetchUsersRatings(client, options, page)

	if err != nil {
		return nil, err
	}

	return items, nil
}
