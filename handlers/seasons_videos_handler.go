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

// SeasonsVideosHandler struct for handler
type SeasonsVideosHandler struct{}

// Handle to handle seasons: videos action
func (m SeasonsVideosHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all videos including trailers, teasers, clips, and featurettes.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptySeasonIDMsg)
	}

	result, _, err := m.fetchSeasonsVideos(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found season videos for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (SeasonsVideosHandler) fetchSeasonsVideos(client *internal.Client, options *str.Options) ([]*str.Video, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetSeasonsVideos(
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
