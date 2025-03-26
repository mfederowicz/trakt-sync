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

// MoviesSummaryHandler struct for handler
type MoviesSummaryHandler struct{}

// Handle to handle movies: summary action
func (m MoviesSummaryHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns a single movie details")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesSummary(client, options)
	
	if err != nil {
		return err
	}	

	printer.Printf("Found movie for id:%s and name:%s \n", options.InternalID, *result.Title)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesSummaryHandler) fetchMoviesSummary(client *internal.Client, options *str.Options) (*str.Movie, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	movie, resp, err := client.Movies.GetMovie(
		context.Background(),
		&options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return movie, resp, nil
}

