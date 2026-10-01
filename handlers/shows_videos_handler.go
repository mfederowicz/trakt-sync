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
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsVideosHandler struct for handler
type ShowsVideosHandler struct{}

// Handle to handle shows: videos action
func (m ShowsVideosHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all videos including trailers, teasers, clips, and featurettes.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsVideos(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found videos for id:%s\n", options.InternalID)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsVideosHandler) fetchShowsVideos(client *trakt.Client, options *str.Options) ([]*str.Video, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Shows.GetShowVideos(
		cli.ContextFromOptions(options),
		options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
