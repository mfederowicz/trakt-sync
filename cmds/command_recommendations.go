// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_recommendationsInternalID        = RecommendationsCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.TraktIDUsage)
	_recommendationsAction            = RecommendationsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_recommendationsHide              = RecommendationsCmd.Flag.Bool("hide", cfg.DefaultConfig().Hide, consts.HideUsage)
	_recommendationsIgnoreCollected   = RecommendationsCmd.Flag.String("ignore_collected", cfg.DefaultConfig().IgnoreCollected, consts.IgnoreCollectedUsage)
	_recommendationsIgnoreWatchlisted = RecommendationsCmd.Flag.String("ignore_watchlisted", cfg.DefaultConfig().IgnoreWatchlisted, consts.IgnoreWatchlistedUsage)
)

// RecommendationsCmd manage movie and shows recommendations for user.
var RecommendationsCmd = &Command{
	Name:    "recommendations",
	Usage:   "",
	Summary: "Recommendations manage movie and shows recommendations for user",
	Help:    `recommendations command`,
}

func recommendationsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.RecommendationsHandler
	var recommendationsHandlers = map[string]handlers.Handler{
		"movies": handlers.RecommendationsMoviesHandler{},
		"shows":  handlers.RecommendationsShowsHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, recommendationsHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{"movies", "shows"})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	recommendationsDumpTemplate = ``
)

func init() {
	RecommendationsCmd.Run = recommendationsFunc
}
