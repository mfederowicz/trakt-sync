// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SeasonsStatsHandler struct for handler
type SeasonsStatsHandler struct{}

// Handle to handle seasons: stats action
func (m SeasonsStatsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns lots of season stats.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsStats(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found season stats for id:%s\n", options.InternalID)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsStatsHandler) fetchSeasonsStats(client *trakt.Client, options *str.Options) (*str.SeasonStats, *str.Response, error) {
	result, resp, err := client.Shows.GetSeasonStats(
		cli.ContextFromOptions(options),
		options.InternalID,
		options.Season,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
