// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CheckinMovieHandler struct for handler
type CheckinMovieHandler struct {
	common CommonLogic
}

// Handle to handle checkin: checkin action
func (h CheckinMovieHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.TraktID == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	checkin, err := h.common.CreateCheckin(client, options)
	if err != nil {
		return fmt.Errorf(consts.CheckinError, err)
	}
	
	result, resp, err := h.common.Checkin(client, checkin)
	if err != nil {
		return fmt.Errorf(consts.CheckinError, err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, movie checkin number:%d \n", result.ID)
	}

	return nil
}
