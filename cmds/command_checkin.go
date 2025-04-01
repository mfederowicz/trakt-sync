// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_checkinAction      = CheckinCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_checkinTraktID     = CheckinCmd.Flag.Int("trakt_id", cfg.DefaultConfig().TraktID, consts.TraktIDUsage)
	_checkinEpisodeAbs  = CheckinCmd.Flag.Int("episode_abs", cfg.DefaultConfig().EpisodeAbs, consts.EpisodeAbsUsage)
	_checkinEpisodeCode = CheckinCmd.Flag.String("episode_code", cfg.DefaultConfig().EpisodeCode, consts.EpisodeCodeUsage)
	_checkinMsg         = CheckinCmd.Flag.String("msg", cfg.DefaultConfig().Msg, consts.CheckInMsgUsage)
	_checkinDelete      = CheckinCmd.Flag.Bool("delete", cfg.DefaultConfig().Delete, consts.DeleteUsage)
)

// CheckinCmd create or delete active checkins.
var CheckinCmd = &Command{
	Name:    "checkin",
	Usage:   "",
	Summary: "Checkin movie,episode,show_episode,delete",
	Help:    `checkin command`,
}

func checkinFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	var handler handlers.CheckinHandler
	allHandlers := map[string]handlers.Handler{
		"movie":        handlers.CheckinMovieHandler{},
		"episode":      handlers.CheckinEpisodeHandler{},
		"show_episode": handlers.CheckinShowEpisodeHandler{},
		"delete":       handlers.CheckinDeleteHandler{},
	}

	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	validActions = []string{"movie", "episode", "show_episode", "delete"}
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
	checkinDumpTemplate = ``
)

func init() {
	CheckinCmd.Run = checkinFunc
}
