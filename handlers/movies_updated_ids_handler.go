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

// MoviesUpdatedIDsHandler struct for handler
type MoviesUpdatedIDsHandler struct{}

// Handle to handle people: updated_ids action
func (m MoviesUpdatedIDsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get recently updated movies Trakt IDs for date:" + options.StartDate)
	date := options.StartDate
	updates, err := m.fetchMoviesUpdatedIDs(client, options, date, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch movies updated ids error:%w", err)
	}

	if len(updates) == consts.ZeroValue {
		return errors.New("empty updated ids lists")
	}

	if len(updates) > consts.ZeroValue {
		printer.Printf("Found %d items \n", len(updates))
		exportJSON := []*int{}
		exportJSON = append(exportJSON, updates...)
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")

		writer.WriteJSON(options, jsonData)
	} else {
		printer.Print("No update ids items to fetch\n")
	}

	return nil
}

func (m MoviesUpdatedIDsHandler) fetchMoviesUpdatedIDs(client *internal.Client, options *str.Options, startDate string, page int) ([]*int, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Movies.GetRecentlyUpdatedMoviesTraktIDs(
		client.BuildCtxFromOptions(options),
		&startDate,
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
		nextPageItems, err := m.fetchMoviesUpdatedIDs(client, options, startDate, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
