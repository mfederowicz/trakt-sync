// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_syncAction     = SyncCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_syncStartAt    = SyncCmd.Flag.String("start_at", cfg.DefaultConfig().StartAt, consts.StartAtUsage)
	_syncEndAt      = SyncCmd.Flag.String("end_at", cfg.DefaultConfig().EndAt, consts.EndAtUsage)
	_syncPlaybackID = SyncCmd.Flag.Int("playback_id", cfg.DefaultConfig().PlaybackID, consts.PlaybackIDUsage)

	validSyncActions = []string{"last_activities"}
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

	var handler handlers.SyncHandler
	allHandlers := map[string]handlers.Handler{
		"last_activities": handlers.SyncLastActivitiesHandler{},
		"playback":        handlers.SyncPlaybackHandler{},
		"remove_playback": handlers.SyncRemovePlaybackHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

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
