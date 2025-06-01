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

// EpisodesWatchingHandler struct for handler
type EpisodesWatchingHandler struct{}

// Handle to handle episodes: watching action
func (m EpisodesWatchingHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all users watching this episode right now.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	result, _, err := m.fetchEpisodesWatching(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found watching for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (EpisodesWatchingHandler) fetchEpisodesWatching(client *internal.Client, options *str.Options) ([]*str.UserProfile, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetEpisodesWatching(
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
