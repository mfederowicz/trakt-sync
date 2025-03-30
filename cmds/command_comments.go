// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_commentsAction         = CommentsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_commentsInternalID     = CommentsCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.TraktIDUsage)
	_commentsTraktID        = CommentsCmd.Flag.String("trakt_id", cfg.DefaultConfig().InternalID, consts.TraktIDUsage)
	_commentsCommentID      = CommentsCmd.Flag.Int("comment_id", cfg.DefaultConfig().CommentID, consts.CommentIDUsage)
	_commentsDelete         = CommentsCmd.Flag.Bool("delete", cfg.DefaultConfig().Delete, consts.DeleteUsage)
	_commentsRemove         = CommentsCmd.Flag.Bool("remove", cfg.DefaultConfig().Remove, consts.RemoveUsage)
	_commentsSpoiler        = CommentsCmd.Flag.Bool("spoiler", cfg.DefaultConfig().Spoiler, consts.SpoilerUsage)
	_commentsIncludeReplies = CommentsCmd.Flag.String("include_replies", cfg.DefaultConfig().IncludeReplies, consts.IncludeRepliesUsage)
	_commentsComment        = CommentsCmd.Flag.String("comment", cfg.DefaultConfig().Comment, consts.CommentUsage)
	_commentsCommentType    = CommentsCmd.Flag.String("comment_type", cfg.DefaultConfig().CommentType, consts.CommentTypeUsage)
	_commentsReply          = CommentsCmd.Flag.String("reply", cfg.DefaultConfig().Reply, consts.ReplyUsage)
)

// CommentsCmd manage all types of comments:movie, show, season, episode, or list.
var CommentsCmd = &Command{
	Name:    "comments",
	Usage:   "",
	Summary: "Comments comments,comment,replies,item,likes,like,trending,recent,updates",
	Help:    `comments command`,
}

func commentsFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	err := cmd.ValidModuleActionType(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	var handler handlers.CommentsHandler
	allHandlers := map[string]handlers.Handler{
		"comments": handlers.CommentsCommentsHandler{},
		"comment":  handlers.CommentsCommentHandler{},
		"replies":  handlers.CommentsRepliesHandler{},
		"item":     handlers.CommentsItemHandler{},
		"likes":    handlers.CommentsLikesHandler{},
		"like":     handlers.CommentsLikeHandler{},
		"trending": handlers.CommentsTrendingHandler{},
		"recent":   handlers.CommentsRecentHandler{},
		"updates":  handlers.CommentsUpdatesHandler{},
	}

	handler, err = cmd.common.GetHandlerForMap(options.Action, allHandlers)

	validActions = []string{
		"comments", "comment", "replies", "item", "likes", "like", 
		"trending", "recent", "updates",
	}
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
	commentsDumpTemplate = ``
)

func init() {
	CommentsCmd.Run = commentsFunc
}
