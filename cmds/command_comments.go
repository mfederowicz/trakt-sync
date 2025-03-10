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
	_commentsAction         = CommentsCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_commentsTraktID        = CommentsCmd.Flag.Int("trakt_id", cfg.DefaultConfig().TraktID, consts.TraktIDUsage)
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
	var handler handlers.CommentsHandler
	switch options.Action {
	case "comments":
		handler = handlers.CommentsCommentsHandler{}
	case "comment":
		handler = handlers.CommentsCommentHandler{}
	case "replies":
		handler = handlers.CommentsRepliesHandler{}
	case "item":
		handler = handlers.CommentsItemHandler{}
	case "likes":
		handler = handlers.CommentsLikesHandler{}
	case "like":
		handler = handlers.CommentsLikeHandler{}
	case "trending":
		err := cmd.ValidType(options)
		if err != nil {
			return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
		}

		handler = handlers.CommentsTrendingHandler{}
	case "recent":
		handler = handlers.CommentsRecentHandler{}
	case "updates":
		handler = handlers.CommentsUpdatesHandler{}

	default:
		printer.Println("possible actons: comments,comment,replies,item,likes,like,trending,recent,updates")
	}
	err := handler.Handle(options, client)
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
