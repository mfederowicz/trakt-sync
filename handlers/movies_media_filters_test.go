// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// the movies list actions send the media filter flags as query parameters.
func TestMoviesHandlersMediaFilters(t *testing.T) {
	const filtered = "certifications=pg-13&countries=us&end_date=2026-12-31&genres=action%2Cdrama&page=1&ratings=75-100&runtimes=90-150&start_date=2026-01-01&subgenres=space&watchnow=free&years=2020-2026"
	filters := str.Options{WatchNow: "free", Genres: "action,drama", Subgenres: "space", Years: "2020-2026", Ratings: "75-100", Runtimes: "90-150", Countries: "us", Certifications: "pg-13", MediaStartDate: "2026-01-01", MediaEndDate: "2026-12-31"}
	handlers := []struct {
		name    string
		handler Handler
		path    string
	}{
		{name: "trending", handler: MoviesTrendingHandler{}, path: "/movies/trending"},
		{name: "anticipated", handler: MoviesAnticipatedHandler{}, path: "/movies/anticipated"},
		{name: "popular", handler: MoviesPopularHandler{}, path: "/movies/popular"},
		{name: "hot", handler: MoviesHotHandler{}, path: "/movies/hot"},
		{name: "streaming", handler: MoviesStreamingHandler{}, path: "/movies/streaming/weekly"},
	}
	cases := []struct {
		name    string
		options str.Options
		query   string
		wantErr string
	}{
		{name: "all filters", options: filters, query: filtered},
		{name: "no filters", query: "page=1"},
		{name: "unknown watchnow", options: str.Options{WatchNow: "cinema"}, wantErr: "watchnow 'cinema' is not valid"},
	}

	for _, h := range handlers {
		for _, tt := range cases {
			h, tt := h, tt
			t.Run(h.name+" "+tt.name, func(t *testing.T) {
				s := setup(t)
				defer s.Teardown()
				queries := []string{}
				s.Mux.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
					test.AssertMethod(t, r, http.MethodGet)
					assert.Equal(t, h.path, r.URL.Path)
					queries = append(queries, r.URL.RawQuery)
					test.SafeFprint(w, `[{"title":"Weapons","movie":{"title":"Weapons"}}]`)
				})

				options := tt.options
				options.Module = "movies"
				options.Action = h.name
				options.Period = "weekly"
				options.Output = filepath.Join(t.TempDir(), "out.json")
				err := h.handler.Handle(&options, s.Client)
				if tt.wantErr != "" {
					assert.ErrorContains(t, err, tt.wantErr)
					assert.Empty(t, queries, "no request is sent")
					return
				}
				test.AssertNilError(t, err)
				assert.Equal(t, []string{tt.query}, queries)
			})
		}
	}
}
