// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

// TestServicesReturnResponse checks that the sync write methods, show progress
// and hidden items return the *str.Response, so callers can read its headers.
func TestServicesReturnResponse(t *testing.T) {
	ctx := context.Background()
	show := "bb"
	opts := &uri.ListOptions{}
	tests := []struct {
		name   string
		method string
		path   string
		call   func(c *Client) (*str.Response, error)
	}{
		{name: "AddItemsToCollection", method: http.MethodPost, path: "/sync/collection", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.AddItemsToCollection(ctx, new(str.ItemsList))
			return resp, err
		}},
		{name: "RemoveItemsFromCollection", method: http.MethodPost, path: "/sync/collection/remove", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.RemoveItemsFromCollection(ctx, new(str.ItemsList))
			return resp, err
		}},
		{name: "AddItemsToHistory", method: http.MethodPost, path: "/sync/history", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.AddItemsToHistory(ctx, new(str.HistoryItems))
			return resp, err
		}},
		{name: "RemoveItemsFromHistory", method: http.MethodPost, path: "/sync/history/remove", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.RemoveItemsFromHistory(ctx, new(str.ItemsToRemove))
			return resp, err
		}},
		{name: "AddItemsToRatings", method: http.MethodPost, path: "/sync/ratings", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.AddItemsToRatings(ctx, new(str.RatingItems))
			return resp, err
		}},
		{name: "RemoveItemsFromRatings", method: http.MethodPost, path: "/sync/ratings/remove", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.RemoveItemsFromRatings(ctx, new(str.ItemsToRemove))
			return resp, err
		}},
		{name: "UpdateWatchlist", method: http.MethodPut, path: "/sync/watchlist", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.UpdateWatchlist(ctx, new(str.PersonalList))
			return resp, err
		}},
		{name: "UpdateWatchlistItem", method: http.MethodPut, path: "/sync/watchlist/7", call: func(c *Client) (*str.Response, error) {
			return c.Sync.UpdateWatchlistItem(ctx, 7, new(str.WatchlistItem))
		}},
		{name: "AddItemsToWatchlist", method: http.MethodPost, path: "/sync/watchlist", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.AddItemsToWatchlist(ctx, new(str.HistoryItems))
			return resp, err
		}},
		{name: "RemoveItemsFromWatchlist", method: http.MethodPost, path: "/sync/watchlist/remove", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.RemoveItemsFromWatchlist(ctx, new(str.ItemsToRemove))
			return resp, err
		}},
		{name: "ReorderWatchlistItems", method: http.MethodPost, path: "/sync/watchlist/reorder", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.ReorderWatchlistItems(ctx, new(str.ItemsToReorder))
			return resp, err
		}},
		{name: "UpdateFavorites", method: http.MethodPut, path: "/sync/favorites", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.UpdateFavorites(ctx, new(str.PersonalList))
			return resp, err
		}},
		{name: "UpdateFavoriteItem", method: http.MethodPut, path: "/sync/favorites/7", call: func(c *Client) (*str.Response, error) {
			return c.Sync.UpdateFavoriteItem(ctx, 7, new(str.FavoriteItem))
		}},
		{name: "AddItemsToFavorites", method: http.MethodPost, path: "/sync/favorites", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.AddItemsToFavorites(ctx, new(str.HistoryItems))
			return resp, err
		}},
		{name: "RemoveItemsFromFavorites", method: http.MethodPost, path: "/sync/favorites/remove", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.RemoveItemsFromFavorites(ctx, new(str.ItemsToRemove))
			return resp, err
		}},
		{name: "ReorderFavoritesItems", method: http.MethodPost, path: "/sync/favorites/reorder", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.ReorderFavoritesItems(ctx, new(str.ItemsToReorder))
			return resp, err
		}},
		{name: "GetShowCollectionProgress", method: http.MethodGet, path: "/shows/bb/progress/collection", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Shows.GetShowCollectionProgress(ctx, show, opts)
			return resp, err
		}},
		{name: "GetShowWatchedProgress", method: http.MethodGet, path: "/shows/bb/progress/watched", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Shows.GetShowWatchedProgress(ctx, show, opts)
			return resp, err
		}},
		{name: "AddHiddenItems", method: http.MethodPost, path: "/users/hidden/calendar", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Users.AddHiddenItems(ctx, new(str.HistoryItems), "calendar")
			return resp, err
		}},
		{name: "RemoveHiddenItems", method: http.MethodPost, path: "/users/hidden/calendar/remove", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Users.RemoveHiddenItems(ctx, new(str.HistoryItems), "calendar")
			return resp, err
		}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, tt.method)
				if r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				w.Header().Set("X-Ratelimit", "ok")
				test.SafeFprint(w, `{}`)
			})

			resp, err := tt.call(setup.Client)
			test.AssertNilError(t, err)
			if resp == nil {
				t.Fatal("resp is nil")
			}
			if got := resp.Header.Get("X-Ratelimit"); got != "ok" {
				t.Errorf("X-Ratelimit header is %q, want %q", got, "ok")
			}
		})
	}
}
