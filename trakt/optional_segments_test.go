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
// segment out (or, for users, becomes "me") and a set value adds it.
func TestOptionalPathSegments(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{}
	tests := []struct {
		name string
		path string
		body string
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
		{name: "movie releases all", path: "/movies/tron/releases", call: func(c *Client) error {
			_, _, err := c.Movies.GetAllMovieReleases(ctx, "tron", "")
			return err
		}},
		{name: "movie releases country", path: "/movies/tron/releases/us", call: func(c *Client) error {
			_, _, err := c.Movies.GetAllMovieReleases(ctx, "tron", "us")
			return err
		}},
		{name: "movie translations all", path: "/movies/tron/translations", call: func(c *Client) error {
			_, _, err := c.Movies.GetAllMovieTranslations(ctx, "tron", "")
			return err
		}},
		{name: "movie translations language", path: "/movies/tron/translations/pl", call: func(c *Client) error {
			_, _, err := c.Movies.GetAllMovieTranslations(ctx, "tron", "pl")
			return err
		}},
		{name: "movie comments default", path: "/movies/tron/comments", call: func(c *Client) error {
			_, _, err := c.Movies.GetAllMovieComments(ctx, "tron", "", opts)
			return err
		}},
		{name: "movie comments sorted", path: "/movies/tron/comments/likes", call: func(c *Client) error {
			_, _, err := c.Movies.GetAllMovieComments(ctx, "tron", "likes", opts)
			return err
		}},
		{name: "movie lists default", path: "/movies/tron/lists", call: func(c *Client) error {
			_, _, err := c.Movies.GetListsContainingMovie(ctx, "tron", "personal", "", opts)
			return err
		}},
		{name: "movie lists typed", path: "/movies/tron/lists/personal/popular", call: func(c *Client) error {
			_, _, err := c.Movies.GetListsContainingMovie(ctx, "tron", "personal", "popular", opts)
			return err
		}},
		{name: "show translations all", path: "/shows/bb/translations", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllShowTranslations(ctx, "bb", "")
			return err
		}},
		{name: "show translations language", path: "/shows/bb/translations/pl", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllShowTranslations(ctx, "bb", "pl")
			return err
		}},
		{name: "season translations all", path: "/shows/bb/seasons/0/translations", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllSeasonTranslations(ctx, "bb", 0, "", opts)
			return err
		}},
		{name: "season translations language", path: "/shows/bb/seasons/1/translations/pl", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllSeasonTranslations(ctx, "bb", 1, "pl", opts)
			return err
		}},
		{name: "episode translations all", path: "/shows/bb/seasons/1/episodes/2/translations", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllEpisodeTranslations(ctx, "bb", 1, 2, "")
			return err
		}},
		{name: "episode translations language", path: "/shows/bb/seasons/1/episodes/2/translations/pl", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllEpisodeTranslations(ctx, "bb", 1, 2, "pl")
			return err
		}},
		{name: "show comments default", path: "/shows/bb/comments", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllShowComments(ctx, "bb", "", opts)
			return err
		}},
		{name: "show comments sorted", path: "/shows/bb/comments/likes", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllShowComments(ctx, "bb", "likes", opts)
			return err
		}},
		{name: "season comments default", path: "/shows/bb/seasons/1/comments", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllSeasonComments(ctx, "bb", 1, "", opts)
			return err
		}},
		{name: "season comments sorted", path: "/shows/bb/seasons/1/comments/likes", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllSeasonComments(ctx, "bb", 1, "likes", opts)
			return err
		}},
		{name: "episode comments default", path: "/shows/bb/seasons/1/episodes/2/comments", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllEpisodeComments(ctx, "bb", 1, 2, "", opts)
			return err
		}},
		{name: "episode comments sorted", path: "/shows/bb/seasons/1/episodes/2/comments/likes", call: func(c *Client) error {
			_, _, err := c.Shows.GetAllEpisodeComments(ctx, "bb", 1, 2, "likes", opts)
			return err
		}},
		{name: "show lists default", path: "/shows/bb/lists", call: func(c *Client) error {
			_, _, err := c.Shows.GetListsContainingShow(ctx, "bb", "", "popular", opts)
			return err
		}},
		{name: "show lists typed", path: "/shows/bb/lists/personal/popular", call: func(c *Client) error {
			_, _, err := c.Shows.GetListsContainingShow(ctx, "bb", "personal", "popular", opts)
			return err
		}},
		{name: "season lists default", path: "/shows/bb/seasons/1/lists", call: func(c *Client) error {
			_, _, err := c.Shows.GetListsContainingSeason(ctx, "bb", 1, "personal", "", opts)
			return err
		}},
		{name: "season lists typed", path: "/shows/bb/seasons/1/lists/personal/popular", call: func(c *Client) error {
			_, _, err := c.Shows.GetListsContainingSeason(ctx, "bb", 1, "personal", "popular", opts)
			return err
		}},
		{name: "episode lists default", path: "/shows/bb/seasons/1/episodes/2/lists", call: func(c *Client) error {
			_, _, err := c.Shows.GetListsContainingEpisode(ctx, "bb", 1, 2, "", "", opts)
			return err
		}},
		{name: "episode lists typed", path: "/shows/bb/seasons/1/episodes/2/lists/personal/popular", call: func(c *Client) error {
			_, _, err := c.Shows.GetListsContainingEpisode(ctx, "bb", 1, 2, "personal", "popular", opts)
			return err
		}},
		{name: "user profile me", path: "/users/me", body: `{}`, call: func(c *Client) error {
			_, _, err := c.Users.GetUserProfile(ctx, "")
			return err
		}},
		{name: "user profile given", path: "/users/sean", body: `{}`, call: func(c *Client) error {
			_, _, err := c.Users.GetUserProfile(ctx, "sean")
			return err
		}},
		{name: "user stats me", path: "/users/me/stats", body: `{}`, call: func(c *Client) error {
			_, _, err := c.Users.GetStats(ctx, "")
			return err
		}},
		{name: "user stats given", path: "/users/sean/stats", body: `{}`, call: func(c *Client) error {
			_, _, err := c.Users.GetStats(ctx, "sean")
			return err
		}},
		{name: "user lists me", path: "/users/me/lists", call: func(c *Client) error {
			_, _, err := c.Users.GetUsersPersonalLists(ctx, "")
			return err
		}},
		{name: "user list items me", path: "/users/me/lists/55/items/movie", call: func(c *Client) error {
			_, _, err := c.Users.GetListItemsByType(ctx, "", "55", "movie")
			return err
		}},
		{name: "user list items given", path: "/users/sean/lists/55/items/show", call: func(c *Client) error {
			_, _, err := c.Users.GetListItemsByType(ctx, "sean", "55", "show")
			return err
		}},
		{name: "user watched me", path: "/users/me/watched/movies", call: func(c *Client) error {
			_, _, err := c.Users.GetWatched(ctx, "", "movies", opts)
			return err
		}},
		{name: "user collaborations me", path: "/users/me/lists/collaborations", call: func(c *Client) error {
			_, _, err := c.Users.GetCollaborations(ctx, "", opts)
			return err
		}},
		{name: "user collaborations given", path: "/users/sean/lists/collaborations", call: func(c *Client) error {
			_, _, err := c.Users.GetCollaborations(ctx, "sean", opts)
			return err
		}},
		{name: "user likes me", path: "/users/me/likes/lists", call: func(c *Client) error {
			_, _, err := c.Users.GetLikes(ctx, "", "lists", opts)
			return err
		}},
		{name: "user likes given", path: "/users/sean/likes/comments", call: func(c *Client) error {
			_, _, err := c.Users.GetLikes(ctx, "sean", "comments", opts)
			return err
		}},
		{name: "user collection given", path: "/users/sean/collection/shows", call: func(c *Client) error {
			_, _, err := c.Users.GetCollection(ctx, "sean", "shows", opts)
			return err
		}},
		{name: "user comments given", path: "/users/sean/comments/reviews/movies", call: func(c *Client) error {
			_, _, err := c.Users.GetComments(ctx, "sean", "reviews", "movies", opts)
			return err
		}},
		{name: "user notes given", path: "/users/sean/notes/movies", call: func(c *Client) error {
			_, _, err := c.Users.GetNotes(ctx, "sean", "movies", opts)
			return err
		}},
		{name: "hidden section", path: "/users/hidden/calendar", call: func(c *Client) error {
			_, _, err := c.Users.GetHiddenItems(ctx, "calendar", opts)
			return err
		}},
		{name: "user history all", path: "/users/sean/history/movies", call: func(c *Client) error {
			_, _, err := c.Users.GetHistory(ctx, "sean", "movies", 0, opts)
			return err
		}},
		{name: "user history id", path: "/users/sean/history/movies/12", call: func(c *Client) error {
			_, _, err := c.Users.GetHistory(ctx, "sean", "movies", 12, opts)
			return err
		}},
		{name: "user ratings all", path: "/users/sean/ratings/movies", call: func(c *Client) error {
			_, _, err := c.Users.GetRatings(ctx, "sean", "movies", "", opts)
			return err
		}},
		{name: "user ratings value", path: "/users/sean/ratings/movies/10", call: func(c *Client) error {
			_, _, err := c.Users.GetRatings(ctx, "sean", "movies", "10", opts)
			return err
		}},
		{name: "data syncs all", path: "/users/syncs", call: func(c *Client) error {
			_, _, err := c.Users.GetDataSyncs(ctx, "", opts)
			return err
		}},
		{name: "data syncs type", path: "/users/syncs/plex", call: func(c *Client) error {
			_, _, err := c.Users.GetDataSyncs(ctx, "plex", opts)
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
				body := tt.body
				if body == "" {
					body = `[]`
				}
				test.SafeFprint(w, body)
			})

			test.AssertNilError(t, tt.call(setup.Client))
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
		})
	}
}
