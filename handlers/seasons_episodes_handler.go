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

// SeasonsEpisodesHandler struct for handler
type SeasonsEpisodesHandler struct{}

// Handle to handle seasons: episodes action
func (m SeasonsEpisodesHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all episodes for a specific season of a show.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsEpisodes(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found episodes for id:%s season:%d\n", options.InternalID, options.Season)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsEpisodesHandler) fetchSeasonsEpisodes(client *trakt.Client, options *str.Options) ([]*str.Episode, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo, Translations: options.Translations}
	result, resp, err := client.Shows.GetAllEpisodesForSingleSeason(
		cli.ContextFromOptions(options),
		&options.InternalID,
		&options.Season,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
