// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// CommentsLikeHandler struct for handler
type CommentsLikeHandler struct{}

// Handle to handle comments: like action
func (h CommentsLikeHandler) Handle(options *str.Options, client *trakt.Client) error {
	if options.CommentID == consts.ZeroValue {
		return errors.New(consts.EmptyCommentIDMsg)
	}

	resp, err := h.likeSingleComment(client, options)
	if resp == nil {
		return fmt.Errorf("like comment error: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found comment for:%d", options.CommentID)
	}

	if err != nil {
		return fmt.Errorf("like comment error: %w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		printer.Print("result: success \n")
	}

	return nil
}

func (CommentsLikeHandler) likeSingleComment(client *trakt.Client, options *str.Options) (*str.Response, error) {
	commentID := options.CommentID

	if !options.Remove {
		resp, err := client.Comments.LikeComment(
			cli.ContextFromOptions(options),
			commentID,
		)
		return resp, err
	}

	resp, err := client.Comments.RemoveLikeComment(
		cli.ContextFromOptions(options),
		commentID,
	)

	return resp, err
}
