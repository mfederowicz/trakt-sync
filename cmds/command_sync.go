// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
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

	validSyncActions = []string{
		"last_activities", "playback", "remove_playback", "get_collection",
		"add_to_collection", "remove_from_collection", "get_watched",
		"get_history", "add_to_history", "remove_from_history",
		"get_ratings", "add_to_ratings", "remove_from_ratings",
		"get_watchlist", "update_watchlist", "add_to_watchlist",
		"remove_from_watchlist", "reorder_watchlist", "update_watchlist_item",
		"get_favorites", "update_favorites", "add_to_favorites"}
)

// SyncCmd returns movies and episodes that a user has watched, sorted by most recent.
var SyncCmd = &Command{
	Name:    "sync",
	Usage:   "",
	Summary: "Syncing with trakt",
	Help:    `sync command`,
}

func syncFunc(cmd *Command, _ ...string) error {
	cmd.UpdateSyncFlagsValues()
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	err := cmd.ValidSort(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}
	var handler handlers.SyncHandler
	allHandlers := map[string]handlers.Handler{
		"last_activities":        handlers.SyncLastActivitiesHandler{},
		"playback":               handlers.SyncPlaybackHandler{},
		"remove_playback":        handlers.SyncRemovePlaybackHandler{},
		"get_collection":         handlers.SyncGetCollectionHandler{},
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
	}
	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validSyncActions)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	syncDumpTemplate = ``
)

func init() {
	SyncCmd.Run = syncFunc
}
