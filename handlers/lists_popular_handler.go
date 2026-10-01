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

// ListsPopularHandler struct for handler
type ListsPopularHandler struct{}

// Handle to handle lists: popular action
func (h ListsPopularHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns the most popular lists. Popularity is calculated using total number of likes and comments.")
	result, err := h.fetchListsPopular(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch lists error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty lists")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.List{}
	exportJSON = append(exportJSON, result...)
	jsonData, err := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (h ListsPopularHandler) fetchListsPopular(client *trakt.Client, options *str.Options, page int) ([]*str.List, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	ctx := cli.ContextFromOptions(options)
	var list []*str.List
	var resp *str.Response
	var err error
	if len(options.Type) > consts.ZeroValue {
		list, resp, err = client.Lists.GetPopularListsByType(ctx, options.Type, &opts)
	} else {
		list, resp, err = client.Lists.GetPopularLists(ctx, &opts)
	}

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		waitPageDelay()

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := h.fetchListsPopular(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
