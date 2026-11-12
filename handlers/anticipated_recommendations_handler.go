// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// exportAnticipatedRecommendations fetches all pages of anticipated recommendations and writes them to the output file.
// The routes take page and limit only, so extended info is not sent.
func exportAnticipatedRecommendations(client *trakt.Client, options *str.Options, fetch pageFetcher[*str.AnticipatedRecommendation]) error {
	printer.Println("Returns unreleased titles ranked for your taste (at most 100).")
	result, err := fetchAllPages(client, options, consts.DefaultPage, func(opts *uri.ListOptions) ([]*str.AnticipatedRecommendation, *str.Response, error) {
		opts.Extended = consts.EmptyString
		return fetch(opts)
	})
	if err != nil {
		return fmt.Errorf("fetch %s %s error: %w", options.Module, options.Action, err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}
