// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsService  handles communication with the comments related
// methods of the Trakt API.
type CommentsService Service

// PostAComment Add a new comment to a movie, show, season, episode, or list.
//
// API docs:https://trakt.docs.apiary.io/#reference/comments/comments/post-a-comment
func (c *CommentsService) PostAComment(ctx context.Context, comment *str.Comment) (*str.Comment, *str.Response, error) {
	var url = "comments"
	printer.Println("create new comment")
	req, err := c.client.NewRequest(http.MethodPost, url, comment)
	if err != nil {
		return nil, nil, err
	}

	com := new(str.Comment)
	resp, err := c.client.Do(ctx, req, com)
	if err != nil {
		return nil, nil, errors.Join(resp.Errors.GetComments())
	}

	return com, resp, nil
}
