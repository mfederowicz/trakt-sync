// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/printer"
)

var (
	_checkinAction  = CheckinCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_checkinTraktID = CheckinCmd.Flag.Int("trakt_id", cfg.DefaultConfig().TraktID, consts.TraktIDUsage)
	_checkinEpisodeAbs = CheckinCmd.Flag.Int("episode_abs", cfg.DefaultConfig().EpisodeAbs, consts.EpisodeAbsUsage)
	_checkinEpisodeCode = CheckinCmd.Flag.String("episode_code", cfg.DefaultConfig().EpisodeCode, consts.EpisodeCodeUsage)
	_checkinMsg = CheckinCmd.Flag.String("msg", cfg.DefaultConfig().Msg, consts.CheckInMsgUsage)
	_checkinDelete  = CheckinCmd.Flag.Bool("delete", cfg.DefaultConfig().Delete, consts.DeleteUsage)
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

	var handler handlers.ListsHandler
	switch options.Action {
	case "movie":
		handler = handlers.CheckinMovieHandler{}
	case "episode":
		handler = handlers.CheckinEpisodeHandler{}
	case "show_episode":
		handler = handlers.CheckinShowEpisodeHandler{}
	case "delete":
		handler = handlers.CheckinDeleteHandler{}
	default:
		printer.Println("possible actions: movie,episode,show_episode,delete")
	}
	err := handler.Handle(options, client)
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
