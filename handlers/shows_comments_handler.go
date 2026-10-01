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

// ShowsCommentsHandler struct for handler
type ShowsCommentsHandler struct{}

// Handle to handle shows: comments action
func (m ShowsCommentsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all top level comments for a show.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	result, _, err := m.fetchShowsComments(client, options, consts.DefaultPage)

	if err != nil {
		return err
	}

	printer.Printf("Found comments for id:%s\n", options.InternalID)

	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)
	return nil
}

func (m ShowsCommentsHandler) fetchShowsComments(client *trakt.Client, options *str.Options, page int) ([]*str.Comment, *str.Response, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Shows.GetAllShowComments(
		cli.ContextFromOptions(options),
		options.InternalID,
		options.Sort,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		waitPageDelay()

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, _, err := m.fetchShowsComments(client, options, nextPage)
		if err != nil {
			return nil, nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, resp, nil
}
