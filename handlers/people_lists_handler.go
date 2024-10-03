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

// PeopleListsHandler struct for handler
type PeopleListsHandler struct{}

// Handle to handle people: shows action
func (p PeopleListsHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return fmt.Errorf(consts.EmptyPersonIDMsg)
	}
	printer.Println("Get lists containing this person")
	result, err := p.fetchListsContainingThisPerson(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch lists error:%v", err)
	}

	if len(result) == consts.ZeroValue {
		return fmt.Errorf("empty lists")
	}

	printer.Printf("Found %d result \n", len(result))
	exportJSON := []*str.PersonalList{}
	exportJSON = append(exportJSON, result...)
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)

	writer.WriteJSON(options, jsonData)

	return nil
}

func (p PeopleListsHandler) fetchListsContainingThisPerson(client *internal.Client, options *str.Options, page int) ([]*str.PersonalList, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.People.GetListsContainingThisPerson(
		context.Background(),
		&options.ID,
		&options.Type,
		&options.Sort,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	pages, _ := strconv.Atoi(resp.Header.Get(internal.HeaderPaginationPageCount))
	// Check if there are more pages
	if client.HavePages(page, pages) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		// Fetch items from the next page
		nextPage := page + consts.NextPageStep
		nextPageItems, err := p.fetchListsContainingThisPerson(client, options, nextPage)
		if err != nil {
			return nil, err
		}
		// Append items from the next page to the current page
		list = append(list, nextPageItems...)
	}
	return list, nil
}
