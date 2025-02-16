// Package handlers used to handle module actions
package handlers

import (
	"context"
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

// ListsLikesHandler struct for handler
type ListsLikesHandler struct{}

// Handle to handle lists: likes action
func (h ListsLikesHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.TraktID == consts.ZeroValue {
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
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (h ListsLikesHandler) fetchListsLikes(client *internal.Client, options *str.Options, page int) ([]*str.UserLike, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Lists.GetAllUsersWhoLikedList(
		context.Background(),
		&opts,
		&options.TraktID,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

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
