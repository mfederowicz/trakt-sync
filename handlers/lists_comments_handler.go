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

// ListsCommentsHandler struct for handler
type ListsCommentsHandler struct{}

// Handle to handle lists: comments action
func (h ListsCommentsHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyListIDMsg)
	}
	printer.Println("Returns all top level comments for a list.")
	result, err := h.fetchListComments(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch list error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty comments list")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.ListComment{}
	exportJSON = append(exportJSON, result...)
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (h ListsCommentsHandler) fetchListComments(client *internal.Client, options *str.Options, page int) ([]*str.ListComment, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Lists.GetListComments(
		client.BuildCtxFromOptions(options),
		&options.InternalID,
		&options.Sort,
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
		nextPageItems, err := h.fetchListComments(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}
