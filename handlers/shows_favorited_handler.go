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

// ShowsFavoritedHandler struct for handler
type ShowsFavoritedHandler struct{}

// Handle to handle shows: favorited action
func (h ShowsFavoritedHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns the most favorited shows in the specified time period")
	result, err := h.fetchShowsFavorited(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch shows error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty shows")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.ShowsItem{}
	exportJSON = append(exportJSON, result...)
	jsonData, err := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (h ShowsFavoritedHandler) fetchShowsFavorited(client *trakt.Client, options *str.Options, page int) ([]*str.ShowsItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	period := options.Period
	list, resp, err := client.Shows.GetFavoritedShows(
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
		nextPageItems, err := h.fetchShowsFavorited(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
