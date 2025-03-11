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

// GetCommentUserLikes Returns all users who liked a comment.
// API docs: https://trakt.docs.apiary.io/#reference/comments/item/get-all-users-who-liked-a-comment 
func (c *CommentsService) GetCommentUserLikes(ctx context.Context, id *int, opts *uri.ListOptions) ([]*str.CommentUserLike, *str.Response, error) {
	var url = fmt.Sprintf("comments/%d/likes", *id)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch likes url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CommentUserLike{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch likes err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// LikeComment Votes help determine popular comments. Only one like is allowed per comment per user.
//
// API docs: https://trakt.docs.apiary.io/#reference/comments/like/like-a-comment 
func (c *CommentsService) LikeComment(ctx context.Context, id *int) (*str.Response, error) {
	var url = fmt.Sprintf("comments/%d/like", *id)
	printer.Println("send like for single comment:" + url)
	req, err := c.client.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveLikeComment Remove a like on a comment.
//
// API docs: https://trakt.docs.apiary.io/#reference/comments/like/remove-like-on-a-comment 
func (c *CommentsService) RemoveLikeComment(ctx context.Context, id *int) (*str.Response, error) {
	var url = fmt.Sprintf("comments/%d/like", *id)
	printer.Println("remove like for single comment:" + url)
	req, err := c.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
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

// GetTrendingComments Returns all comments with the most likes and replies over the last 7 days.
// API docs: https://trakt.docs.apiary.io/#reference/comments/trending/get-trending-comments 
func (c *CommentsService) GetTrendingComments(ctx context.Context, contentType *string, strType *string, opts *uri.ListOptions) ([]*str.CommentItem, *str.Response, error) {
	var url = fmt.Sprintf("comments/trending/%s/%s", *contentType, *strType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch trending url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CommentItem{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch trending err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetRecentComments Returns the most recently written comments across all of Trakt.
// API docs: https://trakt.docs.apiary.io/#reference/comments/recent/get-recently-created-comments
func (c *CommentsService) GetRecentComments(ctx context.Context, contentType *string, strType *string, opts *uri.ListOptions) ([]*str.CommentItem, *str.Response, error) {
	var url = fmt.Sprintf("comments/recent/%s/%s", *contentType, *strType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch recent url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CommentItem{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch recent err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

// GetUpdatedComments Returns the most recently updated comments across all of Trakt.
// API docs: https://trakt.docs.apiary.io/#reference/comments/updates/get-recently-updated-comments
func (c *CommentsService) GetUpdatedComments(ctx context.Context, contentType *string, strType *string, opts *uri.ListOptions) ([]*str.CommentItem, *str.Response, error) {
	var url = fmt.Sprintf("comments/updates/%s/%s", *contentType, *strType)
	url, err := uri.AddQuery(url, opts)
	if err != nil {
		return nil, nil, err
	}
	printer.Println("fetch updated url:" + url)
	req, err := c.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	list := []*str.CommentItem{}
	resp, err := c.client.Do(ctx, req, &list)

	if err != nil {
		printer.Println("fetch updated err:" + err.Error())
		return nil, resp, err
	}

	return list, resp, nil
}

