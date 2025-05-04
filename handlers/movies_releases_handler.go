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

// MoviesReleasesHandler struct for handler
type MoviesReleasesHandler struct{}

// Handle to handle movies: releases action
func (m MoviesReleasesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns all releases for a movie")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesReleases(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found releases for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (MoviesReleasesHandler) fetchMoviesReleases(client *internal.Client, options *str.Options) ([]*str.Release, *str.Response, error) {
	releases, resp, err := client.Movies.GetAllMovieReleases(
		context.Background(),
		&options.InternalID,
		&options.Country,
	)

	if err != nil {
		return nil, nil, err
	}

	return releases, resp, nil
}
