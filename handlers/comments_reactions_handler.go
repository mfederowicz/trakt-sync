// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// CommentsReactionsHandler struct for handler
type CommentsReactionsHandler struct{}

// Handle to handle comments: reactions action
func (CommentsReactionsHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.CommentID == consts.ZeroValue {
		return errors.New(consts.EmptyCommentIDMsg)
	}

	printer.Println("Get all reactions on a comment.")
	commentID := options.CommentID
	result, err := fetchAllPages(client, options, consts.DefaultPage, func(opts *uri.ListOptions) ([]*str.CommentReaction, *str.Response, error) {
		return client.Comments.GetCommentReactions(client.BuildCtxFromOptions(options), &commentID, opts)
	})
	if err != nil {
		return fmt.Errorf("fetch reactions error: %w", err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal reactions error: %w", err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
