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

// ShowsReportHandler struct for handler
type ShowsReportHandler struct{ common CommonLogic }

// Handle to handle shows: report action
func (h ShowsReportHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}
	if len(options.Reason) == consts.ZeroValue {
		return errors.New(consts.EmptyReasonMsg)
	}
	if err := h.common.ValidReason(options); err != nil {
		return err
	}

	report := &str.ShowReport{Reason: &options.Reason}
	if len(options.Msg) > consts.ZeroValue {
		report.Message = &options.Msg
	}

	if _, err := client.Shows.ReportShow(client.BuildCtxFromOptions(options), &options.InternalID, report); err != nil {
		return fmt.Errorf("report error: %w", err)
	}

	printer.Printf("reported show %s\n", options.InternalID)
	return nil
}
