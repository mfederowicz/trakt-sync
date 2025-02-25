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

	connections, _ := h.common.FetchUserConnections(client, options)
	show, err := h.common.FetchShow(client, options)
	if err != nil {
		return fmt.Errorf("fetch show error:%w", err)
	}

	if len(options.EpisodeCode) > consts.ZeroValue {
		return h.CreateCheckinForEpisodeCode(options, client, connections, show)
	}

	if options.EpisodeAbs > consts.ZeroValue {
		return h.CreateCheckinForEpisodeAbs(options, client, connections, show)
	}

	return nil
}

// CreateCheckinForEpisodeCode to handle checkin: episode code
func (h CheckinShowEpisodeHandler) CreateCheckinForEpisodeCode(options *str.Options, client *internal.Client, connections *str.Connections, show *str.Show) error {
	c := new(str.CheckIn)
	season, number, err := h.common.CheckSeasonNumber(options.EpisodeCode)
	if err != nil {
		return fmt.Errorf("check episode error:%w", err)
	}

	c.Show = new(str.Show)
	c.Show = show
	c.Episode = new(str.Episode)
	c.Episode.Season = season
	c.Episode.Number = number
	if len(options.Msg) > consts.ZeroValue {
		c.Message = &options.Msg
	}
	c.Sharing = new(str.Sharing)
	c.Sharing.Tumblr = connections.Tumblr
	c.Sharing.Twitter = connections.Twitter
	c.Sharing.Mastodon = connections.Mastodon

	result, resp, err := h.common.Checkin(client, c)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, season:%d, episode:%d", *show.Title, *season, *number)
	}

	if resp.StatusCode == http.StatusConflict {
		return fmt.Errorf("checkin for show:%s, season:%d, episode:%d exists, expires:%s", *show.Title, *season, *number, result.Expires.Local())
	}

	if err != nil {
		return printer.Errorf("checkin error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode code checkin number:%d \n", result.ID)
	}

	return nil
}

// CreateCheckinForEpisodeAbs to handle checkin: episode abs
func (h CheckinShowEpisodeHandler) CreateCheckinForEpisodeAbs(options *str.Options, client *internal.Client, connections *str.Connections, show *str.Show) error {
	c := new(str.CheckIn)
	c.Show = new(str.Show)
	c.Show = show
	c.Episode = new(str.Episode)
	c.Episode.NumberAbs = &options.EpisodeAbs
	if len(options.Msg) > consts.ZeroValue {
		c.Message = &options.Msg
	}
	c.Sharing = new(str.Sharing)
	c.Sharing.Tumblr = connections.Tumblr
	c.Sharing.Twitter = connections.Twitter
	c.Sharing.Mastodon = connections.Mastodon
	result, resp, err := h.common.Checkin(client, c)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found episode for show:%s, episode_abs:%d", *show.Title, options.EpisodeAbs)
	}

	if resp.StatusCode == http.StatusConflict {
		return fmt.Errorf("checkin for show:%s, episode_abs:%d exists, expires:%s", *show.Title, options.EpisodeAbs, result.Expires.Local())
	}

	if err != nil {
		return printer.Errorf("checkin error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode abs checkin number:%d \n", result.ID)
	}

	return nil
}
