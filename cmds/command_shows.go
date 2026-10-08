// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_showsAction         = ShowsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_showsInternalID     = ShowsCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.MovieIDUsage)
	_showsPeriod         = ShowsCmd.Flag.String("period", cfg.DefaultConfig().ShowsPeriod, consts.ShowsPeriodUsage)
	_showsCountry        = ShowsCmd.Flag.String("country", cfg.DefaultConfig().ShowsCountry, consts.ShowsCountryUsage)
	_showsHidden         = ShowsCmd.Flag.String("hidden", cfg.DefaultConfig().Hidden, consts.ShowsHiddenUsage)
	_showsSpecials       = ShowsCmd.Flag.String("specials", cfg.DefaultConfig().Specials, consts.ShowsSpecialsUsage)
	_showsCountSpecials  = ShowsCmd.Flag.String("count_specials", cfg.DefaultConfig().CountSpecials, consts.ShowsCountSpecialsUsage)
	_showsLanguage       = ShowsCmd.Flag.String("language", cfg.DefaultConfig().ShowsLanguage, consts.ShowsLanguageUsage)
	_showsLastActivity   = ShowsCmd.Flag.String("last_activity", consts.EmptyString, consts.ShowsLastActivityUsage)
	_showsSort           = ShowsCmd.Flag.String("s", cfg.DefaultConfig().ShowsSort, consts.ShowsSortUsage)
	_showsType           = ShowsCmd.Flag.String("t", cfg.DefaultConfig().ShowsType, consts.ShowsTypeUsage)
	_showsStartDate      = ShowsCmd.Flag.String("start_date", "", consts.StartDateUsage)
	_showsUndo           = ShowsCmd.Flag.Bool("undo", cfg.DefaultConfig().Undo, consts.UndoUsage)
	_showsResetAt        = ShowsCmd.Flag.String("reset_at", "", consts.ResetAtUsage)
	_showsReason         = ShowsCmd.Flag.String("r", cfg.DefaultConfig().Reason, consts.ReasonUsage)
	_showsMessage        = ShowsCmd.Flag.String("message", cfg.DefaultConfig().Msg, consts.ReportMsgUsage)
	_showsLinks          = ShowsCmd.Flag.String("links", consts.EmptyString, consts.ShowsLinksUsage)
	_showsEndDate        = ShowsCmd.Flag.String("end_date", consts.EmptyString, consts.EndDateUsage)
	_showsWatchNow       = ShowsCmd.Flag.String("watchnow", consts.EmptyString, consts.WatchNowUsage)
	_showsSubgenres      = ShowsCmd.Flag.String("subgenres", consts.EmptyString, consts.SubgenresUsage)
	_showsRatings        = ShowsCmd.Flag.String("ratings", consts.EmptyString, consts.RatingsFilterUsage)
	_showsCertifications = ShowsCmd.Flag.String("certifications", consts.EmptyString, consts.CertificationsUsage)
	_showsStatus         = ShowsCmd.Flag.String("status", consts.EmptyString, consts.ShowStatusUsage)
	_showsImdbRatings    = ShowsCmd.Flag.String("imdb_ratings", consts.EmptyString, consts.ImdbRatingsUsage)
	_showsRtMeters       = ShowsCmd.Flag.String("rt_meters", consts.EmptyString, consts.RtMetersUsage)
	_showsRtUserMeters   = ShowsCmd.Flag.String("rt_user_meters", consts.EmptyString, consts.RtUserMetersUsage)

	validShowsActions = []string{
		consts.Trending, consts.Popular, consts.Favorited, consts.Played, consts.Watched, consts.Collected,
		consts.Anticipated, consts.Updates, consts.UpdatedIDs, consts.Summary, consts.Aliases, consts.Certifications,
		consts.CollectionProgress, consts.WatchedProgress, consts.ResetShowProgress, consts.Translations, consts.Comments, consts.Lists, consts.People, consts.Ratings,
		consts.Related, consts.Stats, consts.Studios, consts.Watching, consts.NextEpisode, consts.LastEpisode, consts.Videos, consts.Refresh,
		consts.Report, consts.Sentiments, consts.WatchNow, consts.JustwatchLinks, consts.RefreshJustwatch}
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
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	err = cmd.ValidSort(options)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	var handler handlers.ShowsHandler
	allHandlers := map[string]handlers.Handler{
		consts.Trending:           handlers.ShowsTrendingHandler{},
		consts.Popular:            handlers.ShowsPopularHandler{},
		consts.Favorited:          handlers.ShowsFavoritedHandler{},
		consts.Played:             handlers.ShowsPlayedHandler{},
		consts.Watched:            handlers.ShowsWatchedHandler{},
		consts.Collected:          handlers.ShowsCollectedHandler{},
		consts.Anticipated:        handlers.ShowsAnticipatedHandler{},
		consts.Updates:            handlers.ShowsUpdatesHandler{},
		consts.UpdatedIDs:         handlers.ShowsUpdatedIDsHandler{},
		consts.Summary:            handlers.ShowsSummaryHandler{},
		consts.Aliases:            handlers.ShowsAliasesHandler{},
		consts.Certifications:     handlers.ShowsCertificationsHandler{},
		consts.Translations:       handlers.ShowsTranslationsHandler{},
		consts.Comments:           handlers.ShowsCommentsHandler{},
		consts.Lists:              handlers.ShowsListsHandler{},
		consts.CollectionProgress: handlers.ShowsCollectionProgressHandler{},
		consts.WatchedProgress:    handlers.ShowsWatchedProgressHandler{},
		consts.ResetShowProgress:  handlers.ShowsResetShowProgressHandler{},
		consts.People:             handlers.ShowsPeopleHandler{},
		consts.Ratings:            handlers.ShowsRatingsHandler{},
		consts.Related:            handlers.ShowsRelatedHandler{},
		consts.Stats:              handlers.ShowsStatsHandler{},
		consts.Studios:            handlers.ShowsStudiosHandler{},
		consts.Watching:           handlers.ShowsWatchingHandler{},
		consts.NextEpisode:        handlers.ShowsNextEpisodeHandler{},
		consts.LastEpisode:        handlers.ShowsLastEpisodeHandler{},
		consts.Videos:             handlers.ShowsVideosHandler{},
		consts.Refresh:            handlers.ShowsRefreshHandler{},

		consts.Report:     handlers.ShowsReportHandler{},
		consts.Sentiments: handlers.ShowsSentimentsHandler{},

		consts.WatchNow:       handlers.ShowsWatchNowHandler{},
		consts.JustwatchLinks: handlers.ShowsJustwatchLinksHandler{},

		consts.RefreshJustwatch: handlers.ShowsRefreshJustwatchHandler{},
	}
	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validShowsActions)
		return unknownActionError(cmd.Name, options.Action)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	ShowsCmd.Run = showsFunc
}
