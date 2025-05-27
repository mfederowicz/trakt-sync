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

// CommentsCommentsSeasonHandler struct for handler
type CommentsCommentsSeasonHandler struct{ common CommonLogic }

// Handle to handle comments: season type
func (h CommentsCommentsSeasonHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.TraktID == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	connections, _ := h.common.FetchUserConnections(client, options)
	season, err := h.common.FetchSeason(client, options)
	if err != nil {
		return fmt.Errorf("fetch season error:%w", err)
	}

	c := new(str.Comment)
	c.Season = season
	if len(options.Comment) > consts.ZeroValue {
		c.Comment = &options.Comment
	}
	c.Spoiler = &options.Spoiler
	c.Sharing = new(str.Sharing)
	c.Sharing.Tumblr = connections.Tumblr
	c.Sharing.Twitter = connections.Twitter
	c.Sharing.Mastodon = connections.Mastodon

	result, resp, err := h.common.Comment(client, c, options)
	if err != nil {
		return fmt.Errorf("comment error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, season comment number:%d \n", result.ID)
	}

	return nil
}
