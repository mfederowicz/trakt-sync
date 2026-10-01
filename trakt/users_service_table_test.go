// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestUsersServiceRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	const user = "sean"
	const listID = "star-wars"
	const (
		followRequest = `{"id":3,"user":{"username":"sean","private":false}}`
		list          = `{"name":"Star Wars","privacy":"public","item_count":5}`
		items         = `[{"title":"TRON: Legacy","year":2010,"movie":{"title":"TRON: Legacy","year":2010}}]`
		reordered     = `{"updated":2,"skipped_ids":[12]}`
	)
	cases := []serviceCase{
		{name: "GetSavedFilters", method: http.MethodGet, path: "/users/saved_filters/movies", body: `[{"rank":1,"id":101,"section":"movies","name":"Movies: IMDB + TMDB ratings","path":"/movies/recommended/weekly","query":"imdb_ratings=6.9-10.0"}]`,
			call: func(c *Client) (any, error) { r, _, err := c.Users.GetSavedFilters(ctx, "movies"); return r, err }},
		{name: "GetSettings", method: http.MethodGet, path: "/users/settings", body: `{"user":{"username":"sean","private":false},"account":{"timezone":"America/Los_Angeles","time_24hr":false}}`,
			call: func(c *Client) (any, error) { r, _, err := c.Users.GetSettings(ctx); return r, err }},
		{name: "ApproveFollowRequest", method: http.MethodPost, path: "/users/requests/3", body: followRequest,
			call: func(c *Client) (any, error) { r, _, err := c.Users.ApproveFollowRequest(ctx, 3); return r, err }},
		{name: "DenyFollowRequest", method: http.MethodDelete, path: "/users/requests/3", body: followRequest,
			call: func(c *Client) (any, error) { r, _, err := c.Users.DenyFollowRequest(ctx, 3); return r, err }},
		{name: "GetProfile", method: http.MethodGet, path: "/users/sean", body: `{"username":"sean","private":false,"name":"Sean Rudford","vip":true}`,
			call: func(c *Client) (any, error) { r, _, err := c.Users.GetProfile(ctx, user); return r, err }},
		{name: "ReorderLists", method: http.MethodPost, path: "/users/sean/lists/reorder", body: reordered,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.ReorderLists(ctx, user, new(str.ItemsToReorder))
				return r, err
			}},
		{name: "ReorderLists without user", method: http.MethodPost, path: "/users/me/lists/reorder", body: reordered,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.ReorderLists(ctx, "", new(str.ItemsToReorder))
				return r, err
			}},
		{name: "GetList", method: http.MethodGet, path: "/users/sean/lists/star-wars", query: "limit=10", body: list,
			call: func(c *Client) (any, error) { r, _, err := c.Users.GetList(ctx, user, listID, opts); return r, err }},
		{name: "RemoveListItems", method: http.MethodPost, path: "/users/sean/lists/star-wars/items/remove", body: `{"deleted":{"movies":1},"list":{"item_count":4}}`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.RemoveListItems(ctx, user, listID, new(str.HistoryItems))
				return r, err
			}},
		{name: "ReorderListItems", method: http.MethodPost, path: "/users/sean/lists/star-wars/items/reorder", body: reordered,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.ReorderListItems(ctx, user, listID, new(str.ItemsToReorder))
				return r, err
			}},
		{name: "UpdateListItem", method: http.MethodPut, path: "/users/sean/lists/star-wars/items/7", status: http.StatusNoContent,
			call: func(c *Client) (any, error) {
				_, err := c.Users.UpdateListItem(ctx, user, listID, 7, new(str.PersonalListItem))
				return nil, err
			}},
		{name: "Unblock", method: http.MethodDelete, path: "/users/sean/block", status: http.StatusNoContent,
			call: func(c *Client) (any, error) { _, err := c.Users.Unblock(ctx, user); return nil, err }},
		{name: "GetWatchlistComments", method: http.MethodGet, path: "/users/sean/watchlist/comments/newest", query: "limit=10", body: items,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.GetWatchlistComments(ctx, user, "newest", opts)
				return r, err
			}},
		{name: "GetFavorites", method: http.MethodGet, path: "/users/sean/favorites/movies/rank/asc", query: "limit=10", body: items,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.GetFavorites(ctx, user, "movies", "rank", "asc", opts)
				return r, err
			}},
		{name: "GetFavoritesComments", method: http.MethodGet, path: "/users/sean/favorites/comments/likes", query: "limit=10", body: items,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.GetFavoritesComments(ctx, user, "likes", opts)
				return r, err
			}},
		{name: "GetWatchlistBySort", method: http.MethodGet, path: "/users/sean/watchlist/movies/added", query: "limit=10", body: items,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.GetWatchlistBySort(ctx, user, "movies", "added", opts)
				return r, err
			}},
		{name: "GetFavoritesBySort", method: http.MethodGet, path: "/users/me/favorites/shows/rank", query: "limit=10", body: items,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.GetFavoritesBySort(ctx, "", "shows", "rank", opts)
				return r, err
			}},
	}
	runServiceCases(t, cases)
}
