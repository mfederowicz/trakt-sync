// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// EpisodesSummaryHandler struct for handler
type EpisodesSummaryHandler struct{}

// Handle to handle episodes: summary action
func (m EpisodesSummaryHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns a single episode's details.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	result, _, err := m.fetchEpisodesSummary(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found episodes for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (EpisodesSummaryHandler) fetchEpisodesSummary(client *internal.Client, options *str.Options) (*str.Episode, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetSingleEpisodeForShow(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Season,
		&options.Episode,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
