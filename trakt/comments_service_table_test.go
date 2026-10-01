// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestCommentsServiceRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	const id = 417
	const (
		comment  = `{"id":417,"comment":"Agreed, this show is awesome.","spoiler":false,"review":false,"replies":1,"likes":2}`
		comments = `[{"id":417,"parent_id":190,"comment":"Agreed, this show is awesome.","spoiler":false}]`
		items    = `[{"type":"movie","movie":{"title":"TRON: Legacy","year":2010},"comment":{"id":417,"comment":"Great movie!"}}]`
	)
	cases := []serviceCase{
		{name: "UpdateComment", method: http.MethodPut, path: "/comments/417", body: comment,
			call: func(c *Client) (any, error) {
				r, _, err := c.Comments.UpdateComment(ctx, id, new(str.Comment))
				return r, err
			}},
		{name: "GetComment", method: http.MethodGet, path: "/comments/417", body: comment,
			call: func(c *Client) (any, error) { r, _, err := c.Comments.GetComment(ctx, id); return r, err }},
		{name: "GetCommentItem", method: http.MethodGet, path: "/comments/417/item", query: "limit=10", body: `{"type":"movie","movie":{"title":"TRON: Legacy","year":2010}}`,
			call: func(c *Client) (any, error) { r, _, err := c.Comments.GetCommentItem(ctx, id, opts); return r, err }},
		{name: "DeleteComment", method: http.MethodDelete, path: "/comments/417", status: http.StatusNoContent,
			call: func(c *Client) (any, error) { _, err := c.Comments.DeleteComment(ctx, id); return nil, err }},
		{name: "GetRepliesForComment", method: http.MethodGet, path: "/comments/417/replies", query: "limit=10", body: comments,
			call: func(c *Client) (any, error) {
				r, _, err := c.Comments.GetRepliesForComment(ctx, opts, id)
				return r, err
			}},
		{name: "GetCommentUserLikes", method: http.MethodGet, path: "/comments/417/likes", query: "limit=10", body: `[{"user":{"username":"sean","private":false}}]`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Comments.GetCommentUserLikes(ctx, id, opts)
				return r, err
			}},
		{name: "LikeComment", method: http.MethodPost, path: "/comments/417/like", status: http.StatusNoContent,
			call: func(c *Client) (any, error) { _, err := c.Comments.LikeComment(ctx, id); return nil, err }},
		{name: "RemoveLikeComment", method: http.MethodDelete, path: "/comments/417/like", status: http.StatusNoContent,
			call: func(c *Client) (any, error) { _, err := c.Comments.RemoveLikeComment(ctx, id); return nil, err }},
		{name: "ReplyAComment", method: http.MethodPost, path: "/comments/417/replies", status: http.StatusCreated, body: comment,
			call: func(c *Client) (any, error) {
				r, _, err := c.Comments.ReplyAComment(ctx, id, new(str.Comment))
				return r, err
			}},
		{name: "GetTrendingComments", method: http.MethodGet, path: "/comments/trending/reviews/movies", query: "limit=10", body: items,
			call: func(c *Client) (any, error) {
				r, _, err := c.Comments.GetTrendingComments(ctx, "reviews", "movies", opts)
				return r, err
			}},
		{name: "GetRecentComments", method: http.MethodGet, path: "/comments/recent/shouts/shows", query: "limit=10", body: items,
			call: func(c *Client) (any, error) {
				r, _, err := c.Comments.GetRecentComments(ctx, "shouts", "shows", opts)
				return r, err
			}},
		{name: "GetUpdatedComments", method: http.MethodGet, path: "/comments/updates/all/all", query: "limit=10", body: items,
			call: func(c *Client) (any, error) {
				r, _, err := c.Comments.GetUpdatedComments(ctx, "all", "all", opts)
				return r, err
			}},
	}
	runServiceCases(t, cases)
}

func TestCommentsServiceNotFound(t *testing.T) {
	ctx := context.Background()
	const id = 417
	const (
		commentNotFound = "comment not found with commentId:417"
		itemNotFound    = "comment item not found with commentId:417"
	)
	cases := []struct {
		name string
		want string
		call func(c *Client) error
	}{
		{name: "GetComment", want: commentNotFound, call: func(c *Client) error { _, _, err := c.Comments.GetComment(ctx, id); return err }},
		{name: "GetCommentItem", want: itemNotFound, call: func(c *Client) error { _, _, err := c.Comments.GetCommentItem(ctx, id, nil); return err }},
		{name: "DeleteComment", want: commentNotFound, call: func(c *Client) error { _, err := c.Comments.DeleteComment(ctx, id); return err }},
		{name: "ReplyAComment", want: commentNotFound, call: func(c *Client) error {
			_, _, err := c.Comments.ReplyAComment(ctx, id, new(str.Comment))
			return err
		}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})

			err := tc.call(setup.Client)
			if err == nil {
				t.Fatal("error is nil, want a not found error")
			}
			if err.Error() != tc.want {
				t.Errorf("error is %q, want %q", err.Error(), tc.want)
			}
		})
	}
}
