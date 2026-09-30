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

// EpisodesVideosHandler struct for handler
type EpisodesVideosHandler struct{}

// Handle to handle episodes: videos action
func (m EpisodesVideosHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all videos including trailers, teasers, clips, and featurettes.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	result, _, err := m.fetchEpisodesVideos(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found episode videos for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (EpisodesVideosHandler) fetchEpisodesVideos(client *trakt.Client, options *str.Options) ([]*str.Video, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetEpisodeVideos(
		cli.ContextFromOptions(options),
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
