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
	"github.com/mfederowicz/trakt-sync/writer"
)

// CommentsReactionsSummaryHandler struct for handler
type CommentsReactionsSummaryHandler struct{}

// Handle to handle comments: reactions_summary action
func (CommentsReactionsSummaryHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.CommentID == consts.ZeroValue {
		return errors.New(consts.EmptyCommentIDMsg)
	}

	printer.Println("Get reaction totals for a comment.")
	commentID := options.CommentID
	result, _, err := client.Comments.GetCommentReactionsSummary(client.BuildCtxFromOptions(options), &commentID)
	if err != nil {
		return fmt.Errorf("fetch reactions summary error: %w", err)
	}

	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal reactions summary error: %w", err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
