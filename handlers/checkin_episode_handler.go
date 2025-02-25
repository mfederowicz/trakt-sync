// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CheckinEpisodeHandler struct for handler
type CheckinEpisodeHandler struct {
	common CommonLogic
}

// Handle to handle checkin: episode action
func (h CheckinEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.TraktID == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	connections, _ := h.common.FetchUserConnections(client, options)
	episode, _ := h.common.FetchEpisode(client, options)
	c := new(str.CheckIn)
	c.Episode = new(str.Episode)
	c.Episode.IDs = new(str.IDs)
	c.Episode.IDs.Trakt = episode.IDs.Trakt
	if len(options.Msg) > consts.ZeroValue {
		c.Message = &options.Msg
	}
	c.Sharing = new(str.Sharing)
	c.Sharing.Tumblr = connections.Tumblr
	c.Sharing.Twitter = connections.Twitter
	c.Sharing.Mastodon = connections.Mastodon

	result, resp, err := h.common.Checkin(client, c)
	if err != nil {
		return printer.Errorf("checkin error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode checkin number:%d \n", result.ID)
	}


	return nil
}
