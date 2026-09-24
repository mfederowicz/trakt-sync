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

	validSyncActions = []string{
		"last_activities", "playback", "remove_playback", "get_collection", "get_minimal_collection",
		"get_up_next", "get_watched_progress",
		"add_to_collection", "remove_from_collection", "get_watched",
		"get_history", "add_to_history", "remove_from_history",
		"get_ratings", "add_to_ratings", "remove_from_ratings",
		"get_watchlist", "update_watchlist", "add_to_watchlist",
		"remove_from_watchlist", "reorder_watchlist", "update_watchlist_item",
		"get_favorites", "update_favorites", "add_to_favorites",
		"remove_from_favorites", "reorder_favorites", "update_favorite_item"}
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
		"last_activities":        handlers.SyncLastActivitiesHandler{},
		"playback":               handlers.SyncPlaybackHandler{},
		"remove_playback":        handlers.SyncRemovePlaybackHandler{},
		"get_collection":         handlers.SyncGetCollectionHandler{},
		"get_minimal_collection": handlers.SyncGetMinimalCollectionHandler{},
		"get_up_next":            handlers.SyncGetUpNextHandler{},
		"get_watched_progress":   handlers.SyncGetWatchedProgressHandler{},
		"add_to_collection":      handlers.SyncAddToCollectionHandler{},
		"remove_from_collection": handlers.SyncRemoveFromCollectionHandler{},
		"get_watched":            handlers.SyncGetWatchedHandler{},
		"get_history":            handlers.SyncGetHistoryHandler{},
		"add_to_history":         handlers.SyncAddToHistoryHandler{},
		"remove_from_history":    handlers.SyncRemoveFromHistoryHandler{},
		"get_ratings":            handlers.SyncGetRatingsHandler{},
		"add_to_ratings":         handlers.SyncAddToRatingsHandler{},
		"remove_from_ratings":    handlers.SyncRemoveFromRatingsHandler{},
		"get_watchlist":          handlers.SyncGetWatchlistHandler{},
		"update_watchlist":       handlers.SyncUpdateWatchlistHandler{},
		"add_to_watchlist":       handlers.SyncAddToWatchlistHandler{},
		"remove_from_watchlist":  handlers.SyncRemoveFromWatchlistHandler{},
		"reorder_watchlist":      handlers.SyncReorderWatchlistHandler{},
		"update_watchlist_item":  handlers.SyncUpdateWatchlistItemHandler{},
		"get_favorites":          handlers.SyncGetFavoritesHandler{},
		"update_favorites":       handlers.SyncUpdateFavoritesHandler{},
		"add_to_favorites":       handlers.SyncAddToFavoritesHandler{},
		"remove_from_favorites":  handlers.SyncRemoveFromFavoritesHandler{},
		"reorder_favorites":      handlers.SyncReorderFavoritesHandler{},
		"update_favorite_item":   handlers.SyncUpdateFavoriteItemHandler{},
	}
	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validSyncActions)
		return nil
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

// syncProgressSort sends sort_by and sort_how for up next and watched progress only when set on the command line,
// so the global defaults (rank, asc) do not override the API's own order
func syncProgressSort(options *str.Options, sortBySet bool, sortHowSet bool) {
	if options.Action != consts.GetUpNext && options.Action != consts.GetWatchedProgress {
		return
	}
	if !sortBySet {
		options.SortBy = consts.EmptyString
	}
	if !sortHowSet {
		options.SortHow = consts.EmptyString
	}
}
