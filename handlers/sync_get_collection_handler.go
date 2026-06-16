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

// SyncGetCollectionHandler struct for handler
type SyncGetCollectionHandler struct{ common CommonLogic }

// Handle to handle sync: get_collection action
func (s SyncGetCollectionHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get collection type:", options.Type)
	items, err := s.syncGetCollection(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get collection error:%w", err)
	}

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (s SyncGetCollectionHandler) syncGetCollection(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	if options.Type == consts.Seasons {
		items, err := s.syncGetCollectedSeasons(client, options, page)
		if err != nil {
			return nil, err
		}

		return items, nil
	}

	items, err := s.syncGetCollected(client, options, page)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s SyncGetCollectionHandler) syncGetCollected(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Sync.GetCollection(
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
		nextPageItems, err := s.syncGetCollected(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

func (s SyncGetCollectionHandler) syncGetCollectedSeasons(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Sync.GetCollectedSeasons(
		client.BuildCtxFromOptions(options),
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
		nextPageItems, err := s.syncGetCollectedSeasons(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
