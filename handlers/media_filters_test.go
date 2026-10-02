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

// the media list actions send the media filter flags as query parameters; an unknown -watchnow stops before the request.
func TestMediaAndRecommendationsMediaFilters(t *testing.T) {
	const filtered = "certifications=pg-13&countries=us&end_date=2026-12-31&genres=action%2Cdrama&page=1&ratings=75-100&runtimes=90-150&start_date=2026-01-01&subgenres=space&watchnow=free&years=2020-2026"
	filters := str.Options{
		WatchNow: "free", Genres: "action,drama", Subgenres: "space", Years: "2020-2026", Ratings: "75-100", Runtimes: "90-150", Countries: "us",
		Certifications: "pg-13", MediaStartDate: "2026-01-01", MediaEndDate: "2026-12-31",
	}
	handlers := []struct {
		name    string
		handler Handler
		path    string
	}{
		{name: "media trending", handler: MediaTrendingHandler{}, path: "/media/trending"},
		{name: "media anticipated", handler: MediaAnticipatedHandler{}, path: "/media/anticipated"},
		{name: "media popular", handler: MediaPopularHandler{}, path: "/media/popular"},
		{name: "recommendations movies", handler: RecommendationsMoviesHandler{}, path: "/recommendations/movies"},
		{name: "recommendations shows", handler: RecommendationsShowsHandler{}, path: "/recommendations/shows"},
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
				s.Mux.HandleFunc(h.path, func(w http.ResponseWriter, r *http.Request) {
					test.AssertMethod(t, r, http.MethodGet)
					queries = append(queries, r.URL.RawQuery)
					test.SafeFprint(w, `[{"title":"Andor","year":2022,"ids":{"trakt":1}}]`)
				})

				options := tt.options
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
