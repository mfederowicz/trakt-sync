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

// EpisodesStatsHandler struct for handler
type EpisodesStatsHandler struct{ common CommonLogic }

// Handle to handle episode: stats action
func (m EpisodesStatsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns lots of episodes stats.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	result, _, err := m.fetchEpisodesStats(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found episodes stats for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (EpisodesStatsHandler) fetchEpisodesStats(client *internal.Client, options *str.Options) (*str.EpisodeStats, *str.Response, error) {
	result, resp, err := client.Shows.GetEpisodeStats(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Season,
		&options.Episode,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
