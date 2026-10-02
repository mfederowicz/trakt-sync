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

// MoviesWatchedHandler struct for handler
type MoviesWatchedHandler struct{}

// Handle to handle movies: watched action
func (h MoviesWatchedHandler) Handle(options *str.Options, client *trakt.Client) error {
	if err := checkMediaFilters(options); err != nil {
		return err
	}
	printer.Println("Returns the most watched (unique users) movies in the specified time period, defaulting to weekly.")
	result, err := h.fetchMoviesWatched(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch movies error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty movies")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.MoviesItem{}
	exportJSON = append(exportJSON, result...)
	jsonData, err := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (h MoviesWatchedHandler) fetchMoviesWatched(client *trakt.Client, options *str.Options, page int) ([]*str.MoviesItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, Filters: mediaFilters(options)}
	period := options.Period
	list, resp, err := client.Movies.GetWatchedMovies(
		cli.ContextFromOptions(options),
		&opts,
		period,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		waitPageDelay()

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := h.fetchMoviesWatched(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
