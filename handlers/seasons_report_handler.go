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

// SeasonsReportHandler struct for handler
type SeasonsReportHandler struct{ common CommonLogic }

// Handle to handle seasons: report action
func (h SeasonsReportHandler) Handle(options *str.Options, client *internal.Client) error {
	byID := len(options.ID) > consts.ZeroValue
	if !byID && len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}
	if len(options.Reason) == consts.ZeroValue {
		return errors.New(consts.EmptyReasonMsg)
	}
	if err := h.common.ValidReason(options); err != nil {
		return err
	}

	report := &str.SeasonReport{Reason: &options.Reason}
	if len(options.Msg) > consts.ZeroValue {
		report.Message = &options.Msg
	}

	if byID {
		if _, err := client.Seasons.ReportSeason(client.BuildCtxFromOptions(options), &options.ID, report); err != nil {
			return fmt.Errorf("report error: %w", err)
		}
		printer.Printf("reported season %s\n", options.ID)
		return nil
	}

	if _, err := client.Shows.ReportSeason(client.BuildCtxFromOptions(options), &options.InternalID, &options.Season, report); err != nil {
		return fmt.Errorf("report error: %w", err)
	}

	printer.Printf("reported season %d of show %s\n", options.Season, options.InternalID)
	return nil
}
