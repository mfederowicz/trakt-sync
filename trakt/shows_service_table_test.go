// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

// serviceCase is one request/response check: the method must send method+path(+query)
// and decode body into its result type without losing or renaming a field.
type serviceCase struct {
	name   string
	method string
	path   string
	query  string
	status int
	body   string
	call   func(c *Client) (any, error)
}

// runServiceCases runs each case against the mock server. A non-empty body is decoded by the method;
// the result is encoded again and compared with body, which fails when a field is dropped or renamed.
func runServiceCases(t *testing.T, cases []serviceCase) {
	t.Helper()
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			calls := 0
			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, tc.method)
				if r.URL.Path != tc.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tc.path)
				}
				if r.URL.RawQuery != tc.query {
					t.Errorf("query is %q, want %q", r.URL.RawQuery, tc.query)
				}
				if tc.status != 0 {
					w.WriteHeader(tc.status)
				}
				if tc.body != "" {
					test.SafeFprint(w, tc.body)
				}
			})

			got, err := tc.call(setup.Client)
			test.AssertNilError(t, err)
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
			if tc.body == "" {
				return
			}
			assertSameJSON(t, tc.body, got)
		})
	}
}

// assertSameJSON checks got encodes to the same JSON value as want.
func assertSameJSON(t *testing.T, want string, got any) {
	t.Helper()
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("encode result: %v", err)
	}
	var wantValue, gotValue any
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("bad test body: %v", err)
	}
	if err := json.Unmarshal(encoded, &gotValue); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if !reflect.DeepEqual(wantValue, gotValue) {
		t.Errorf("decoded result is %s, want %s", encoded, want)
	}
}

func TestShowsServiceRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	const id = "breaking-bad"
	const (
		showsItems = `[{"watchers":5,"show":{"title":"Breaking Bad","year":2008}}]`
		shows      = `[{"title":"Breaking Bad","year":2008}]`
		episode    = `{"season":1,"number":2,"title":"Cat's in the Bag..."}`
		episodes   = `[{"season":1,"number":1,"title":"Pilot"}]`
		users      = `[{"username":"sean","private":false}]`
		videos     = `[{"title":"Trailer","site":"youtube","type":"trailer"}]`
		ratings    = `{"rating":9.3,"votes":120}`
		stats      = `{"watchers":10,"plays":20,"collectors":3}`
		people     = `{"cast":[{"character":"Walter White"}]}`
	)
	cases := []serviceCase{
		{name: "GetSingleEpisodeForShow", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/episodes/2", query: "limit=10", body: episode,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetSingleEpisodeForShow(ctx, id, 1, 2, opts)
				return r, err
			}},
		{name: "GetTrendingShows", method: http.MethodGet, path: "/shows/trending", query: "limit=10", body: showsItems,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetTrendingShows(ctx, opts); return r, err }},
		{name: "GetPopularShows", method: http.MethodGet, path: "/shows/popular", query: "limit=10", body: shows,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetPopularShows(ctx, opts); return r, err }},
		{name: "GetFavoritedShows", method: http.MethodGet, path: "/shows/favorited/weekly", query: "limit=10", body: showsItems,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetFavoritedShows(ctx, opts, "weekly")
				return r, err
			}},
		{name: "GetPlayedShows", method: http.MethodGet, path: "/shows/played/monthly", query: "limit=10", body: showsItems,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetPlayedShows(ctx, opts, "monthly"); return r, err }},
		{name: "GetWatchedShows", method: http.MethodGet, path: "/shows/watched/yearly", query: "limit=10", body: showsItems,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetWatchedShows(ctx, opts, "yearly"); return r, err }},
		{name: "GetCollectedShows", method: http.MethodGet, path: "/shows/collected/all", query: "limit=10", body: showsItems,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetCollectedShows(ctx, opts, "all"); return r, err }},
		{name: "GetAnticipatedShows", method: http.MethodGet, path: "/shows/anticipated", query: "limit=10", body: showsItems,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetAnticipatedShows(ctx, opts); return r, err }},
		{name: "GetRecentlyUpdatedShows", method: http.MethodGet, path: "/shows/updates/2026-10-01", query: "limit=10", body: showsItems,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetRecentlyUpdatedShows(ctx, "2026-10-01", opts)
				return r, err
			}},
		{name: "GetRecentlyUpdatedShowsTraktIDs", method: http.MethodGet, path: "/shows/updates/id/2026-10-01", query: "limit=10", body: `[1388,1390]`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetRecentlyUpdatedShowsTraktIDs(ctx, "2026-10-01", opts)
				return r, err
			}},
		{name: "GetAllShowAliases", method: http.MethodGet, path: "/shows/breaking-bad/aliases", body: `[{"title":"Breaking Bad","country":"us"}]`,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetAllShowAliases(ctx, id); return r, err }},
		{name: "GetAllShowCertifications", method: http.MethodGet, path: "/shows/breaking-bad/certifications", body: `[{"country":"us","certification":"TV-MA"}]`,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetAllShowCertifications(ctx, id); return r, err }},
		{name: "ResetShowProgress", method: http.MethodPost, path: "/shows/breaking-bad/progress/watched/reset", body: `{"aired":62,"completed":10}`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.ResetShowProgress(ctx, id, new(str.WatchedProgress))
				return r, err
			}},
		{name: "UndoResetShowProgress", method: http.MethodDelete, path: "/shows/breaking-bad/progress/watched/reset", status: http.StatusNoContent,
			call: func(c *Client) (any, error) { _, err := c.Shows.UndoResetShowProgress(ctx, id); return nil, err }},
		{name: "GetShowRatings", method: http.MethodGet, path: "/shows/breaking-bad/ratings", body: ratings,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetShowRatings(ctx, id); return r, err }},
		{name: "GetRelatedShows", method: http.MethodGet, path: "/shows/breaking-bad/related", query: "limit=10", body: shows,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetRelatedShows(ctx, id, opts); return r, err }},
		{name: "GetShowStats", method: http.MethodGet, path: "/shows/breaking-bad/stats", body: stats,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetShowStats(ctx, id); return r, err }},
		{name: "GetShowStudios", method: http.MethodGet, path: "/shows/breaking-bad/studios", body: `[{"name":"Sony Pictures Television","country":"us"}]`,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetShowStudios(ctx, id); return r, err }},
		{name: "GetShowWatching", method: http.MethodGet, path: "/shows/breaking-bad/watching", query: "limit=10", body: users,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetShowWatching(ctx, id, opts); return r, err }},
		{name: "GetShowVideos", method: http.MethodGet, path: "/shows/breaking-bad/videos", query: "limit=10", body: videos,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetShowVideos(ctx, id, opts); return r, err }},
		{name: "RefreshShowMetadata", method: http.MethodPost, path: "/shows/breaking-bad/refresh", status: http.StatusCreated,
			call: func(c *Client) (any, error) { _, err := c.Shows.RefreshShowMetadata(ctx, id); return nil, err }},
		{name: "GetNextEpisode", method: http.MethodGet, path: "/shows/breaking-bad/next_episode", query: "limit=10", body: episode,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetNextEpisode(ctx, id, opts); return r, err }},
		{name: "GetLastEpisode", method: http.MethodGet, path: "/shows/breaking-bad/last_episode", query: "limit=10", body: episode,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetLastEpisode(ctx, id, opts); return r, err }},
		{name: "GetAllSeasonsForShow", method: http.MethodGet, path: "/shows/breaking-bad/seasons", query: "limit=10", body: `[{"number":1,"title":"Season 1"}]`,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetAllSeasonsForShow(ctx, id, opts); return r, err }},
		{name: "GetSingleSeasonsForShow", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/info", query: "limit=10", body: `{"number":1,"title":"Season 1"}`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetSingleSeasonsForShow(ctx, id, 1, opts)
				return r, err
			}},
		{name: "GetAllEpisodesForSingleSeason", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1", query: "limit=10", body: episodes,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetAllEpisodesForSingleSeason(ctx, id, 1, opts)
				return r, err
			}},
		{name: "GetAllPeopleForSeason", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/people", query: "limit=10", body: people,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetAllPeopleForSeason(ctx, id, 1, opts)
				return r, err
			}},
		{name: "GetAllPeopleForEpisode", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/episodes/2/people", query: "limit=10", body: people,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetAllPeopleForEpisode(ctx, id, 1, 2, opts)
				return r, err
			}},
		{name: "GetSeasonRatings", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/ratings", body: ratings,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetSeasonRatings(ctx, id, 1); return r, err }},
		{name: "GetEpisodeRatings", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/episodes/2/ratings", body: ratings,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetEpisodeRatings(ctx, id, 1, 2); return r, err }},
		{name: "GetSeasonStats", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/stats", body: stats,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetSeasonStats(ctx, id, 1); return r, err }},
		{name: "GetEpisodeStats", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/episodes/2/stats", body: stats,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetEpisodeStats(ctx, id, 1, 2); return r, err }},
		{name: "GetSeasonsWatching", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/watching", query: "limit=10", body: users,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetSeasonsWatching(ctx, id, 1, opts); return r, err }},
		{name: "GetEpisodesWatching", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/episodes/2/watching", query: "limit=10", body: users,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetEpisodesWatching(ctx, id, 1, 2, opts)
				return r, err
			}},
		{name: "GetSeasonsVideos", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/videos", query: "limit=10", body: videos,
			call: func(c *Client) (any, error) { r, _, err := c.Shows.GetSeasonsVideos(ctx, id, 1, opts); return r, err }},
		{name: "GetEpisodeVideos", method: http.MethodGet, path: "/shows/breaking-bad/seasons/1/episodes/2/videos", query: "limit=10", body: videos,
			call: func(c *Client) (any, error) {
				r, _, err := c.Shows.GetEpisodeVideos(ctx, id, 1, 2, opts)
				return r, err
			}},
	}
	runServiceCases(t, cases)
}
