// Package cmds used for commands modules
package cmds

import (
	"flag"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

var (
	username   = "me"
	exportData []*str.PersonalList

	_listID = flag.String("i", cfg.DefaultConfig().ID, consts.UserlistUsage)
	_usersAction = UsersCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
)

// UsersCmd Returns all personal lists for a user.
var UsersCmd = &Command{
	Name:    "users",
	Usage:   "",
	Summary: "Returns all data for a users.",
	Help:    `users command`,
}

func usersListsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.UsersHandler
	switch options.Action {
	case "lists":
		handler = handlers.UsersListsHandler{}
	case "saved_filters":
		handler = handlers.UsersSavedFiltersHandler{}
	case "stats":
		handler = handlers.UsersStatsHandler{}
	default:
		printer.Println("possible actions: lists, saved_filters, stats")
	}
	err := handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	usersListItemsDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	UsersCmd.Run = usersListsFunc
}
