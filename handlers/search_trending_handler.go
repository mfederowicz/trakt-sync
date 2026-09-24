// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SearchTrendingHandler struct for handler
type SearchTrendingHandler struct{}

// Handle to handle search: trending action
func (s SearchTrendingHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := checkSearchSingleType(options, cfg.SearchTrendingTypes); err != nil {
		return err
	}

	printer.Println("Get trending searches for: " + options.SearchType[consts.ZeroValue])
	result, err := s.fetchSearchTrending(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch %s search error: %w", options.Action, err)
	}

	printer.Printf("Found %d result\n", len(result))
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}

	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

func (s SearchTrendingHandler) fetchSearchTrending(client *internal.Client, options *str.Options, page int) ([]*str.SearchTrendingItem, error) {
	searchType := options.SearchType[consts.ZeroValue]
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, Query: options.Query}
	list, resp, err := client.Search.GetTrendingSearches(
		client.BuildCtxFromOptions(options),
		&searchType,
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
		nextPageItems, err := s.fetchSearchTrending(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
