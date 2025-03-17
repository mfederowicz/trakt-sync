// Package handlers used to handle module actions
package handlers

import (
	"context"
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

// MoviesUpdatesHandler struct for handler
type MoviesUpdatesHandler struct{}

// Handle to handle movies: updates action
func (p MoviesUpdatesHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get recently updated movies for date:"+options.StartDate)
	date := options.StartDate
	updates, err := p.fetchMoviesUpdates(client, options, date, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch peoples updates error:%w", err)
	}

	if len(updates) == consts.ZeroValue {
		return errors.New("empty updates lists")
	}

	if len(updates) > consts.ZeroValue {
		printer.Printf("Found %d items \n", len(updates))
		exportJSON := []*str.MoviesItem{}
		exportJSON = append(exportJSON, updates...)
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")
		writer.WriteJSON(options, jsonData)
	} else {
		printer.Print("No update items to fetch\n")
	}
	return nil
}

func (p MoviesUpdatesHandler) fetchMoviesUpdates(client *internal.Client, options *str.Options, startDate string, page int) ([]*str.MoviesItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Movies.GetRecentlyUpdatedMovies(
		context.Background(),
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
		nextPageItems, err := p.fetchMoviesUpdates(client, options, startDate, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

