// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// CheckinEpisodeHandler struct for handler
type CheckinEpisodeHandler struct {
	common CommonLogic
}

// Handle to handle checkin: episode action
func (h CheckinEpisodeHandler) Handle(options *str.Options, client *trakt.Client) error {
	if options.TraktID == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	checkin, err := h.common.CreateCheckin(client, options)
	if err != nil {
		return printer.Errorf(consts.CheckinError, err)
	}

	result, resp, err := h.common.Checkin(client, checkin, options)
	if err != nil {
		return printer.Errorf(consts.CheckinError, err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode checkin number:%d \n", result.ID)
	}

	return nil
}
