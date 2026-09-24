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

// EpisodesReportHandler struct for handler
type EpisodesReportHandler struct{ common CommonLogic }

// Handle to handle episodes: report action
func (h EpisodesReportHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}
	if options.Episode < consts.FirstEpisodeNumber {
		return errors.New(consts.EmptyEpisodeMsg)
	}
	if len(options.Reason) == consts.ZeroValue {
		return errors.New(consts.EmptyReasonMsg)
	}
	if err := h.common.ValidReason(options); err != nil {
		return err
	}

	report := &str.EpisodeReport{Reason: &options.Reason}
	if len(options.Msg) > consts.ZeroValue {
		report.Message = &options.Msg
	}

	if _, err := client.Shows.ReportEpisode(client.BuildCtxFromOptions(options), &options.InternalID, &options.Season, &options.Episode, report); err != nil {
		return fmt.Errorf("report error: %w", err)
	}

	printer.Printf("reported episode %dx%d of show %s\n", options.Season, options.Episode, options.InternalID)
	return nil
}
