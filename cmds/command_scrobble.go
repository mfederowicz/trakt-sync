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
		consts.Start: handlers.ScrobbleStartHandler{},
		consts.Pause: handlers.ScrobblePauseHandler{},
		consts.Stop:  handlers.ScrobbleStopHandler{},
	}

	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	validActions := []string{consts.Start, consts.Pause, consts.Stop}
	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validActions)
		return unknownActionError(cmd.Name, options.Action)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	ScrobbleCmd.Run = scrobbleFunc
}
