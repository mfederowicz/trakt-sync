// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_younifyAction    = YounifyCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_younifyServiceID = YounifyCmd.Flag.String("service_id", consts.EmptyString, consts.ServiceIDUsage)
	_younifyReturnURL = YounifyCmd.Flag.String("return_url", consts.DefaultReturnURL, consts.ReturnURLUsage)
	_younifyAllData   = YounifyCmd.Flag.Bool("all_data", false, consts.AllDataUsage)
)

// YounifyCmd manages streaming service connections (younify).
var YounifyCmd = &Command{
	Name:    "younify",
	Usage:   "",
	Summary: "Streaming service connections: connections, connect, refresh, disconnect.",
	Help:    `younify command`,
}

func younifyFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)
	var handler handlers.YounifyHandler
	allHandlers := map[string]handlers.Handler{
		consts.Connections: handlers.YounifyConnectionsHandler{},
		consts.Connect:     handlers.YounifyConnectHandler{},
		consts.Refresh:     handlers.YounifyRefreshHandler{},
		consts.Disconnect:  handlers.YounifyDisconnectHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.Connections, consts.Connect, consts.Refresh, consts.Disconnect})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	YounifyCmd.Run = younifyFunc
}
