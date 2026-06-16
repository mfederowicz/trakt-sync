// Package handlers used to handle module actions
package handlers

import (
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

// SyncGetWatchedHandler struct for handler
type SyncGetWatchedHandler struct{ common CommonLogic }

// Handle to handle sync: get_watched action
func (m SyncGetWatchedHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	printer.Println("Get watched type:", options.Type)
	items, err := m.syncGetWatchedItems(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get watched error:%w", err)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (m SyncGetWatchedHandler) syncGetWatchedItems(client *internal.Client, options *str.Options, page int) ([]*str.UserWatched, error) {
	opts := uri.ListOptions{Page: page, Limit: consts.PerPage, Extended: options.ExtendedInfo}
	items, resp, err := client.Sync.GetWatched(
		client.BuildCtxFromOptions(options),
		&options.Type,
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
		nextPageItems, err := m.syncGetWatchedItems(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		items = append(items, nextPageItems...)
	}

	return items, nil
}
