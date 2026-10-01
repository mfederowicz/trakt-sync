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
	"github.com/mfederowicz/trakt-sync/writer"
)

// MoviesStudiosHandler struct for handler
type MoviesStudiosHandler struct{}

// Handle to handle movies: studios action
func (m MoviesStudiosHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all studios for a movie")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesStudios(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found studios for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesStudiosHandler) fetchMoviesStudios(client *trakt.Client, options *str.Options) ([]*str.Studio, *str.Response, error) {
	result, resp, err := client.Movies.GetMovieStudios(
		cli.ContextFromOptions(options),
		options.InternalID,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
