// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_smartListsAction            = SmartListsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_smartListsID                = SmartListsCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.SmartListIDUsage)
	_smartListsWatchNow          = SmartListsCmd.Flag.String("watchnow", consts.EmptyString, consts.WatchNowUsage)
	_smartListsSubgenres         = SmartListsCmd.Flag.String("subgenres", consts.EmptyString, consts.SubgenresUsage)
	_smartListsRatings           = SmartListsCmd.Flag.String("ratings", consts.EmptyString, consts.RatingsFilterUsage)
	_smartListsCertifications    = SmartListsCmd.Flag.String("certifications", consts.EmptyString, consts.CertificationsUsage)
	_smartListsIgnoreWatched     = SmartListsCmd.Flag.String("ignore_watched", cfg.DefaultConfig().IgnoreWatched, consts.IgnoreWatchedUsage)
	_smartListsIgnoreWatchlisted = SmartListsCmd.Flag.String("ignore_watchlisted", cfg.DefaultConfig().IgnoreWatchlisted, consts.IgnoreWatchlistedUsage)
)

// SmartListsCmd returns smart list definitions and the items they resolve to.
var SmartListsCmd = &Command{
	Name:    "smart_lists",
	Usage:   "",
	Summary: "Returns a smart list definition (summary) or the items it resolves to (items).",
	Help:    `smart_lists command`,
}

func smartListsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.SmartListsHandler
	allHandlers := map[string]handlers.Handler{
		consts.Summary: handlers.SmartListsSummaryHandler{},
		consts.Items:   handlers.SmartListsItemsHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.Summary, consts.Items})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	SmartListsCmd.Run = smartListsFunc
}
