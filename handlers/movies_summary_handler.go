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

// MoviesSummaryHandler struct for handler
type MoviesSummaryHandler struct{}

// Handle to handle movies: summary action
func (m MoviesSummaryHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns a single movie details")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesSummary(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found movie for id:%s and name:%s \n", options.InternalID, *result.Title)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesSummaryHandler) fetchMoviesSummary(client *trakt.Client, options *str.Options) (*str.Movie, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	movie, resp, err := client.Movies.GetMovie(
		cli.ContextFromOptions(options),
		options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return movie, resp, nil
}
