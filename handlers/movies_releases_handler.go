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

// MoviesReleasesHandler struct for handler
type MoviesReleasesHandler struct{}

// Handle to handle movies: releases action
func (m MoviesReleasesHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all releases for a movie")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesReleases(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found releases for id:%s\n", options.InternalID)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesReleasesHandler) fetchMoviesReleases(client *trakt.Client, options *str.Options) ([]*str.Release, *str.Response, error) {
	releases, resp, err := client.Movies.GetAllMovieReleases(
		cli.ContextFromOptions(options),
		options.InternalID,
		options.Country,
	)

	if err != nil {
		return nil, nil, err
	}

	return releases, resp, nil
}
