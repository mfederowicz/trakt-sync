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
	"github.com/mfederowicz/trakt-sync/writer"
)

// EpisodesRatingsHandler struct for handler
type EpisodesRatingsHandler struct{ common CommonLogic }

// Handle to handle episodes: ratings action
func (m EpisodesRatingsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns rating (between 0 and 10) and distribution for a episode.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	result, _, err := m.fetchEpisodesRatings(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found ratings for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (EpisodesRatingsHandler) fetchEpisodesRatings(client *trakt.Client, options *str.Options) (*str.EpisodeRatings, *str.Response, error) {
	result, resp, err := client.Shows.GetEpisodeRatings(
		cli.ContextFromOptions(options),
		&options.InternalID,
		&options.Season,
		&options.Episode,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
