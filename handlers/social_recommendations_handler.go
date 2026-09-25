// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SocialRecommendationsHandler interface to handle social_recommendations module action
type SocialRecommendationsHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}

// exportSocialRecommendations fetches social recommendations with the ignore and watch_window options and writes them to the output file.
func exportSocialRecommendations(client *internal.Client, options *str.Options, fetch pageFetcher[*str.Recommendation]) error {
	result, err := fetchAllPages(client, options, consts.DefaultPage, func(opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
		opts.IgnoreCollected = options.IgnoreCollected
		opts.IgnoreWatched = options.IgnoreWatched
		opts.IgnoreWatchlisted = options.IgnoreWatchlisted
		opts.WatchWindow = options.WatchWindow
		return fetch(opts)
	})
	if err != nil {
		return fmt.Errorf("fetch social recommendations %s error: %w", options.Action, err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal social recommendations %s error: %w", options.Action, err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
