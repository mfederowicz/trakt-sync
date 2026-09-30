// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

// TestOptionalPathSegments checks that an empty value leaves an optional path
// segment out and a set value adds it.
func TestOptionalPathSegments(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{}
	tests := []struct {
		name string
		path string
		call func(c *Client) error
	}{
		{name: "collection all", path: "/sync/collection", call: func(c *Client) error {
			_, _, err := c.Sync.GetCollection(ctx, "", opts)
			return err
		}},
		{name: "collection type", path: "/sync/collection/movies", call: func(c *Client) error {
			_, _, err := c.Sync.GetCollection(ctx, "movies", opts)
			return err
		}},
		{name: "history all", path: "/sync/history", call: func(c *Client) error {
			_, _, err := c.Sync.GetWatchedHistory(ctx, 0, "", opts)
			return err
		}},
		{name: "history type and id", path: "/sync/history/movies/12", call: func(c *Client) error {
			_, _, err := c.Sync.GetWatchedHistory(ctx, 12, "movies", opts)
			return err
		}},
		{name: "watchlist all", path: "/sync/watchlist", call: func(c *Client) error {
			_, _, err := c.Sync.GetWatchlist(ctx, "movies", "", "", opts)
			return err
		}},
		{name: "watchlist sorted", path: "/sync/watchlist/movies/rank/asc", call: func(c *Client) error {
			_, _, err := c.Sync.GetWatchlist(ctx, "movies", "rank", "asc", opts)
			return err
		}},
		{name: "favorites all", path: "/sync/favorites", call: func(c *Client) error {
			_, _, err := c.Sync.GetFavorites(ctx, "", "rank", "asc", opts)
			return err
		}},
		{name: "favorites sorted", path: "/sync/favorites/shows/rank/asc", call: func(c *Client) error {
			_, _, err := c.Sync.GetFavorites(ctx, "shows", "rank", "asc", opts)
			return err
		}},
		{name: "playback all", path: "/sync/playback", call: func(c *Client) error {
			_, _, err := c.Sync.GetPlaybackProgress(ctx, "", opts)
			return err
		}},
		{name: "playback type", path: "/sync/playback/episodes", call: func(c *Client) error {
			_, _, err := c.Sync.GetPlaybackProgress(ctx, "episodes", opts)
			return err
		}},
		{name: "ratings all", path: "/sync/ratings/movies", call: func(c *Client) error {
			_, _, err := c.Sync.GetRatings(ctx, "movies", "", opts)
			return err
		}},
		{name: "ratings value", path: "/sync/ratings/movies/10", call: func(c *Client) error {
			_, _, err := c.Sync.GetRatings(ctx, "movies", "10", opts)
			return err
		}},
		{name: "list items all", path: "/lists/55/items", call: func(c *Client) error {
			_, _, err := c.Lists.GetListItems(ctx, "55", "", opts)
			return err
		}},
		{name: "list items type", path: "/lists/55/items/movie", call: func(c *Client) error {
			_, _, err := c.Lists.GetListItems(ctx, "55", "movie", opts)
			return err
		}},
		{name: "list comments default", path: "/lists/55/comments", call: func(c *Client) error {
			_, _, err := c.Lists.GetListComments(ctx, "55", "", opts)
			return err
		}},
		{name: "list comments sorted", path: "/lists/55/comments/newest", call: func(c *Client) error {
			_, _, err := c.Lists.GetListComments(ctx, "55", "newest", opts)
			return err
		}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			calls := 0
			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				if r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				test.SafeFprint(w, `[]`)
			})

			test.AssertNilError(t, tt.call(setup.Client))
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
		})
	}
}
