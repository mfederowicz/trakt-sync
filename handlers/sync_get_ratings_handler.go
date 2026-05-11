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

// SyncGetRatingsHandler struct for handler
type SyncGetRatingsHandler struct{ common CommonLogic }

// Handle to handle sync: get_ratings action
func (m SyncGetRatingsHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Get user's ratings for type:", options.Type)

	items, err := m.syncGetRatings(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get ratings error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m SyncGetRatingsHandler) syncGetRatings(client *internal.Client, options *str.Options, page int) ([]*str.RatingListItem, error) {
	items, err := m.common.FetchRatings(client, options, page)

	if err != nil {
		return nil, err
	}

	return items, nil
}
