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

// SeasonsStatsHandler struct for handler
type SeasonsStatsHandler struct{ common CommonLogic }

// Handle to handle seasons: stats action
func (m SeasonsStatsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns lots of season stats.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsStats(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found season stats for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsStatsHandler) fetchSeasonsStats(client *internal.Client, options *str.Options) (*str.SeasonStats, *str.Response, error) {
	result, resp, err := client.Shows.GetSeasonStats(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Season,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
