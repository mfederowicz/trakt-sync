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

// PeopleReportHandler struct for handler
type PeopleReportHandler struct{ common CommonLogic }

// Handle to handle people: report action
func (h PeopleReportHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.EmptyPersonIDMsg)
	}
	if len(options.Reason) == consts.ZeroValue {
		return errors.New(consts.EmptyReasonMsg)
	}
	if err := h.common.ValidReason(options); err != nil {
		return err
	}

	report := &str.PersonReport{Reason: &options.Reason}
	if len(options.Msg) > consts.ZeroValue {
		report.Message = &options.Msg
	}

	if _, err := client.People.ReportPerson(client.BuildCtxFromOptions(options), &options.ID, report); err != nil {
		return fmt.Errorf("report error: %w", err)
	}

	printer.Printf("reported person %s\n", options.ID)
	return nil
}
