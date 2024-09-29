// Package cmds used for commands modules
package cmds

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strconv"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

var (
	username   = "me"
	exportData []*str.PersonalList

	_listID = flag.String("i", cfg.DefaultConfig().ID, consts.UserlistUsage)
)

// UsersListItemsCmd Returns all personal lists for a user.
var UsersListItemsCmd = &Command{
	Name:    "lists",
	Usage:   "",
	Summary: "Returns all personal lists for a user.",
	Help:    `lists command`,
}

func usersListItemsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	intID, _ := strconv.Atoi(*_listID)

	fmt.Println("fetch private lists for:" + options.UserName)

	username = options.UserName
	personalLists, _, err := fetchUsersPersonalLists(client, &username)
	if err != nil {
		return fmt.Errorf("fetch user list error:%w", err)
	}

	if len(personalLists) == consts.ZeroValue {
		return fmt.Errorf("empty personal lists")
	}

	fmt.Printf("Found %d user list\n", len(personalLists))

	var avLists []int

	for _, data := range personalLists {
		fmt.Printf("Found list id %d name '%s' with %d items own by %s\n", *data.IDs.Trakt, *data.Name, *data.ItemCount, *data.User.Name)
		avLists = append(avLists, int(*data.IDs.Trakt))
	}

	if intID == consts.ZeroValue {
		return fmt.Errorf("please set personal listid")
	}

	if !str.ContainInt(intID, avLists) {
		return fmt.Errorf("unknown listid:%d", intID)
	}

	fmt.Printf("ListId to fetch:%d\n", intID)

	if len(*_output) > consts.ZeroValue {
		options.Output = *_output
	} else {
		options.Output = fmt.Sprintf("export_%s_%s.json", options.Module, options.Type)
	}

	if intID > consts.ZeroValue && str.ContainInt(intID, avLists) {
		options.ID = strconv.Itoa(intID)
		itemsExportData, _, itemsErr := fetchUsersPersonalList(client, options)
		if itemsErr == nil {
			if len(itemsExportData) > consts.ZeroValue {
				fmt.Printf("Found %d items \n", len(itemsExportData))
				exportJSON := []*str.UserListItem{}
				exportJSON = append(exportJSON, itemsExportData...)
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")
				writer.WriteJSON(options, jsonData)
			} else {
				fmt.Printf("No %s items in list %d to fetch\n", options.Type, intID)
			}
		}
	}
	return nil
}

var (
	usersListItemsDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	UsersListItemsCmd.Run = usersListItemsFunc
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
	lists, resp, err := client.Users.GetItemstOnAPersonalList(
		context.Background(),
		&username,
		&listIDString,
		&options.Type,
	)

	return lists, resp, err
}
