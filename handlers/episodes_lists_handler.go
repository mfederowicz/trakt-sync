// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// EpisodesListsHandler struct for handler
type EpisodesListsHandler struct{ common CommonLogic }

// Handle to handle seasons: lists action
func (m EpisodesListsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("Returns all lists that contain this season.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	err := m.common.CheckSortAndTypes(options)

	if err != nil {
		return err
	}

	result, _, err := m.fetchEpisodesLists(client, options, consts.DefaultPage)

	if err != nil {
		return err
	}

	printer.Printf("Found lists for id:%s\n", options.InternalID)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (m EpisodesListsHandler) fetchEpisodesLists(client *trakt.Client, options *str.Options, page int) ([]*str.PersonalList, *str.Response, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.Shows.GetListsContainingEpisode(
		cli.ContextFromOptions(options),
		&options.InternalID,
		&options.Season,
		&options.Episode,
		&options.Type,
		&options.Sort,
		&opts,
	)

	if err != nil {
		return nil, nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, _, err := m.fetchEpisodesLists(client, options, nextPage)
		if err != nil {
			return nil, nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, resp, nil
}
