// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_mediaAction = MediaCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
)

// MediaCmd returns trending, popular and anticipated movies and shows together.
var MediaCmd = &Command{
	Name:    "media",
	Usage:   "",
	Summary: "Returns movies and shows together: trending, popular, anticipated.",
	Help:    `media command`,
}

func mediaFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.MediaHandler
	allHandlers := map[string]handlers.Handler{
		consts.Trending:    handlers.MediaTrendingHandler{},
		consts.Popular:     handlers.MediaPopularHandler{},
		consts.Anticipated: handlers.MediaAnticipatedHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.Trending, consts.Popular, consts.Anticipated})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	MediaCmd.Run = mediaFunc
}
