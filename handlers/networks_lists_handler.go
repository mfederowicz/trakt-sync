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

// NetworksListsHandler struct for handler
type NetworksListsHandler struct{}

// Handle to handle networks: lists action
func (p NetworksListsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Get a list of all TV networks")
	result, err := p.fetchNetworksList(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch lists error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty lists")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.TvNetwork{}
	exportJSON = append(exportJSON, result...)
	jsonData, err := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (p NetworksListsHandler) fetchNetworksList(client *trakt.Client, options *str.Options, page int) ([]*str.TvNetwork, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Networks.GetNetworksList(
		cli.ContextFromOptions(options),
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		waitPageDelay()
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := p.fetchNetworksList(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}
