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

// ShowsWatchedProgressHandler struct for handler
type ShowsWatchedProgressHandler struct{ common CommonLogic }

// Handle to handle shows: watched_progress action
func (m ShowsWatchedProgressHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns watched progress for a show including details on all aired seasons and episodes.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	err := m.common.CheckSortAndTypes(options)

	if err != nil {
		return err
	}

	result, err := m.fetchShowsWatchedProgress(client, options)

	if err != nil {
		return err
	}

	printer.Printf("Found watched for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (ShowsWatchedProgressHandler) fetchShowsWatchedProgress(client *internal.Client, options *str.Options) (*str.WatchedProgress, error) {
	opts := uri.ListOptions{Hidden: options.Hidden, Specials: options.Specials, CountSpecials: options.CountSpecials}

	result, err := client.Shows.GetShowWatchedProgress(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&opts,
	)
	if err != nil {
		return nil, fmt.Errorf(consts.WatchedProgressError, err)
	}

	return result, nil
}
