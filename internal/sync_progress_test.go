// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestSyncServiceShowProgress(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		opts      uri.SyncProgressOptions
		wantQuery string
		call      func(c *Client, opts *uri.SyncProgressOptions) ([]*str.ShowProgress, *str.Response, error)
	}{
		{
			name:      "up next",
			path:      "/sync/progress/up_next",
			opts:      uri.SyncProgressOptions{Page: 1, IncludeStats: true},
			wantQuery: "include_stats=true&page=1",
			call: func(c *Client, opts *uri.SyncProgressOptions) ([]*str.ShowProgress, *str.Response, error) {
				return c.Sync.GetUpNext(context.Background(), opts)
			},
		},
		{
			name:      "watched progress",
			path:      "/sync/progress/watched",
			opts:      uri.SyncProgressOptions{Page: 1, SortBy: "watched", OnlyRewatching: true},
			wantQuery: "only_rewatching=true&page=1&sort_by=watched",
			call: func(c *Client, opts *uri.SyncProgressOptions) ([]*str.ShowProgress, *str.Response, error) {
				return c.Sync.GetWatchedProgress(context.Background(), opts)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				test.AssertNoDiff(t, tt.wantQuery, r.URL.RawQuery)
				test.SafeFprint(w, `[{"show":{"title":"Reacher"},"progress":{"aired":24,"completed":20,"next_episode":{"season":3,"number":5},"stats":{"play_count":20,"minutes_watched":960,"minutes_left":192}}}]`)
			})

			got, _, err := tt.call(setup.Client, &tt.opts)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, []*str.ShowProgress{{
				Show: &str.Show{Title: str.String("Reacher")},
				Progress: &str.WatchedProgress{
					Aired:       test.Ptr(24),
					Completed:   test.Ptr(20),
					NextEpisode: &str.Episode{Season: test.Ptr(3), Number: test.Ptr(5)},
					Stats:       &str.ProgressStats{PlayCount: test.Ptr(20), MinutesWatched: test.Ptr(960), MinutesLeft: test.Ptr(192)},
				},
			}}, got)
		})
	}
}

func TestSyncServiceGetUpNextNitro(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/sync/progress/up_next_nitro", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.AssertNoDiff(t, "genres=action%2Cdrama&intent=continue&page=1&watchnow=subscriptions&years=2020-2026", r.URL.RawQuery)
		test.SafeFprint(w, `[{"show":{"title":"Reacher"},"progress":{"aired":24,"completed":20}}]`)
	})

	got, _, err := setup.Client.Sync.GetUpNextNitro(context.Background(), &uri.UpNextNitroOptions{
		Page: 1, Intent: "continue", WatchNow: "subscriptions", Genres: "action,drama", Years: "2020-2026",
	})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.ShowProgress{{
		Show:     &str.Show{Title: str.String("Reacher")},
		Progress: &str.WatchedProgress{Aired: test.Ptr(24), Completed: test.Ptr(20)},
	}}, got)
}
