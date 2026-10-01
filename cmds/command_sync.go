// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/str"
)

var (
	_syncAction               = SyncCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_syncStartAt              = SyncCmd.Flag.String("start_at", cfg.DefaultConfig().StartAt, consts.StartAtUsage)
	_syncEndAt                = SyncCmd.Flag.String("end_at", cfg.DefaultConfig().EndAt, consts.EndAtUsage)
	_syncPlaybackID           = SyncCmd.Flag.Int("playback_id", cfg.DefaultConfig().PlaybackID, consts.PlaybackIDUsage)
	_syncListItemID           = SyncCmd.Flag.Int("list_item_id", cfg.DefaultConfig().ListItemID, consts.ListItemIDUsage)
	_syncItems                = SyncCmd.Flag.String("items", consts.EmptyString, consts.ItemsUsage)
	_syncID                   = SyncCmd.Flag.Int("i", cfg.DefaultConfig().TraktID, consts.TraktIDUsage)
	_syncWatchlistDescription = SyncCmd.Flag.String("description", cfg.DefaultConfig().Description, consts.WatchlistDescriptionUsage)
	_syncWatchlistNotes       = SyncCmd.Flag.String("notes", cfg.DefaultConfig().Notes, consts.WatchlistNotesUsage)
	_syncAvailableOn          = SyncCmd.Flag.String("available_on", consts.EmptyString, consts.AvailableOnUsage)
	_syncIncludeStats         = SyncCmd.Flag.Bool("include_stats", false, consts.IncludeStatsUsage)
	_syncLifetimeStats        = SyncCmd.Flag.Bool("lifetime_stats", false, consts.LifetimeStatsUsage)
	_syncHideCompleted        = SyncCmd.Flag.Bool("hide_completed", false, consts.HideCompletedUsage)
	_syncHideNotCompleted     = SyncCmd.Flag.Bool("hide_not_completed", false, consts.HideNotCompletedUsage)
	_syncOnlyRewatching       = SyncCmd.Flag.Bool("only_rewatching", false, consts.OnlyRewatchingUsage)
	_syncIntent               = SyncCmd.Flag.String("intent", consts.EmptyString, consts.IntentUsage)
	_syncWatchNow             = SyncCmd.Flag.String("watchnow", consts.EmptyString, consts.WatchNowUsage)
	_syncSubgenres            = SyncCmd.Flag.String("subgenres", consts.EmptyString, consts.SubgenresUsage)
	_syncRatings              = SyncCmd.Flag.String("ratings", consts.EmptyString, consts.RatingsFilterUsage)
	_syncCertifications       = SyncCmd.Flag.String("certifications", consts.EmptyString, consts.CertificationsUsage)
	_syncStartDate            = SyncCmd.Flag.String("start_date", consts.EmptyString, consts.StartDateUsage)
	_syncEndDate              = SyncCmd.Flag.String("end_date", consts.EmptyString, consts.EndDateUsage)

	validSyncActions = []string{
		consts.LastActivities, consts.Playback, consts.RemovePlayback, consts.GetCollection, consts.GetMinimalCollection,
		consts.GetUpNext, consts.GetUpNextNitro, consts.GetWatchedProgress,
		consts.AddToCollection, consts.RemoveFromCollection, consts.GetWatched,
		consts.GetHistory, consts.AddToHistory, consts.RemoveFromHistory,
		consts.GetRatings, consts.AddToRatings, consts.RemoveFromRatings,
		consts.GetWatchlist, consts.UpdateWatchlist, consts.AddToWatchlist,
		consts.RemoveFromWatchlist, consts.ReorderWatchlist, consts.UpdateWatchlistItem,
		consts.GetFavorites, consts.UpdateFavorites, consts.AddToFavorites,
		consts.RemoveFromFavorites, consts.ReorderFavorites, consts.UpdateFavoriteItem}
)

// SyncCmd returns movies and episodes that a user has watched, sorted by most recent.
var SyncCmd = &Command{
	Name:    "sync",
	Usage:   "",
	Summary: "Sync data useful for mediacenters: activities, playbacks, collections, ratings, watchlists, favorites.",
	Help:    `sync command`,
}

func syncFunc(cmd *Command, _ ...string) error {
	cmd.UpdateSyncFlagsValues()
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	options.Type = syncPlaybackType(options, cmd.flagIsSet("t"))

	err := cmd.ValidSort(options)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}
	syncProgressSort(options, cmd.flagIsSet("sort_by"), cmd.flagIsSet("sort_how"))
	var handler handlers.SyncHandler
	allHandlers := map[string]handlers.Handler{
		consts.LastActivities:       handlers.SyncLastActivitiesHandler{},
		consts.Playback:             handlers.SyncPlaybackHandler{},
		consts.RemovePlayback:       handlers.SyncRemovePlaybackHandler{},
		consts.GetCollection:        handlers.SyncGetCollectionHandler{},
		consts.GetMinimalCollection: handlers.SyncGetMinimalCollectionHandler{},
		consts.GetUpNext:            handlers.SyncGetUpNextHandler{},
		consts.GetUpNextNitro:       handlers.SyncGetUpNextNitroHandler{},
		consts.GetWatchedProgress:   handlers.SyncGetWatchedProgressHandler{},
		consts.AddToCollection:      handlers.SyncAddToCollectionHandler{},
		consts.RemoveFromCollection: handlers.SyncRemoveFromCollectionHandler{},
		consts.GetWatched:           handlers.SyncGetWatchedHandler{},
		consts.GetHistory:           handlers.SyncGetHistoryHandler{},
		consts.AddToHistory:         handlers.SyncAddToHistoryHandler{},
		consts.RemoveFromHistory:    handlers.SyncRemoveFromHistoryHandler{},
		consts.GetRatings:           handlers.SyncGetRatingsHandler{},
		consts.AddToRatings:         handlers.SyncAddToRatingsHandler{},
		consts.RemoveFromRatings:    handlers.SyncRemoveFromRatingsHandler{},
		consts.GetWatchlist:         handlers.SyncGetWatchlistHandler{},
		consts.UpdateWatchlist:      handlers.SyncUpdateWatchlistHandler{},
		consts.AddToWatchlist:       handlers.SyncAddToWatchlistHandler{},
		consts.RemoveFromWatchlist:  handlers.SyncRemoveFromWatchlistHandler{},
		consts.ReorderWatchlist:     handlers.SyncReorderWatchlistHandler{},
		consts.UpdateWatchlistItem:  handlers.SyncUpdateWatchlistItemHandler{},
		consts.GetFavorites:         handlers.SyncGetFavoritesHandler{},
		consts.UpdateFavorites:      handlers.SyncUpdateFavoritesHandler{},
		consts.AddToFavorites:       handlers.SyncAddToFavoritesHandler{},
		consts.RemoveFromFavorites:  handlers.SyncRemoveFromFavoritesHandler{},
		consts.ReorderFavorites:     handlers.SyncReorderFavoritesHandler{},
		consts.UpdateFavoriteItem:   handlers.SyncUpdateFavoriteItemHandler{},
	}
	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validSyncActions)
		return unknownActionError(cmd.Name, options.Action)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

var (
	syncDumpTemplate = ``
)

func init() {
	SyncCmd.Run = syncFunc
}

// syncPlaybackType makes playback without -t cover movies and episodes (the untyped sync/playback route)
func syncPlaybackType(options *str.Options, typeSet bool) string {
	if options.Action == consts.Playback && !typeSet {
		return consts.ActionTypeAll
	}

	return options.Type
}

// syncProgressSort sends sort_by and sort_how for up next (nitro) and watched progress only when set on the command line,
// so the global defaults (rank, asc) do not override the API's own order
func syncProgressSort(options *str.Options, sortBySet bool, sortHowSet bool) {
	if options.Action != consts.GetUpNext && options.Action != consts.GetUpNextNitro && options.Action != consts.GetWatchedProgress {
		return
	}
	if !sortBySet {
		options.SortBy = consts.EmptyString
	}
	if !sortHowSet {
		options.SortHow = consts.EmptyString
	}
}
