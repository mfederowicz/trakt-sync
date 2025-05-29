// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_seasonsAction     = SeasonsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_seasonsInternalID = SeasonsCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.InternalIDUsage)
	_seasonsSeason     = SeasonsCmd.Flag.Int("season", cfg.DefaultConfig().Season, consts.SeasonUsage)
	_seasonsLanguage   = SeasonsCmd.Flag.String("language", cfg.DefaultConfig().Language, consts.LanguageUsage)
	_seasonsSort       = SeasonsCmd.Flag.String("s", cfg.DefaultConfig().SeasonsSort, consts.SortUsage)
	_seasonsType       = SeasonsCmd.Flag.String("t", cfg.DefaultConfig().SeasonsType, consts.TypeUsage)

	_seasonsActions = []string{
		"summary", "season", "episodes", "translations", "comments", "lists",
		"people", "ratings", "stats", "watching", "videos"}
)

// SeasonsCmd returns seasons and episodes that a user has watched, sorted by most recent.
var SeasonsCmd = &Command{
	Name:    "seasons",
	Usage:   "",
	Summary: "Returns data about seasons: summary, season, episodes, translations, comments etc...",
	Help:    `seasons command`,
}

func seasonsFunc(cmd *Command, _ ...string) error {
	cmd.UpdateSeasonFlagsValues()
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	err := cmd.ValidPeriodForModule(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	err = cmd.ValidSort(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	var handler handlers.SeasonsHandler
	allHandlers := map[string]handlers.Handler{
		"summary":      handlers.SeasonsSummaryHandler{},
		"season":       handlers.SeasonsSeasonHandler{},
		"episodes":     handlers.SeasonsEpisodesHandler{},
		"translations": handlers.SeasonsTranslationsHandler{},
		"comments":     handlers.SeasonsCommentsHandler{},
		"lists":        handlers.SeasonsListsHandler{},
		"people":       handlers.SeasonsPeopleHandler{},
		"ratings":      handlers.SeasonsRatingsHandler{},
		"stats":        handlers.SeasonsStatsHandler{},
		"watching":     handlers.SeasonsWatchingHandler{},
		// "videos":       handlers.SeasonsVideosHandler{},
	}
	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, _seasonsActions)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	seasonsDumpTemplate = ``
)

func init() {
	SeasonsCmd.Run = seasonsFunc
}
