// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsReportHandler struct for handler
type CommentsReportHandler struct{ common CommonLogic }

// Handle to handle comments: report action
func (h CommentsReportHandler) Handle(options *str.Options, client *internal.Client) error {
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
	if _, err := client.Comments.ReportComment(client.BuildCtxFromOptions(options), &commentID, report); err != nil {
		return fmt.Errorf("report error: %w", err)
	}

	printer.Printf("reported comment %d\n", commentID)
	return nil
}
