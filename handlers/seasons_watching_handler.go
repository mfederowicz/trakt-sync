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

// SeasonsWatchingHandler struct for handler
type SeasonsWatchingHandler struct{}

// Handle to handle seasons: watching action
func (m SeasonsWatchingHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all users watching this season right now.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsWatching(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found watching for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsWatchingHandler) fetchSeasonsWatching(client *internal.Client, options *str.Options) ([]*str.UserProfile, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetSeasonsWatching(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Season,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
