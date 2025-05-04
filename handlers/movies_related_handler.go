// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// MoviesRelatedHandler struct for handler
type MoviesRelatedHandler struct{ common CommonLogic }

// Handle to handle movies: related action
func (m MoviesRelatedHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns related and similar movies.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	result, _, err := m.fetchMoviesRelated(client, options, consts.DefaultPage)

	if err != nil {
		return err
	}

	printer.Printf("Found related for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (m MoviesRelatedHandler) fetchMoviesRelated(client *internal.Client, options *str.Options, page int) ([]*str.Movie, *str.Response, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Movies.GetRelatedMovies(
		context.Background(),
		&options.InternalID,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, _, err := m.fetchMoviesRelated(client, options, nextPage)
		if err != nil {
			return nil, nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, resp, nil
}
