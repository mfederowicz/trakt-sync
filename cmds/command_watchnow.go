// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_watchNowAction  = WatchNowCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_watchNowCountry = WatchNowCmd.Flag.String("country", cfg.DefaultConfig().MoviesCountry, consts.WatchNowCountryUsage)
)

// WatchNowCmd returns watch now sources (streaming providers) supported by Trakt.
var WatchNowCmd = &Command{
	Name:    "watchnow",
	Usage:   "",
	Summary: "Returns watch now sources (streaming providers), all or by country.",
	Help:    `watchnow command`,
}

func watchNowFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.WatchNowHandler
	allHandlers := map[string]handlers.Handler{
		consts.Sources: handlers.WatchNowSourcesHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.Sources})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	WatchNowCmd.Run = watchNowFunc
}
