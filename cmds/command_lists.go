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
	_listsAction    = ListsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_listTraktID    = ListsCmd.Flag.Int("trakt_id", cfg.DefaultConfig().TraktID, consts.ListIDUsage)
)

// ListsCmd returns movies and episodes that a user has watched, sorted by most recent.
var ListsCmd = &Command{
	Name:    "lists",
	Usage:   "",
	Summary: "Returns data about lists: trending, popular, list, likes, like, items, comments.",
	Help:    `lists command`,
}

func listsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	var handler handlers.ListsHandler
	switch options.Action {
	case "trending":
		handler = handlers.ListsTrendingHandler{}
	case "popular":
		handler = handlers.ListsPopularHandler{}
	case "list":
		handler = handlers.ListsListHandler{}
	case "likes":
		handler = handlers.ListsLikesHandler{}
	case "like":
		handler = handlers.ListsLikeHandler{}
	case "items":
		handler = handlers.ListsItemsHandler{}
	case "comments":
		handler = handlers.ListsCommentsHandler{}
	default:
		printer.Println("possible actions: trending, popular, list, likes, like, items, comments")
	}
	err := handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	listsDumpTemplate = ``
)

func init() {
	ListsCmd.Run = listsFunc
}
