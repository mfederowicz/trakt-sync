// Package cmds used for commands modules
package cmds

import (
	"flag"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/str"
)

var (
	username   = "me"
	exportData []*str.PersonalList

	_usersListID = flag.String("i", cfg.DefaultConfig().ID, consts.UserlistUsage)
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
	err := cmd.ValidModuleActionType(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}
	var handler handlers.UsersHandler
	allHandlers := map[string]handlers.Handler{
		"settings":      handlers.UsersSettingsHandler{},
		"lists":         handlers.UsersListsHandler{},
		"saved_filters": handlers.UsersSavedFiltersHandler{},
		"stats":         handlers.UsersStatsHandler{},
		"watched":       handlers.UsersWatchedHandler{},
	}

	handler, err = cmd.GetHandlerForMap(options.Action, allHandlers)

	validActions = []string{"lists", "saved_filters", "stats", "watched"}
	if err != nil {
		cmd.GenActionsUsage(validActions)
		return nil
	}

	err = handler.Handle(options, client)
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
