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

// ListsLikesHandler struct for handler
type ListsLikesHandler struct{}

// Handle to handle lists: likes action
func (h ListsLikesHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyListIDMsg)
	}
	printer.Println("Returns all users who liked a list.")
	result, err := h.fetchListsLikes(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch users error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New("empty users")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.UserLike{}
	exportJSON = append(exportJSON, result...)
	jsonData, err := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (h ListsLikesHandler) fetchListsLikes(client *trakt.Client, options *str.Options, page int) ([]*str.UserLike, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Lists.GetAllUsersWhoLikedList(
		cli.ContextFromOptions(options),
		&opts,
		options.InternalID,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		waitPageDelay()

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := h.fetchListsLikes(client, options, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}
