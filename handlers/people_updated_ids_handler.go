// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// PeopleUpdatedIDsHandler struct for handler
type PeopleUpdatedIDsHandler struct{}

// Handle to handle people: updated_ids action
func (p PeopleUpdatedIDsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get recently updated people Trakt IDs for date:"+options.StartDate)
	date := options.StartDate
	updates, err := p.fetchPeoplesUpdatedIDs(client, options, date, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch peoples updated ids error:%w", err)
	}

	if len(updates) == consts.ZeroValue {
		return fmt.Errorf("empty updated ids lists")
	}

	if len(updates) > consts.ZeroValue {
		printer.Printf("Found %d items \n", len(updates))
		exportJSON := []*int{}
		exportJSON = append(exportJSON, updates...)
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")

		writer.WriteJSON(options, jsonData)
	} else {
		printer.Print("No update ids items to fetch\n")
	}

	return nil
}

func (p PeopleUpdatedIDsHandler) fetchPeoplesUpdatedIDs(client *internal.Client, options *str.Options, startDate string, page int) ([]*int, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.People.GetRecentlyUpdatedPeopleTraktIDs(
		context.Background(),
		&startDate,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	pages, _ := strconv.Atoi(resp.Header.Get(internal.HeaderPaginationPageCount))
	if client.HavePages(page, pages) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)

		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := p.fetchPeoplesUpdatedIDs(client, options, startDate, nextPage)
		if err != nil {
			return nil, err
		}

		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}

	return list, nil
}

