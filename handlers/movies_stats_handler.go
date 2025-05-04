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
	"github.com/mfederowicz/trakt-sync/writer"
)

// MoviesStatsHandler struct for handler
type MoviesStatsHandler struct{ common CommonLogic }

// Handle to handle movies: ratings action
func (m MoviesStatsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns lots of movie stats.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesStats(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found stats for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesStatsHandler) fetchMoviesStats(client *internal.Client, options *str.Options) (*str.MovieStats, *str.Response, error) {
	result, resp, err := client.Movies.GetMovieStats(
		context.Background(),
		&options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
