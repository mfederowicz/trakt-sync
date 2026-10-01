// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncAddToRatingsHandler struct for handler
type SyncAddToRatingsHandler struct{ common CommonLogic }

// Handle to handle sync: add_to_ratings action
func (m SyncAddToRatingsHandler) Handle(options *str.Options, client *trakt.Client) error {
	items, err := m.common.ReadInput(*options)
	if err != nil {
		return err
	}
	printer.Println("add to ratings")
	toRatings := m.common.CreateItemsToAddRatings(items)
	addResult, err := m.syncAddToRatings(client, options, &toRatings)
	if err != nil {
		return fmt.Errorf("add to ratings error:%w", err)
	}

	jsonDataResult, err := json.MarshalIndent(addResult, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write result to:" + options.Output)
	writer.WriteJSON(options, jsonDataResult)
	return nil
}

func (SyncAddToRatingsHandler) syncAddToRatings(client *trakt.Client, options *str.Options, items *str.RatingItems) (*str.AddResult, error) {
	result, _, err := client.Sync.AddItemsToRatings(
		cli.ContextFromOptions(options),
		items,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
