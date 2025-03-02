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

// CommentsCommentsListHandler struct for handler
type CommentsCommentsListHandler struct{ common CommonLogic }

// Handle to handle comments: list type
func (h CommentsCommentsListHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.TraktID == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	connections, _ := h.common.FetchUserConnections(client, options)
	list, err := h.common.FetchList(client, options)
	if err != nil {
		return fmt.Errorf("fetch list error:%w", err)
	}

	c := new(str.Comment)
	c.List = list
	if len(options.Comment) > consts.ZeroValue {
		c.Comment = &options.Comment
	}
	c.Spoiler = &options.Spoiler
	c.Sharing = new(str.Sharing)
	c.Sharing.Tumblr = connections.Tumblr
	c.Sharing.Twitter = connections.Twitter
	c.Sharing.Mastodon = connections.Mastodon

	result, resp, err := h.common.Comment(client, c)
	if err != nil {
		return fmt.Errorf("comment error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, list comment number:%d \n", result.ID)
	}
	return nil
}
