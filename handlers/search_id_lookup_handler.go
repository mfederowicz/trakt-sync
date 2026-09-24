// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SearchIDLookupHandler struct for handler
type SearchIDLookupHandler struct{}

// Handle to handle search: id_lookup action
func (s SearchIDLookupHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get sarch: " + options.Action)
	printer.Println("search id_type: " + options.SearchIDType)
	printer.Println("search id: " + options.ID)
	printer.Println("search item_type: " + options.SearchType.String())

	result, err := s.fetchSearchIDLookup(client, options)
	if err != nil {
		return fmt.Errorf("fetch %s search error: %w", options.Action, err)
	}

	if result == nil {
		return errors.New("empty result")
	}

	printer.Print("Found " + options.Action + " search data \n")
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")

	writer.WriteJSON(options, jsonData)
	return nil
}

func (SearchIDLookupHandler) fetchSearchIDLookup(client *internal.Client, options *str.Options) ([]*str.SearchListItem, error) {
	err := checkSearchRequiredFields(options)

	if err != nil {
		return nil, err
	}

	searchType := options.SearchType.String()

	opts := uri.ListOptions{Extended: options.ExtendedInfo, Type: searchType}
	list, _, err := client.Search.GetIDLookupResults(
		client.BuildCtxFromOptions(options),
		&options.SearchIDType,
		&options.ID,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return list, nil
}
