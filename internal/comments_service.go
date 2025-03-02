// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"fmt"
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

// UpdateComment to update a single comment.
// API docs: https://trakt.docs.apiary.io/#reference/comments/comment/update-a-comment-or-reply
func (c *CommentsService) UpdateComment(ctx context.Context, id *int, comment *str.Comment) (*str.Comment, *str.Response, error) {
	var url = fmt.Sprintf("comments/%d", *id)
	printer.Println("update comment")
	req, err := c.client.NewRequest(http.MethodPut, url, comment)
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

// GetComment Returns comment object.
func (c *CommentsService) GetComment(ctx context.Context, id *int) (*str.Comment, *str.Response, error) {
	var url = fmt.Sprintf("comments/%d", *id)
	printer.Println("fetch comment url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.Comment)
	resp, err := c.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("comment not found with commentId:%d", *id)
	}

	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteComment to delete a single comment.
// API docs: https://trakt.docs.apiary.io/#reference/comments/comment/delete-a-comment-or-reply
func (c *CommentsService) DeleteComment(ctx context.Context, id *int) (*str.Response, error) {
	var url = fmt.Sprintf("comments/%d", *id)
	printer.Println("delete comment")
	req, err := c.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(ctx, req, nil)
	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("comment not found with commentId:%d", *id)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}
