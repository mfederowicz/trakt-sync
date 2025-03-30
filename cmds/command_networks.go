// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_networksAction = NetworksCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
)

// NetworksCmd returns networks and episodes that a user has watched, sorted by most recent.
var NetworksCmd = &Command{
	Name:    "networks",
	Usage:   "",
	Summary: "Get a list of all TV networks",
	Help:    `networks command`,
}

func networksFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.NetworksHandler
	var networksHandlers = map[string]handlers.Handler{
		"list": handlers.NetworksListsHandler{},
	}
	handler, err := cmd.GetHandlerForMap(options.Action, networksHandlers)

	if err != nil {
		cmd.GenActionsUsage([]string{"list"})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	networksDumpTemplate = ``
)

func init() {
	NetworksCmd.Run = networksFunc
}
