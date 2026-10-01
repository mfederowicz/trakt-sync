// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsPopularHandler struct for handler
type ShowsPopularHandler struct{}

// Handle to handle shows: popular action
func (h ShowsPopularHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns the most popular shows. Popularity is calculated using the rating percentage and the number of ratings.")
	result, err := h.fetchShowsPopular(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch shows error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty lists")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.Show{}
	exportJSON = append(exportJSON, result...)
	jsonData, err := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (h ShowsPopularHandler) fetchShowsPopular(client *trakt.Client, options *str.Options, page int) ([]*str.Show, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Shows.GetPopularShows(
		cli.ContextFromOptions(options),
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := h.fetchShowsPopular(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
