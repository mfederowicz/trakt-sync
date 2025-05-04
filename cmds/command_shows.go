// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_showsAction     = ShowsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_showsInternalID = ShowsCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.MovieIDUsage)
	_showsPeriod     = ShowsCmd.Flag.String("period", cfg.DefaultConfig().ShowsPeriod, consts.ShowsPeriodUsage)
	_showsCountry    = ShowsCmd.Flag.String("country", cfg.DefaultConfig().ShowsCountry, consts.ShowsCountryUsage)
	_showsLanguage   = ShowsCmd.Flag.String("language", cfg.DefaultConfig().ShowsLanguage, consts.ShowsLanguageUsage)
	_showsSort       = ShowsCmd.Flag.String("s", cfg.DefaultConfig().ShowsSort, consts.ShowsSortUsage)
	_showsType       = ShowsCmd.Flag.String("t", cfg.DefaultConfig().ShowsType, consts.ShowsTypeUsage)
	_showsStartDate  = ShowsCmd.Flag.String("start_date", "", consts.StartDateUsage)

	validShowsActions = []string{
		"trending", "popular", "favorited", "played", "watched", "collected",
		"anticipated", "boxoffice", "updated", "updated_ids", "summary", "aliases",
		"releases", "translations", "comments", "lists", "people", "ratings",
		"releated", "stats", "studios", "watching", "videos", "refresh"}
)

// ShowsCmd returns movies and episodes that a user has watched, sorted by most recent.
var ShowsCmd = &Command{
	Name:    "shows",
	Usage:   "",
	Summary: "Returns data about shows: trending, popular, list, likes, like, items, comments etc...",
	Help:    `shows command`,
}

func showsFunc(cmd *Command, _ ...string) error {
	cmd.UpdateShowFlagsValues()
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

	var handler handlers.ShowsHandler
	allHandlers := map[string]handlers.Handler{
		"trending":  handlers.ShowsTrendingHandler{},
		"popular":   handlers.ShowsPopularHandler{},
		"favorited": handlers.ShowsFavoritedHandler{},
		"played":    handlers.ShowsPlayedHandler{},
		"watched":   handlers.ShowsWatchedHandler{},
		"collected": handlers.ShowsCollectedHandler{},
		// "anticipated":  handlers.ShowsAnticipatedHandler{},
		// "boxoffice":    handlers.ShowsBoxofficeHandler{},
		// "updates":      handlers.ShowsUpdatesHandler{},
		// "updated_ids":  handlers.ShowsUpdatedIDsHandler{},
		// "summary":      handlers.ShowsSummaryHandler{},
		// "aliases":      handlers.ShowsAliasesHandler{},
		// "releases":     handlers.ShowsReleasesHandler{},
		// "translations": handlers.ShowsTranslationsHandler{},
		// "comments":     handlers.ShowsCommentsHandler{},
		// "lists":        handlers.ShowsListsHandler{},
		// "people":       handlers.ShowsPeopleHandler{},
		// "ratings":      handlers.ShowsRatingsHandler{},
		// "related":      handlers.ShowsRelatedHandler{},
		// "stats":        handlers.ShowsStatsHandler{},
		// "studios":      handlers.ShowsStudiosHandler{},
		// "watching":     handlers.ShowsWatchingHandler{},
		// "videos":       handlers.ShowsVideosHandler{},
		// "refresh":      handlers.ShowsRefreshHandler{},
	}
	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validShowsActions)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	showsDumpTemplate = ``
)

func init() {
	ShowsCmd.Run = showsFunc
}
