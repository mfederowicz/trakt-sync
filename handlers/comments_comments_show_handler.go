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

// CommentsCommentsShowHandler struct for handler
type CommentsCommentsShowHandler struct{ common CommonLogic }

// Handle to handle comments: show type
func (h CommentsCommentsShowHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}

	connections, _ := h.common.FetchUserConnections(client, options)
	show, err := h.common.FetchShow(client, options)
	if err != nil {
		return fmt.Errorf("fetch show error:%w", err)
	}

	c := new(str.Comment)
	c.Show = show
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
		printer.Printf("result: success, show comment number:%d \n", result.ID)
	}

	return nil
}
