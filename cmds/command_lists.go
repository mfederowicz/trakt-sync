// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_listsAction    = ListsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_listTraktID    = ListsCmd.Flag.String("trakt_id", cfg.DefaultConfig().InternalID, consts.ListIDUsage)
	_listInternalID = ListsCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.ListIDUsage)
	_listSort       = ListsCmd.Flag.String("s", cfg.DefaultConfig().CommentsSort, consts.ListCommentSortUsage)
	_listLikeRemove = ListsCmd.Flag.Bool("remove", cfg.DefaultConfig().Remove, consts.ListLikeRemoveUsage)
	_listType       = ListsCmd.Flag.String("t", cfg.DefaultConfig().ListsType, consts.ListTypeUsage)
	_listReason     = ListsCmd.Flag.String("r", cfg.DefaultConfig().Reason, consts.ReasonUsage)
	_listMessage    = ListsCmd.Flag.String("message", cfg.DefaultConfig().Msg, consts.ReportMsgUsage)
)

// ListsCmd returns movies and episodes that a user has watched, sorted by most recent.
var ListsCmd = &Command{
	Name:    "lists",
	Usage:   "",
	Summary: "Returns data about lists: trending, popular, list, likes, like, items, comments, report.",
	Help:    `lists command`,
}

func listsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	var handler handlers.ListsHandler
	allHandlers := map[string]handlers.Handler{
		"trending": handlers.ListsTrendingHandler{},
		"popular":  handlers.ListsPopularHandler{},
		"list":     handlers.ListsListHandler{},
		"likes":    handlers.ListsLikesHandler{},
		"like":     handlers.ListsLikeHandler{},
		"items":    handlers.ListsItemsHandler{},
		"comments": handlers.ListsCommentsHandler{},
		"report":   handlers.ListsReportHandler{},
	}

	handler, err := cmd.common.GetHandlerForMap(options.Action, allHandlers)

	validActions = []string{"trending", "popular", "list", "likes", "like", "items", "comments", "report"}
	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, validActions)
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

var (
	listsDumpTemplate = ``
)

func init() {
	ListsCmd.Run = listsFunc
}
