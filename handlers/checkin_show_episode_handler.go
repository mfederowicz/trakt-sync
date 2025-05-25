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

// CheckinShowEpisodeHandler struct for handler
type CheckinShowEpisodeHandler struct {
	common CommonLogic
}

// Handle to handle checkin: episode action
func (h CheckinShowEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.TraktID == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	if options.EpisodeAbs > consts.ZeroValue && len(options.EpisodeCode) > consts.ZeroValue {
		return errors.New("only episode_abs or episode_code at time")
	}

	if len(options.EpisodeCode) > consts.ZeroValue {
		return h.CreateCheckinForEpisodeCode(options, client)
	}

	if options.EpisodeAbs > consts.ZeroValue {
		return h.CreateCheckinForEpisodeAbs(options, client)
	}

	return nil
}

// CreateCheckinForEpisodeCode to handle checkin: episode code
func (h CheckinShowEpisodeHandler) CreateCheckinForEpisodeCode(options *str.Options, client *internal.Client) error {
	checkin, err := h.common.CreateCheckin(client, options)
	if err != nil {
		return printer.Errorf(consts.CheckinError, err)
	}

	result, resp, err := h.common.Checkin(client, checkin, options)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, season:%d, episode:%d", *checkin.Show.Title, *checkin.Episode.Season, *checkin.Episode.Number)
	}

	if resp.StatusCode == http.StatusConflict {
		return fmt.Errorf("checkin for show:%s, season:%d, episode:%d exists, expires:%s", *checkin.Show.Title, *checkin.Episode.Season, *checkin.Episode.Number, result.Expires.Local())
	}

	if err != nil {
		return printer.Errorf(consts.CheckinError, err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode code checkin number:%d \n", result.ID)
	}

	return nil
}

// CreateCheckinForEpisodeAbs to handle checkin: episode abs
func (h CheckinShowEpisodeHandler) CreateCheckinForEpisodeAbs(options *str.Options, client *internal.Client) error {
	checkin, err := h.common.CreateCheckin(client, options)
	if err != nil {
		return printer.Errorf(consts.CheckinError, err)
	}

	result, resp, err := h.common.Checkin(client, checkin, options)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, episode_abs:%d", *checkin.Show.Title, options.EpisodeAbs)
	}

	if resp.StatusCode == http.StatusConflict {
		return fmt.Errorf("checkin for show:%s, episode_abs:%d exists, expires:%s", *checkin.Show.Title, options.EpisodeAbs, result.Expires.Local())
	}

	if err != nil {
		return printer.Errorf(consts.CheckinError, err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode abs checkin number:%d \n", result.ID)
	}

	return nil
}
