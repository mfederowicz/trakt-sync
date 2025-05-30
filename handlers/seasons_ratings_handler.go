// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SeasonsRatingsHandler struct for handler
type SeasonsRatingsHandler struct{ common CommonLogic }

// Handle to handle seasons: ratings action
func (m SeasonsRatingsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns rating (between 0 and 10) and distribution for a season.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsRatings(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found ratings for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsRatingsHandler) fetchSeasonsRatings(client *internal.Client, options *str.Options) (*str.SeasonRatings, *str.Response, error) {
	result, resp, err := client.Shows.GetSeasonRatings(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Season,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
