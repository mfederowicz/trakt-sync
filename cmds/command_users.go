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

	_usersListID          = flag.String("i", cfg.DefaultConfig().ID, consts.UserlistUsage)
	_usersAction          = UsersCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_usersType            = UsersCmd.Flag.String("t", cfg.DefaultConfig().UsersType, consts.UsersTypeUsage)
	_usersSection         = UsersCmd.Flag.String("s", cfg.DefaultConfig().UsersSection, consts.UsersSectionUsage)
	_usersDeny            = UsersCmd.Flag.Bool("deny", cfg.DefaultConfig().Deny, consts.DenyUsage)
	_usersFollowerRequest = UsersCmd.Flag.Int("follower_request", cfg.DefaultConfig().FollowerRequest, consts.FollowerRequestUsage)
	_usersItems           = UsersCmd.Flag.String("items", consts.EmptyString, consts.ItemsUsage)
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

	err = cmd.ValidSection(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	var handler handlers.UsersHandler
	allHandlers := map[string]handlers.Handler{
		"settings":            handlers.UsersSettingsHandler{},
		"following_requests":  handlers.UsersFollowingRequestsHandler{},
		"follower_requests":   handlers.UsersFollowerRequestsHandler{},
		"saved_filters":       handlers.UsersSavedFiltersHandler{},
		"hidden_items":        handlers.UsersHiddenItemsHandler{},
		"add_hidden_items":    handlers.UsersAddHiddenItemsHandler{},
		"remove_hidden_items": handlers.UsersRemoveHiddenItemsHandler{},
		"profile":             handlers.UsersProfileHandler{},
		"lists":               handlers.UsersListsHandler{},
		"stats":               handlers.UsersStatsHandler{},
		"watched":             handlers.UsersWatchedHandler{},
	}

	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	validActions = []string{"settings", "following_requests", "follower_requests",
		"follow_request", "saved_filters", "hidden_items", "add_hidden_items",
		"remove_hidden_items", "profile", "lists", "stats", "watched"}
	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validActions)
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
