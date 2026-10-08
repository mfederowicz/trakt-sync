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

// MoviesFavoritedHandler struct for handler
type MoviesFavoritedHandler struct{}

// Handle to handle movies: favorited action
func (h MoviesFavoritedHandler) Handle(options *str.Options, client *trakt.Client) error {
	if err := checkMovieFilters(options); err != nil {
		return err
	}
	printer.Println("Returns the most favorited movies in the specified time period")
	result, err := h.fetchMoviesFavorited(client, options, consts.DefaultPage)
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

func (h MoviesFavoritedHandler) fetchMoviesFavorited(client *trakt.Client, options *str.Options, page int) ([]*str.MoviesItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, Filters: mediaFilters(options), Status: movieStatus(options)}
	period := options.Period
	list, resp, err := client.Movies.GetFavoritedMovies(
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
		nextPageItems, err := h.fetchMoviesFavorited(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
