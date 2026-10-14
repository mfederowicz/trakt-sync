// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersListsHandler struct for handler
type UsersListsHandler struct{}

// Handle to handle users: lists action
func (UsersListsHandler) Handle(options *str.Options, client *trakt.Client) error {
	printer.Println("fetch private lists for:" + options.UserName)
	personalLists, err := fetchUsersPersonalLists(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch user list error:%w", err)
	}

	if len(personalLists) == consts.ZeroValue {
		return errors.New("empty personal lists")
	}
	// the lists overview keeps its own file; options.Output (-o) is for the list items
	overview := *options
	overview.Output = fmt.Sprintf(consts.DefaultOutputFormat2, options.Module, consts.Lists)
	printer.Printf("Found %d user list\n", len(personalLists))
	jsonData, err := json.MarshalIndent(personalLists, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + overview.Output + "\n")
	writer.WriteJSON(&overview, jsonData)

	avLists := getAvlistsFromPersonals(personalLists)

	// without -i the overview is the result
	if len(options.ID) == consts.ZeroValue {
		printer.Println(consts.UsersListsItemsHint)
		return nil
	}

	intID, _ := strconv.Atoi(options.ID)
	if intID == consts.ZeroValue {
		return errors.New("please set personal listid")
	}

	if !str.ContainInt(intID, avLists) {
		return fmt.Errorf("unknown listid:%d", intID)
	}

	printer.Printf("ListId to fetch:%d\n", intID)

	itemsExportData, _, itemsErr := fetchUsersPersonalList(client, options)

	if itemsErr != nil {
		return fmt.Errorf("users personal list error %s", itemsErr)
	}

	if len(itemsExportData) == consts.ZeroValue {
		return fmt.Errorf("no %s items in list %d to fetch", options.Type, intID)
	}

	printer.Printf("Found %d items \n", len(itemsExportData))
	exportJSON := []*str.UserListItem{}
	exportJSON = append(exportJSON, itemsExportData...)
	printer.Println("write data to:" + options.Output)
	if len(exportJSON) > 0 {
		jsonData, err := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JSONDataFormat)
		if err != nil {
			return fmt.Errorf("encode %s result: %w", options.Action, err)
		}
		writer.WriteJSON(options, jsonData)
	}

	return nil
}

func getAvlistsFromPersonals(personalLists []*str.PersonalList) []int {
	var avLists []int

	for _, data := range personalLists {
		printer.Printf("Found list id %d name '%s' with %d items own by %s\n", *data.IDs.Trakt, *data.Name, *data.ItemCount, *data.User.Name)
		avLists = append(avLists, int(*data.IDs.Trakt))
	}
	return avLists
}

func fetchUsersPersonalLists(client *trakt.Client, options *str.Options, page int) ([]*str.PersonalList, error) {
	username := options.UserName
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, SortBy: options.SortBy, SortHow: options.SortHow, Extended: options.ExtendedInfo}
	lists, resp, err := client.Users.GetUsersPersonalLists(
		cli.ContextFromOptions(options),
		username,
		&opts,
	)
	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		waitPageDelay()
		nextPageItems, err := fetchUsersPersonalLists(client, options, page+consts.NextPageStep)
		if err != nil {
			return nil, err
		}
		lists = append(lists, nextPageItems...)
	}

	return lists, nil
}

func fetchUsersPersonalList(client *trakt.Client, options *str.Options) ([]*str.UserListItem, *str.Response, error) {
	listIDString := options.ID
	username := options.UserName
	lists, resp, err := client.Users.GetListItemsByType(
		cli.ContextFromOptions(options),
		username,
		listIDString,
		options.Type,
	)

	return lists, resp, err
}
