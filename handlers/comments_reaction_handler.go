// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// CommentsReactionHandler struct for handler
type CommentsReactionHandler struct{ common CommonLogic }

// Handle to handle comments: reaction action, adds a reaction or removes it with -remove
func (h CommentsReactionHandler) Handle(options *str.Options, client *trakt.Client) error {
	if options.CommentID == consts.ZeroValue {
		return errors.New(consts.EmptyCommentIDMsg)
	}
	if len(options.Reaction) == consts.ZeroValue {
		return errors.New(consts.EmptyReactionMsg)
	}
	if err := h.common.ValidReaction(options); err != nil {
		return err
	}

	commentID := options.CommentID
	ctx := cli.ContextFromOptions(options)
	if options.Remove {
		if _, err := client.Comments.RemoveCommentReaction(ctx, &commentID, &options.Reaction); err != nil {
			return fmt.Errorf("remove reaction error: %w", err)
		}
		printer.Printf("removed reaction %s from comment %d\n", options.Reaction, commentID)
		return nil
	}

	if _, err := client.Comments.AddCommentReaction(ctx, &commentID, &options.Reaction); err != nil {
		return fmt.Errorf("add reaction error: %w", err)
	}
	printer.Printf("added reaction %s to comment %d\n", options.Reaction, commentID)
	return nil
}
