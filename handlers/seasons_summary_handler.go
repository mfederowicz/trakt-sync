// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SeasonsSummaryHandler struct for handler
type SeasonsSummaryHandler struct{}

// Handle to handle seasons: summary action
func (m SeasonsSummaryHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all seasons for a show including the number of episodes in each season.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsSummary(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found seasons for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsSummaryHandler) fetchSeasonsSummary(client *trakt.Client, options *str.Options) ([]*str.Season, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetAllSeasonsForShow(
		cli.ContextFromOptions(options),
		options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
