// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// ShowsResetShowProgressHandler struct for handler
type ShowsResetShowProgressHandler struct{ common CommonLogic }

// Handle to handle shows: reset_show_progress action
func (m ShowsResetShowProgressHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Reset a show's progress when the user started re-watching the show.")
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyShowIDMsg)
	}

	if options.Delete {
		return m.handleUndoResetShowProgress(options, client)
	}

	return m.handleResetShowProgress(options, client)
}

// HandleDelete process delete
func (ShowsResetShowProgressHandler) handleUndoResetShowProgress(options *str.Options, client *internal.Client) error {
	resp, err := client.Shows.UndoResetShowProgress(client.BuildCtxFromOptions(options), &options.InternalID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if resp.StatusCode == http.StatusUpgradeRequired {
		return cli.HandleUpgrade(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		printer.Printf("result: success, undo reset show progress for:%s \n", options.InternalID)
	}
	return nil
}

// HandleModifyNotes handle modify exiting notes
func (m ShowsResetShowProgressHandler) handleResetShowProgress(options *str.Options, client *internal.Client) error {
	showProgress := new(str.WatchedProgress)

	if len(options.ResetAt) > consts.ZeroValue {
		showProgress.ResetAt = m.common.ToTimestamp(options.ResetAt)
	}

	result, resp, err := client.Shows.ResetShowProgress(client.BuildCtxFromOptions(options), &options.InternalID, showProgress)
	if err != nil {
		return fmt.Errorf("reset progress error:%w", err)
	}

	if resp.StatusCode == http.StatusUpgradeRequired {
		return cli.HandleUpgrade(resp)
	}

	if resp.StatusCode == http.StatusOK {
		printer.Printf("result: success, reset show progress:%v \n", result.ResetAt)
	}

	return nil
}
