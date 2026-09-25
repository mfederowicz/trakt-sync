// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_socialRecommendationsAction            = SocialRecommendationsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_socialRecommendationsIgnoreCollected   = SocialRecommendationsCmd.Flag.String("ignore_collected", cfg.DefaultConfig().IgnoreCollected, consts.IgnoreCollectedUsage)
	_socialRecommendationsIgnoreWatched     = SocialRecommendationsCmd.Flag.String("ignore_watched", cfg.DefaultConfig().IgnoreWatched, consts.IgnoreWatchedUsage)
	_socialRecommendationsIgnoreWatchlisted = SocialRecommendationsCmd.Flag.String("ignore_watchlisted", cfg.DefaultConfig().IgnoreWatchlisted, consts.IgnoreWatchlistedUsage)
	_socialRecommendationsWatchWindow       = SocialRecommendationsCmd.Flag.Int("watch_window", cfg.DefaultConfig().WatchWindow, consts.WatchWindowUsage)
)

// SocialRecommendationsCmd returns movie and show recommendations based on the people the user follows.
var SocialRecommendationsCmd = &Command{
	Name:    "social_recommendations",
	Usage:   "",
	Summary: "Movie and show recommendations based on the people you follow.",
	Help:    `social_recommendations command`,
}

func socialRecommendationsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.SocialRecommendationsHandler
	allHandlers := map[string]handlers.Handler{
		consts.Movies: handlers.SocialRecommendationsMoviesHandler{},
		consts.Shows:  handlers.SocialRecommendationsShowsHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.Movies, consts.Shows})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	SocialRecommendationsCmd.Run = socialRecommendationsFunc
}
