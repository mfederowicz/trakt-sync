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

// the shows list actions send the media filter flags and -status as query parameters.
func TestShowsHandlersMediaFilters(t *testing.T) {
	const filtered = "certifications=tv-14&countries=us&end_date=2026-12-31&genres=action%2Cdrama&imdb_ratings=8.5-10.0&languages=en%2Cpl&page=1&ratings=75-100" +
		"&rt_meters=90-100&rt_user_meters=80-100&runtimes=30-60" +
		"&start_date=2026-01-01&statuses=returning+series%2Cended&subgenres=space&watchnow=free&years=2020-2026"
	filters := str.Options{
		WatchNow: "free", Genres: "action,drama", Subgenres: "space", Years: "2020-2026", Ratings: "75-100", Runtimes: "30-60", Countries: "us",
		Certifications: "tv-14", MediaStartDate: "2026-01-01", MediaEndDate: "2026-12-31", ShowStatus: "returning series,ended", Languages: "en,pl",
		ImdbRatings: "8.5-10.0", RtMeters: "90-100", RtUserMeters: "80-100",
	}
	handlers := []struct {
		name    string
		handler Handler
		path    string
	}{
		{name: "trending", handler: ShowsTrendingHandler{}, path: "/shows/trending"},
		{name: "anticipated", handler: ShowsAnticipatedHandler{}, path: "/shows/anticipated"},
		{name: "popular", handler: ShowsPopularHandler{}, path: "/shows/popular"},
		{name: "watched", handler: ShowsWatchedHandler{}, path: "/shows/watched/weekly"},
		{name: "favorited", handler: ShowsFavoritedHandler{}, path: "/shows/favorited/weekly"},
		{name: "played", handler: ShowsPlayedHandler{}, path: "/shows/played/weekly"},
		{name: "collected", handler: ShowsCollectedHandler{}, path: "/shows/collected/weekly"},
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
		{name: "unknown status", options: str.Options{ShowStatus: "ended,running"}, wantErr: "status 'running' is not valid"},
	}

	for _, h := range handlers {
		for _, tt := range cases {
			h, tt := h, tt
			t.Run(h.name+" "+tt.name, func(t *testing.T) {
				s := setup(t)
				defer s.Teardown()
				queries := []string{}
				s.Mux.HandleFunc("/shows/", func(w http.ResponseWriter, r *http.Request) {
					test.AssertMethod(t, r, http.MethodGet)
					assert.Equal(t, h.path, r.URL.Path)
					queries = append(queries, r.URL.RawQuery)
					test.SafeFprint(w, `[{"title":"Dark","show":{"title":"Dark"}}]`)
				})

				options := tt.options
				options.Module = "shows"
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
