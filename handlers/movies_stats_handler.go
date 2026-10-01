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
	"github.com/mfederowicz/trakt-sync/writer"
)

// MoviesStatsHandler struct for handler
type MoviesStatsHandler struct{ common CommonLogic }

// Handle to handle movies: ratings action
func (m MoviesStatsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns lots of movie stats.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesStats(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found stats for id:%s\n", options.InternalID)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesStatsHandler) fetchMoviesStats(client *trakt.Client, options *str.Options) (*str.MovieStats, *str.Response, error) {
	result, resp, err := client.Movies.GetMovieStats(
		cli.ContextFromOptions(options),
		options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
