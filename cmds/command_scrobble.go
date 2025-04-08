// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_scrobbleAction      = ScrobbleCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_scrobbleType        = ScrobbleCmd.Flag.String("t", cfg.DefaultConfig().Type, consts.TypeUsage)
	_scrobbleInternalID  = ScrobbleCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.TraktIDUsage)
	_scrobbleProgress    = ScrobbleCmd.Flag.Float64("progress", cfg.DefaultConfig().Progress, consts.ProgressUsage)
	_scrobbleEpisodeAbs  = ScrobbleCmd.Flag.Int("episode_abs", cfg.DefaultConfig().EpisodeAbs, consts.EpisodeAbsUsage)
	_scrobbleEpisodeCode = ScrobbleCmd.Flag.String("episode_code", cfg.DefaultConfig().EpisodeCode, consts.EpisodeCodeUsage)
	_scrobbleDelete      = ScrobbleCmd.Flag.Bool("delete", cfg.DefaultConfig().Delete, consts.DeleteUsage)
)

// ScrobbleCmd start/pause/stop what is user watching.
var ScrobbleCmd = &Command{
	Name:    "scrobble",
	Usage:   "",
	Summary: "Scrobble for start/pause/stop movie,show,episode",
	Help:    `scrobble command`,
}

func scrobbleFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	var handler handlers.ScrobbleHandler
	allHandlers := map[string]handlers.Handler{
		"start": handlers.ScrobbleStartHandler{},
		"pause": handlers.ScrobblePauseHandler{},
		"stop":  handlers.ScrobbleStopHandler{},
	}

	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	validActions = []string{"start", "pause", "stop"}
	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validActions)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	scrobbleDumpTemplate = ``
)

func init() {
	ScrobbleCmd.Run = scrobbleFunc
}
