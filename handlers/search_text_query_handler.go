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

// SearchTextQueryHandler struct for handler
type SearchTextQueryHandler struct{}

// Handle to handle search: text_query action
func (s SearchTextQueryHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get search: " + options.Action)
	printer.Printf("search_type: %v\n", options.SearchType.String())
	printer.Printf("search_field: %v\n", options.SearchField.String())
	printer.Println("search id_type: " + options.SearchIDType)

	result, err := s.fetchSearchTextQuery(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch %s search error: %w", options.Action, err)
	}

	if result == nil {
		return errors.New(consts.EmptyResult)
	}
	printer.Println("Found " + options.Action + " search data")
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}

	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (s SearchTextQueryHandler) fetchSearchTextQuery(client *internal.Client, options *str.Options, page int) ([]*str.SearchListItem, error) {
	err := checkSearchRequiredFields(options)
	if err != nil {
		return nil, err
	}

	searchType := options.SearchType.String()
	searchField := options.SearchField.String()
	opts := uri.ListOptions{
		Page:     page,
		Limit:    options.PerPage,
		Extended: options.ExtendedInfo,
		Query:    options.Query,
		Field:    searchField}
	list, resp, err := client.Search.GetTextQueryResults(
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
		nextPageItems, err := s.fetchSearchTextQuery(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
