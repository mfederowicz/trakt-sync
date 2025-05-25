// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// ShowsCollectedHandler struct for handler
type ShowsCollectedHandler struct{}

// Handle to handle shows: collected action
func (h ShowsCollectedHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns the most collected (unique users) shows in the specified time period, defaulting to weekly.")
	result, err := h.fetchShowsCollected(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch shows error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty shows")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.ShowsItem{}
	exportJSON = append(exportJSON, result...)
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (h ShowsCollectedHandler) fetchShowsCollected(client *internal.Client, options *str.Options, page int) ([]*str.ShowsItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	period := options.Period
	list, resp, err := client.Shows.GetCollectedShows(
		client.BuildCtxFromOptions(options),
		&opts,
		&period,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := h.fetchShowsCollected(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
