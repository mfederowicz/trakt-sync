// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// UsersCommentReactionsHandler struct for handler
type UsersCommentReactionsHandler struct{}

// Handle to handle users: comment_reactions action
func (UsersCommentReactionsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns comments you have reacted to.")
	result, err := fetchAllPages(client, options, consts.DefaultPage, func(opts *uri.ListOptions) ([]*str.CommentReaction, *str.Response, error) {
		return client.Users.GetCommentReactions(client.BuildCtxFromOptions(options), opts)
	})
	if err != nil {
		return fmt.Errorf("fetch comment reactions error: %w", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	return writeResult(options, result)
}
