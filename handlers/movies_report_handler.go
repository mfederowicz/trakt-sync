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

// MoviesReportHandler struct for handler
type MoviesReportHandler struct{ common CommonLogic }

// Handle to handle movies: report action
func (h MoviesReportHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}
	if len(options.Reason) == consts.ZeroValue {
		return errors.New(consts.EmptyReasonMsg)
	}
	if err := h.common.ValidReason(options); err != nil {
		return err
	}

	report := &str.MovieReport{Reason: &options.Reason}
	if len(options.Msg) > consts.ZeroValue {
		report.Message = &options.Msg
	}

	if _, err := client.Movies.ReportMovie(cli.ContextFromOptions(options), &options.InternalID, report); err != nil {
		return fmt.Errorf("report error: %w", err)
	}

	printer.Printf("reported movie %s\n", options.InternalID)
	return nil
}
