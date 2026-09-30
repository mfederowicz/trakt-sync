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

// CommentsReportHandler struct for handler
type CommentsReportHandler struct{ common CommonLogic }

// Handle to handle comments: report action
func (h CommentsReportHandler) Handle(options *str.Options, client *trakt.Client) error {
	if options.CommentID == consts.ZeroValue {
		return errors.New(consts.EmptyCommentIDMsg)
	}
	if len(options.Reason) == consts.ZeroValue {
		return errors.New(consts.EmptyReasonMsg)
	}
	if err := h.common.ValidReason(options); err != nil {
		return err
	}

	report := &str.CommentReport{Reason: &options.Reason}
	if len(options.Msg) > consts.ZeroValue {
		report.Message = &options.Msg
	}

	commentID := options.CommentID
	if _, err := client.Comments.ReportComment(cli.ContextFromOptions(options), &commentID, report); err != nil {
		return fmt.Errorf("report error: %w", err)
	}

	printer.Printf("reported comment %d\n", commentID)
	return nil
}
