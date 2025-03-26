// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// MoviesVideosHandler struct for handler
type MoviesVideosHandler struct{}

// Handle to handle movies: videos action
func (m MoviesVideosHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all videos including trailers, teasers, clips, and featurettes.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesVideos(client, options)
	
	if err != nil {
		return err
	}	

	printer.Printf("Found videos for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesVideosHandler) fetchMoviesVideos(client *internal.Client, options *str.Options) ([]*str.Video, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Movies.GetMovieVideos(
		context.Background(),
		&options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}

