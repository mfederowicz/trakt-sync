// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestListsServiceRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	const id = "2143363"
	const lists = `[{"like_count":5,"comment_count":2,"list":{"name":"Star Wars","privacy":"public","item_count":5}}]`
	cases := []serviceCase{
		{name: "GetTrendingLists", method: http.MethodGet, path: "/lists/trending", query: "limit=10", body: lists,
			call: func(c *Client) (any, error) { r, _, err := c.Lists.GetTrendingLists(ctx, opts); return r, err }},
		{name: "GetPopularLists", method: http.MethodGet, path: "/lists/popular", query: "limit=10", body: lists,
			call: func(c *Client) (any, error) { r, _, err := c.Lists.GetPopularLists(ctx, opts); return r, err }},
		{name: "GetList", method: http.MethodGet, path: "/lists/2143363", body: `{"name":"Star Wars","privacy":"public","item_count":5}`,
			call: func(c *Client) (any, error) { r, _, err := c.Lists.GetList(ctx, id); return r, err }},
		{name: "GetAllUsersWhoLikedList", method: http.MethodGet, path: "/lists/2143363/likes", query: "limit=10", body: `[{"user":{"username":"sean","private":false}}]`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Lists.GetAllUsersWhoLikedList(ctx, opts, id)
				return r, err
			}},
		{name: "LikeList", method: http.MethodPost, path: "/lists/2143363/like", status: http.StatusNoContent,
			call: func(c *Client) (any, error) { _, err := c.Lists.LikeList(ctx, id); return nil, err }},
		{name: "RemoveLikeList", method: http.MethodDelete, path: "/lists/2143363/like", status: http.StatusNoContent,
			call: func(c *Client) (any, error) { _, err := c.Lists.RemoveLikeList(ctx, id); return nil, err }},
	}
	runServiceCases(t, cases)
}

func TestListsServiceGetListNotFound(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/lists/2143363", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusNotFound)
	})

	_, _, err := setup.Client.Lists.GetList(context.Background(), "2143363")
	if err == nil {
		t.Fatal("error is nil, want a not found error")
	}
	const want = "list not found with traktId:2143363"
	if err.Error() != want {
		t.Errorf("error is %q, want %q", err.Error(), want)
	}
}
