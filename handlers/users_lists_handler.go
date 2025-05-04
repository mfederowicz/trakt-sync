// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersListsHandler struct for handler
type UsersListsHandler struct{}

// Handle to handle users: lists action
func (UsersListsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("fetch private lists for:" + options.UserName)

	username := options.UserName
	personalLists, _, err := fetchUsersPersonalLists(client, &username)
	if err != nil {
		return fmt.Errorf("fetch user list error:%w", err)
	}

	if len(personalLists) == consts.ZeroValue {
		return errors.New("empty personal lists")
	}

	printer.Printf("Found %d user list\n", len(personalLists))

	avLists := getAvlistsFromPersonals(personalLists)

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
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")
	writer.WriteJSON(options, jsonData)

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

func fetchUsersPersonalLists(client *internal.Client, username *string) ([]*str.PersonalList, *str.Response, error) {
	lists, resp, err := client.Users.GetUsersPersonalLists(
		context.Background(),
		username,
	)

	return lists, resp, err
}

func fetchUsersPersonalList(client *internal.Client, options *str.Options) ([]*str.UserListItem, *str.Response, error) {
	listIDString := options.ID
	username := options.UserName
	lists, resp, err := client.Users.GetItemstOnAPersonalList(
		context.Background(),
		&username,
		&listIDString,
		&options.Type,
	)

	return lists, resp, err
}
