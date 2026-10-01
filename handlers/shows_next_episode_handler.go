// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsNextEpisodeHandler struct for handler
type ShowsNextEpisodeHandler struct{}

// Handle to handle shows: next_episode action
func (m ShowsNextEpisodeHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns the next scheduled to air episode.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, resp, err := m.fetchShowsNextEpisode(client, options)

	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNoContent {
		printer.Print("Not found next episode")
		return nil
	}

	printer.Printf("Found next episode for id:%s and name:%s \n", options.InternalID, episodeTitle(result))

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsNextEpisodeHandler) fetchShowsNextEpisode(client *trakt.Client, options *str.Options) (*str.Episode, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	show, resp, err := client.Shows.GetNextEpisode(
		cli.ContextFromOptions(options),
		options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return show, resp, nil
}
