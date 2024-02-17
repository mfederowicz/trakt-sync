package cmds

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"trakt-sync/cfg"
	"trakt-sync/internal"
	"trakt-sync/str"
	"trakt-sync/writer"
)

var (
	username    = "me"
	export_data []*str.PersonalList

	_listId = flag.String("i", cfg.DefaultConfig().Id, UserlistUsage)
)

var UsersListItemsCmd = &Command{
	Name:    "lists",
	Usage:   "",
	Summary: "Returns all personal lists for a user.",
	Help:    `lists command`,
}

func usersListItemsFunc(cmd *Command, args ...string) {
	options := cmd.Options
	client := cmd.Client
	intId, _ := strconv.Atoi(*_listId)

	fmt.Println("fetch private lists for:" + options.UserName)

	username = options.UserName
	personal_lists, _, err := fetchUsersPersonalLists(client, &username)
	if err != nil {
		fmt.Printf("fetch user list error:%v", err)
		os.Exit(0)
	}

	if len(personal_lists) == 0 {
		fmt.Print("empty personal lists")
		os.Exit(0)
	}

	fmt.Printf("Found %d user list\n", len(personal_lists))

	var av_lists []int

	for _, data := range personal_lists {
		fmt.Printf("Found list id %d name '%s' with %d items own by %s\n", *data.Ids.Trakt, *data.Name, *data.ItemCount, *data.User.Name)
		av_lists = append(av_lists, int(*data.Ids.Trakt))
	}

	if intId == 0 {
		fmt.Print("please set personal listid")
		os.Exit(0)
	}

	if !str.ContainInt(intId, av_lists) {
		fmt.Printf("unknown listid:%d\n", intId)
		os.Exit(0)
	}

	fmt.Printf("ListId to fetch:%d\n", intId)

	if len(*_output) > 0 {
		options.Output = *_output
	} else {
		options.Output = fmt.Sprintf("export_%s_%s.json", options.Module, options.Type)
	}

	if intId > 0 && str.ContainInt(intId, av_lists) {

		options.Id = strconv.Itoa(intId)
		items_export_data, _, items_err := fetchUsersPersonalList(client, options)
		if items_err == nil {
			if len(items_export_data) > 0 {
				fmt.Printf("Found %d items \n", len(items_export_data))
				export_json := []*str.UserListItem{}
				export_json = append(export_json, items_export_data...)
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(export_json, "", "  ")
				writer.WriteJson(options, jsonData)
			} else {
				fmt.Printf("No %s items in list %d to fetch\n", options.Type, intId)
			}

		}
	}

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

	listIdString := options.Id
	lists, resp, err := client.Users.GetItemstOnAPersonalList(
		context.Background(),
		&username,
		&listIdString,
		&options.Type,
	)

	return lists, resp, err

}
