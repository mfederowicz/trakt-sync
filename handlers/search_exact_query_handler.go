// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
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

// SearchExactQueryHandler struct for handler
type SearchExactQueryHandler struct{}

// Handle to handle search: exact_query action
func (s SearchExactQueryHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := checkSearchSingleType(options, cfg.SearchExactTypes); err != nil {
		return err
	}
	if len(options.Query) == consts.ZeroValue {
		return errors.New(consts.EmptySearchQueryMsg)
	}

	printer.Printf("Get exact search results for: %s\n", options.Query)
	result, err := s.fetchSearchExactQuery(client, options, consts.DefaultPage)
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

func (s SearchExactQueryHandler) fetchSearchExactQuery(client *internal.Client, options *str.Options, page int) ([]*str.SearchListItem, error) {
	searchType := options.SearchType[consts.ZeroValue]
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo, Query: options.Query}
	list, resp, err := client.Search.GetExactTextQueryResults(
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
		nextPageItems, err := s.fetchSearchExactQuery(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
