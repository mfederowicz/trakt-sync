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
	connections, _ := h.common.FetchUserConnections(client, options)
	movie, _ := h.common.FetchMovie(client, options)
	c := new(str.CheckIn)
	c.Movie = movie
	if len(options.Msg) > consts.ZeroValue {
		c.Message = &options.Msg
	}
	c.Sharing = new(str.Sharing)
	c.Sharing.Tumblr = connections.Tumblr
	c.Sharing.Twitter = connections.Twitter
	c.Sharing.Mastodon = connections.Mastodon

	result, resp, err := h.common.Checkin(client, c)
	if err != nil {
		return fmt.Errorf("checkin error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, movie checkin number:%d \n", result.ID)
	}

	return nil
}
