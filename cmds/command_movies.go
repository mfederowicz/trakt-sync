// Package cmds used for commands modules
package cmds

import (
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
	_moviesReason     = MoviesCmd.Flag.String("r", cfg.DefaultConfig().Reason, consts.ReasonUsage)
	_moviesMessage    = MoviesCmd.Flag.String("message", cfg.DefaultConfig().Msg, consts.ReportMsgUsage)
	_moviesLinks      = MoviesCmd.Flag.String("links", consts.EmptyString, consts.MoviesLinksUsage)

	validMoviesActions = []string{
		consts.Trending, consts.Popular, consts.Favorited, consts.Played, consts.Watched, consts.Collected,
		consts.Anticipated, consts.Boxoffice, consts.Updates, consts.UpdatedIDs, consts.Summary, consts.Aliases,
		consts.Releases, consts.Translations, consts.Comments, consts.Lists, consts.People, consts.Ratings,
		consts.Related, consts.Stats, consts.Studios, consts.Watching, consts.Videos, consts.Refresh,
		consts.Hot, consts.Streaming, consts.Report, consts.RefreshJustwatch, consts.Sentiments,
		consts.WatchNow, consts.JustwatchLinks}
)

// MoviesCmd returns movies and episodes that a user has watched, sorted by most recent.
var MoviesCmd = &Command{
	Name:    "movies",
	Usage:   "",
	Summary: "Returns data about movies: trending, popular, list, likes, like, items, comments etc...",
	Help:    `movies command`,
}

func moviesFunc(cmd *Command, _ ...string) error {
	cmd.UpdateMovieFlagsValues()
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	err := cmd.ValidPeriodForModule(options)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	err = cmd.ValidSort(options)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	var handler handlers.MoviesHandler
	allHandlers := map[string]handlers.Handler{
		consts.Trending:     handlers.MoviesTrendingHandler{},
		consts.Popular:      handlers.MoviesPopularHandler{},
		consts.Favorited:    handlers.MoviesFavoritedHandler{},
		consts.Played:       handlers.MoviesPlayedHandler{},
		consts.Watched:      handlers.MoviesWatchedHandler{},
		consts.Collected:    handlers.MoviesCollectedHandler{},
		consts.Anticipated:  handlers.MoviesAnticipatedHandler{},
		consts.Boxoffice:    handlers.MoviesBoxofficeHandler{},
		consts.Updates:      handlers.MoviesUpdatesHandler{},
		consts.UpdatedIDs:   handlers.MoviesUpdatedIDsHandler{},
		consts.Summary:      handlers.MoviesSummaryHandler{},
		consts.Aliases:      handlers.MoviesAliasesHandler{},
		consts.Releases:     handlers.MoviesReleasesHandler{},
		consts.Translations: handlers.MoviesTranslationsHandler{},
		consts.Comments:     handlers.MoviesCommentsHandler{},
		consts.Lists:        handlers.MoviesListsHandler{},
		consts.People:       handlers.MoviesPeopleHandler{},
		consts.Ratings:      handlers.MoviesRatingsHandler{},
		consts.Related:      handlers.MoviesRelatedHandler{},
		consts.Stats:        handlers.MoviesStatsHandler{},
		consts.Studios:      handlers.MoviesStudiosHandler{},
		consts.Watching:     handlers.MoviesWatchingHandler{},
		consts.Videos:       handlers.MoviesVideosHandler{},
		consts.Refresh:      handlers.MoviesRefreshHandler{},

		consts.Hot:       handlers.MoviesHotHandler{},
		consts.Streaming: handlers.MoviesStreamingHandler{},

		consts.Report:           handlers.MoviesReportHandler{},
		consts.RefreshJustwatch: handlers.MoviesRefreshJustwatchHandler{},
		consts.Sentiments:       handlers.MoviesSentimentsHandler{},

		consts.WatchNow:       handlers.MoviesWatchNowHandler{},
		consts.JustwatchLinks: handlers.MoviesJustwatchLinksHandler{},
	}
	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validMoviesActions)
		return unknownActionError(cmd.Name, options.Action)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

var (
	moviesDumpTemplate = ``
)

func init() {
	MoviesCmd.Run = moviesFunc
}
