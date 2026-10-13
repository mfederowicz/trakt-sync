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
		comments      = `[{"id":8,"comment":"Can't wait to watch everything on this epic list!","spoiler":false,"likes":2}]`
		reordered     = `{"updated":2,"skipped_ids":[12]}`
	)
	cases := []serviceCase{
		{name: "GetSavedFilters", method: http.MethodGet, path: "/users/saved_filters/movies", body: `[{"rank":1,"id":101,"section":"movies","name":"Movies: IMDB + TMDB ratings","path":"/movies/recommended/weekly","query":"imdb_ratings=6.9-10.0"}]`,
			call: func(c *Client) (any, error) { r, _, err := c.Users.GetSavedFilters(ctx, "movies"); return r, err }},
		{name: "GetYearInReviewPersona", method: http.MethodGet, path: "/users/sean/yir/2026/persona",
			body: `{"persona":"cinephile","runner_up":"critic","confidence":"strong","rarity":12,"card_number":7,"traits":["night-owl","doc-nerd"],` +
				`"highlights":[{"kind":"movies","value":212},{"kind":"avg-rating","value":7.4}],"runner_up_highlights":[{"kind":"ratings","value":180}],` +
				`"scores":{"cinephile":91,"critic":77},"streak":{"longest":21,"started_at":"2026-03-01T10:00:00Z"},"monthly":[{"month":1,"persona":"cinephile"}]}`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.GetYearInReviewPersona(ctx, user, 2026)
				return r, err
			}},
		{name: "GetSmartListDrafts", method: http.MethodGet, path: "/users/sean/smart-list-drafts",
			body: `{"convertible_count":2,"best_confidence":0.92,"lists":[{"ids":{"trakt":55,"slug":"star-wars"},"name":"Star Wars","item_count":12,"confidence":0.92}]}`,
			call: func(c *Client) (any, error) { r, _, err := c.Users.GetSmartListDrafts(ctx, user); return r, err }},
		{name: "GetSmartListDraft", method: http.MethodGet, path: "/users/sean/lists/star-wars/smart-list-draft",
			body: `{"list":{"ids":{"trakt":55,"slug":"star-wars"}},"origin_url":"https://example.com/lists/star-wars","source":"popular","media_type":"movies",` +
				`"filters":{"genres":["science-fiction"],"years":[1977,2019]},"unsupported":["collection"],"confidence":0.92,` +
				`"converted":{"smart_list_id":7,"converted_at":"2026-10-12T08:30:00Z"}}`,
			call: func(c *Client) (any, error) { r, _, err := c.Users.GetSmartListDraft(ctx, user, listID); return r, err }},
		{name: "ConvertListToSmartList", method: http.MethodPost, path: "/users/sean/lists/star-wars/smart-list-draft/convert", status: http.StatusCreated,
			body: `{"ids":{"trakt":7,"slug":"star-wars-smart"}}`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.ConvertListToSmartList(ctx, user, listID, &str.SmartListWrite{Name: str.String("Star Wars"), Source: str.String("popular"), MediaType: str.String("movies")})
				return r, err
			}},
		{name: "GetSettings", method: http.MethodGet, path: "/users/settings", body: `{"user":{"username":"sean","private":false},"account":{"timezone":"America/Los_Angeles","time_24hr":false}}`,
			call: func(c *Client) (any, error) { r, _, err := c.Users.GetSettings(ctx); return r, err }},
		{name: "ApproveFollowRequest", method: http.MethodPost, path: "/users/requests/3", body: followRequest,
			call: func(c *Client) (any, error) { r, _, err := c.Users.ApproveFollowRequest(ctx, 3); return r, err }},
		{name: "DenyFollowRequest", method: http.MethodDelete, path: "/users/requests/3", body: followRequest,
			call: func(c *Client) (any, error) { r, _, err := c.Users.DenyFollowRequest(ctx, 3); return r, err }},
		{name: "GetHiddenItems users", method: http.MethodGet, path: "/users/hidden/comments", query: "limit=10&type=user", body: `[{"type":"user","user":{"username":"sean","private":false,"ids":{"slug":"sean"}}}]`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.GetHiddenItems(ctx, "comments", &uri.ListOptions{Limit: 10, Type: "user"})
				return r, err
			}},
		{name: "AddHiddenItems users", method: http.MethodPost, path: "/users/hidden/comments", status: http.StatusCreated, body: `{"added":{"seasons":0,"users":1},"not_found":{"episodes":[],"people":[],"users":[{"ids":{"slug":"nobody"}}]}}`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.AddHiddenItems(ctx, new(str.HistoryItems), "comments")
				return r, err
			}},
		{name: "RemoveHiddenItems users", method: http.MethodPost, path: "/users/hidden/comments/remove", body: `{"deleted":{"users":1},"not_found":{"people":[],"users":[]}}`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.RemoveHiddenItems(ctx, new(str.HistoryItems), "comments")
				return r, err
			}},
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
		{name: "GetWatchlistComments", method: http.MethodGet, path: "/users/sean/watchlist/comments/newest", query: "limit=10", body: comments,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.GetWatchlistComments(ctx, user, "newest", opts)
				return r, err
			}},
		{name: "GetFavorites", method: http.MethodGet, path: "/users/sean/favorites/movies/rank/asc", query: "limit=10", body: items,
			call: func(c *Client) (any, error) {
				r, _, err := c.Users.GetFavorites(ctx, user, "movies", "rank", "asc", opts)
				return r, err
			}},
		{name: "GetFavoritesComments", method: http.MethodGet, path: "/users/sean/favorites/comments/likes", query: "limit=10", body: comments,
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
