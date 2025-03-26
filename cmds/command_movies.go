// Package cmds used for commands modules
package cmds

import (
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_moviesAction     = MoviesCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_moviesInternalID = MoviesCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.MovieIDUsage)
	_moviesPeriod     = MoviesCmd.Flag.String("period", cfg.DefaultConfig().MoviesPeriod, consts.MoviesPeriodUsage)
	_moviesCountry    = MoviesCmd.Flag.String("country", cfg.DefaultConfig().MoviesCountry, consts.MoviesCountryUsage)
	_moviesLanguage   = MoviesCmd.Flag.String("language", cfg.DefaultConfig().MoviesLanguage, consts.MoviesLanguageUsage)
	_moviesSort       = MoviesCmd.Flag.String("s", cfg.DefaultConfig().MoviesSort, consts.MoviesSortUsage)
	_moviesType       = MoviesCmd.Flag.String("t", cfg.DefaultConfig().MoviesType, consts.MoviesTypeUsage)
	_moviesStartDate  = MoviesCmd.Flag.String("start_date", "", consts.StartDateUsage)
)

// MoviesCmd returns movies and episodes that a user has watched, sorted by most recent.
var MoviesCmd = &Command{
	Name:    "movies",
	Usage:   "",
	Summary: "Returns data about movies: trending, popular, list, likes, like, items, comments.",
	Help:    `movies command`,
}

func moviesFunc(cmd *Command, _ ...string) error {
	cmd.UpdateMovieFlagsValues()
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

	var handler handlers.MoviesHandler

	handler, err = getHandlerFromAction(options.Action)

	if err != nil {
		cmd.GenActionsUsage(validActions)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

// validActions holds the list of acceptable actions.
var validActions = []string{
	"trending", "popular", "favorited", "played", "watched", "collected",
	"anticipated", "boxoffice", "updated", "updated_ids", "summary", "aliases",
	"releases", "translations", "comments", "lists", "people", "ratings",
	"releated", "stats", "studios", "watching", "videos", "refresh"}

// getHandlerFromAction choose handler from all
func getHandlerFromAction(action string) (handlers.MoviesHandler, error) {
	allHandlers := map[string]handlers.MoviesHandler{
		"trending":     handlers.MoviesTrendingHandler{},
		"popular":      handlers.MoviesPopularHandler{},
		"favorited":    handlers.MoviesFavoritedHandler{},
		"played":       handlers.MoviesPlayedHandler{},
		"watched":      handlers.MoviesWatchedHandler{},
		"collected":    handlers.MoviesCollectedHandler{},
		"anticipated":  handlers.MoviesAnticipatedHandler{},
		"boxoffice":    handlers.MoviesBoxofficeHandler{},
		"updates":      handlers.MoviesUpdatesHandler{},
		"updated_ids":  handlers.MoviesUpdatedIDsHandler{},
		"summary":      handlers.MoviesSummaryHandler{},
		"aliases":      handlers.MoviesAliasesHandler{},
		"releases":     handlers.MoviesReleasesHandler{},
		"translations": handlers.MoviesTranslationsHandler{},
		"comments":     handlers.MoviesCommentsHandler{},
		"lists":        handlers.MoviesListsHandler{},
		"people":       handlers.MoviesPeopleHandler{},
		"ratings":      handlers.MoviesRatingsHandler{},
		"related":      handlers.MoviesRelatedHandler{},
		"stats":        handlers.MoviesStatsHandler{},
		"studios":      handlers.MoviesStudiosHandler{},
		"watching":     handlers.MoviesWatchingHandler{},
		"videos":       handlers.MoviesVideosHandler{},
		"refresh":      handlers.MoviesRefreshHandler{},
	}

	// Lookup and execute handler
	if handler, found := allHandlers[action]; found {
		return handler, nil
	}

	return nil, errors.New("unknown handler")
}

var (
	moviesDumpTemplate = ``
)

func init() {
	MoviesCmd.Run = moviesFunc
}
