// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_teamAction = TeamCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
)

// TeamCmd returns Trakt team members.
var TeamCmd = &Command{
	Name:    "team",
	Usage:   "",
	Summary: "Returns Trakt team members.",
	Help:    `team command`,
}

func teamFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.TeamHandler
	allHandlers := map[string]handlers.Handler{
		consts.Members: handlers.TeamMembersHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.Members})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	TeamCmd.Run = teamFunc
}
