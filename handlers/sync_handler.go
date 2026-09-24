// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// SyncHandler interface to handle sync module action
type SyncHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}

// syncProgressFetch fetches one page of up next or watched progress
type syncProgressFetch func(ctx context.Context, opts *uri.SyncProgressOptions) ([]*str.ShowProgress, *str.Response, error)

// fetchSyncProgress fetches every page of up next or watched progress
func fetchSyncProgress(client *internal.Client, options *str.Options, opts uri.SyncProgressOptions, page int, fetch syncProgressFetch) ([]*str.ShowProgress, error) {
	opts.Page = page
	list, resp, err := fetch(client.BuildCtxFromOptions(options), &opts)
	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPageItems, err := fetchSyncProgress(client, options, opts, page+consts.NextPageStep, fetch)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

// writeSyncProgress writes up next or watched progress to the output file
func writeSyncProgress(options *str.Options, result []*str.ShowProgress) error {
	printer.Printf("Found %d shows\n", len(result))
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}

	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}
