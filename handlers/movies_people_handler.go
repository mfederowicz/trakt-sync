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

// MoviesPeopleHandler struct for handler
type MoviesPeopleHandler struct{ common CommonLogic }

// Handle to handle movies: people action
func (m MoviesPeopleHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all people related with movie.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesPeople(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found people for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesPeopleHandler) fetchMoviesPeople(client *internal.Client, options *str.Options) (*str.MoviePeople, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Movies.GetAllPeopleForMovie(
		context.Background(),
		&options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	return result, resp, nil
}
