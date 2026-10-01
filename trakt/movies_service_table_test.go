// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/uri"
)

func TestMoviesServiceRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	const id = "tron-legacy-2010"
	const (
		moviesItems = `[{"watcher_count":5,"movie":{"title":"TRON: Legacy","year":2010}}]`
		movies      = `[{"title":"TRON: Legacy","year":2010}]`
		boxoffice   = `[{"revenue":48000000,"movie":{"title":"TRON: Legacy","year":2010}}]`
		ids         = `[1,20]`
		aliases     = `[{"title":"Tron 2","country":"us"}]`
		people      = `{"cast":[{"character":"Kevin Flynn"}]}`
		ratings     = `{"rating":7.3,"votes":120,"distribution":{"10":20}}`
		stats       = `{"watchers":10,"plays":20,"collectors":3,"lists":4}`
		studios     = `[{"name":"Walt Disney Pictures","country":"us"}]`
		users       = `[{"username":"sean","private":false}]`
		videos      = `[{"title":"Trailer","site":"youtube","type":"trailer"}]`
	)
	cases := []serviceCase{
		{name: "GetTrendingMovies", method: http.MethodGet, path: "/movies/trending", query: "limit=10", body: moviesItems,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetTrendingMovies(ctx, opts); return r, err }},
		{name: "GetPopularMovies", method: http.MethodGet, path: "/movies/popular", query: "limit=10", body: movies,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetPopularMovies(ctx, opts); return r, err }},
		{name: "GetFavoritedMovies", method: http.MethodGet, path: "/movies/favorited/weekly", query: "limit=10", body: moviesItems,
			call: func(c *Client) (any, error) {
				r, _, err := c.Movies.GetFavoritedMovies(ctx, opts, "weekly")
				return r, err
			}},
		{name: "GetPlayedMovies", method: http.MethodGet, path: "/movies/played/monthly", query: "limit=10", body: moviesItems,
			call: func(c *Client) (any, error) {
				r, _, err := c.Movies.GetPlayedMovies(ctx, opts, "monthly")
				return r, err
			}},
		{name: "GetWatchedMovies", method: http.MethodGet, path: "/movies/watched/yearly", query: "limit=10", body: moviesItems,
			call: func(c *Client) (any, error) {
				r, _, err := c.Movies.GetWatchedMovies(ctx, opts, "yearly")
				return r, err
			}},
		{name: "GetCollectedMovies", method: http.MethodGet, path: "/movies/collected/all", query: "limit=10", body: moviesItems,
			call: func(c *Client) (any, error) {
				r, _, err := c.Movies.GetCollectedMovies(ctx, opts, "all")
				return r, err
			}},
		{name: "GetAnticipatedMovies", method: http.MethodGet, path: "/movies/anticipated", query: "limit=10", body: moviesItems,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetAnticipatedMovies(ctx, opts); return r, err }},
		{name: "GetBoxoffice", method: http.MethodGet, path: "/movies/boxoffice", query: "limit=10", body: boxoffice,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetBoxoffice(ctx, opts); return r, err }},
		{name: "GetRecentlyUpdatedMovies", method: http.MethodGet, path: "/movies/updates/2024-01-01", query: "limit=10", body: moviesItems,
			call: func(c *Client) (any, error) {
				r, _, err := c.Movies.GetRecentlyUpdatedMovies(ctx, "2024-01-01", opts)
				return r, err
			}},
		{name: "GetRecentlyUpdatedMoviesTraktIDs", method: http.MethodGet, path: "/movies/updates/id/2024-01-01", query: "limit=10", body: ids,
			call: func(c *Client) (any, error) {
				r, _, err := c.Movies.GetRecentlyUpdatedMoviesTraktIDs(ctx, "2024-01-01", opts)
				return r, err
			}},
		{name: "GetAllMovieAliases", method: http.MethodGet, path: "/movies/tron-legacy-2010/aliases", body: aliases,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetAllMovieAliases(ctx, id); return r, err }},
		{name: "GetAllPeopleForMovie", method: http.MethodGet, path: "/movies/tron-legacy-2010/people", query: "limit=10", body: people,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetAllPeopleForMovie(ctx, id, opts); return r, err }},
		{name: "GetMovieRatings", method: http.MethodGet, path: "/movies/tron-legacy-2010/ratings", body: ratings,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetMovieRatings(ctx, id); return r, err }},
		{name: "GetRelatedMovies", method: http.MethodGet, path: "/movies/tron-legacy-2010/related", query: "limit=10", body: movies,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetRelatedMovies(ctx, id, opts); return r, err }},
		{name: "GetMovieStats", method: http.MethodGet, path: "/movies/tron-legacy-2010/stats", body: stats,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetMovieStats(ctx, id); return r, err }},
		{name: "GetMovieStudios", method: http.MethodGet, path: "/movies/tron-legacy-2010/studios", body: studios,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetMovieStudios(ctx, id); return r, err }},
		{name: "GetMovieWatching", method: http.MethodGet, path: "/movies/tron-legacy-2010/watching", query: "limit=10", body: users,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetMovieWatching(ctx, id, opts); return r, err }},
		{name: "GetMovieVideos", method: http.MethodGet, path: "/movies/tron-legacy-2010/videos", query: "limit=10", body: videos,
			call: func(c *Client) (any, error) { r, _, err := c.Movies.GetMovieVideos(ctx, id, opts); return r, err }},
		{name: "RefreshMovieMetadata", method: http.MethodPost, path: "/movies/tron-legacy-2010/refresh", status: http.StatusCreated,
			call: func(c *Client) (any, error) { _, err := c.Movies.RefreshMovieMetadata(ctx, id); return nil, err }},
	}
	runServiceCases(t, cases)
}
