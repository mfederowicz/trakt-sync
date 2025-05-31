// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_episodesAction     = EpisodesCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_episodesInternalID = EpisodesCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.InternalIDUsage)
	_episodesEpisode    = EpisodesCmd.Flag.Int("episode", cfg.DefaultConfig().Episode, consts.EpisodeUsage)
	_episodesSeason     = EpisodesCmd.Flag.Int("season", cfg.DefaultConfig().Season, consts.SeasonUsage)
	_episodesLanguage   = EpisodesCmd.Flag.String("language", cfg.DefaultConfig().Language, consts.LanguageUsage)
	_episodesSort       = EpisodesCmd.Flag.String("s", cfg.DefaultConfig().EpisodesSort, consts.SortUsage)
	_episodesType       = EpisodesCmd.Flag.String("t", cfg.DefaultConfig().EpisodesType, consts.TypeUsage)

	_episodesActions = []string{
		"summary", "translations", "comments", "lists",
		"people", "ratings", "stats", "watching", "videos"}
)

// EpisodesCmd returns episodes and episodes that a user has watched, sorted by most recent.
var EpisodesCmd = &Command{
	Name:    "episodes",
	Usage:   "",
	Summary: "Returns data about episodes: summary, season, episodes, translations, comments etc...",
	Help:    `episodes command`,
}

func episodesFunc(cmd *Command, _ ...string) error {
	cmd.UpdateEpisodeFlagsValues()
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

	var handler handlers.EpisodesHandler
	allHandlers := map[string]handlers.Handler{
		"summary":      handlers.EpisodesSummaryHandler{},
		"translations": handlers.EpisodesTranslationsHandler{},
		"comments":     handlers.EpisodesCommentsHandler{},
		"lists":        handlers.EpisodesListsHandler{},
		"people":       handlers.EpisodesPeopleHandler{},
		// "ratings":      handlers.EpisodesRatingsHandler{},
		// "stats":        handlers.EpisodesStatsHandler{},
		// "watching":     handlers.EpisodesWatchingHandler{},
	}
	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, _episodesActions)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	episodesDumpTemplate = ``
)

func init() {
	EpisodesCmd.Run = episodesFunc
}
