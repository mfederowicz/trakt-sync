// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
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

// GetCommentItem Returns comment media item object.
// API docs: https://trakt.docs.apiary.io/#reference/comments/item/get-the-attached-media-item
func (c *CommentsService) GetCommentItem(ctx context.Context, id *int, opts *uri.ListOptions) (*str.CommentMediaItem, *str.Response, error) {
	var url = fmt.Sprintf("comments/%d/item", *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}

	printer.Println("fetch comment madia item url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.CommentMediaItem)
	resp, err := c.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("comment item not found with commentId:%d", *id)
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

// GetRepliesForComment Returns all replies for a comment.
// API docs: https://trakt.docs.apiary.io/#reference/comments/replies/get-replies-for-a-comment
func (c *CommentsService) GetRepliesForComment(ctx context.Context, opts *uri.ListOptions, id *int) ([]*str.Comment, *str.Response, error) {
	var url = fmt.Sprintf("comments/%d/replies", *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch replies url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.Comment{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch replies err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// ReplyAComment Add a new reply to an existing comment.
// API docs:https://trakt.docs.apiary.io/#reference/comments/replies/post-a-reply-for-a-comment
func (c *CommentsService) ReplyAComment(ctx context.Context, id *int, reply *str.Comment) (*str.Comment, *str.Response, error) {
	var url = fmt.Sprintf("comments/%d/replies", *id)
	printer.Println("reply comment")
	req, err := c.client.NewRequest(http.MethodPost, url, reply)
	if err != nil {
		return nil, nil, err
	}

	com := new(str.Comment)
	resp, err := c.client.Do(ctx, req, com)

	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("comment not found with commentId:%d", *id)
	}

	if err != nil {
		return nil, nil, err
	}

	return com, resp, nil
}
