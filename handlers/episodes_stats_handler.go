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

// EpisodesStatsHandler struct for handler
type EpisodesStatsHandler struct{ common CommonLogic }

// Handle to handle episode: stats action
func (m EpisodesStatsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns lots of episodes stats.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	result, _, err := m.fetchEpisodesStats(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found episodes stats for id:%s\n", options.InternalID)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (EpisodesStatsHandler) fetchEpisodesStats(client *trakt.Client, options *str.Options) (*str.EpisodeStats, *str.Response, error) {
	result, resp, err := client.Shows.GetEpisodeStats(
		cli.ContextFromOptions(options),
		options.InternalID,
		options.Season,
		options.Episode,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
