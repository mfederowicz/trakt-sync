// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestRecommendationsHandlersQuery(t *testing.T) {
	tests := []struct {
		name    string
		handler Handler
		action  string
		path    string
		options str.Options
		query   string
	}{
		{
			name: "movies default", handler: RecommendationsMoviesHandler{}, action: consts.Movies, path: "/recommendations/movies",
			options: str.Options{PerPage: 10},
			query:   "limit=10&page=1",
		},
		{
			name: "movies all flags", handler: RecommendationsMoviesHandler{}, action: consts.Movies, path: "/recommendations/movies",
			options: str.Options{PerPage: 10, IgnoreCollected: "true", IgnoreWatched: "true", IgnoreWatchlisted: "false", WatchWindow: 30},
			query:   "ignore_collected=true&ignore_watched=true&ignore_watchlisted=false&limit=10&page=1&watch_window=30",
		},
		{
			name: "shows all flags", handler: RecommendationsShowsHandler{}, action: consts.Shows, path: "/recommendations/shows",
			options: str.Options{PerPage: 10, IgnoreWatched: "true", WatchWindow: 7},
			query:   "ignore_watched=true&limit=10&page=1&watch_window=7",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				if r.URL.RawQuery != tt.query {
					t.Errorf("query is %q, want %q", r.URL.RawQuery, tt.query)
				}
				test.SafeFprint(w, `[{"title":"Andor","year":2022,"ids":{"trakt":1}}]`)
			})

			options := tt.options
			options.Action = tt.action
			options.Output = filepath.Join(t.TempDir(), "out.json")
			test.AssertNilError(t, tt.handler.Handle(&options, s.Client))
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
			if _, err := os.Stat(options.Output); err != nil {
				t.Errorf("output file was not written: %v", err)
			}
		})
	}
}
