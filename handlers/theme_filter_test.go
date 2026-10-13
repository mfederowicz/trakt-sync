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

// the trending, popular, anticipated and recommendations actions send -theme; other list actions do not take it.
func TestHandlersTheme(t *testing.T) {
	const movies = `[{"title":"Weapons","movie":{"title":"Weapons"}}]`
	const shows = `[{"title":"Andor","show":{"title":"Andor"}}]`
	const media = `[{"title":"Andor","year":2022,"ids":{"trakt":1}}]`
	handlers := []struct {
		name    string
		handler Handler
		module  string
		action  string
		path    string
		body    string
		noTheme bool
	}{
		{name: "movies trending", handler: MoviesTrendingHandler{}, module: "movies", action: "trending", path: "/movies/trending", body: movies},
		{name: "movies popular", handler: MoviesPopularHandler{}, module: "movies", action: "popular", path: "/movies/popular", body: movies},
		{name: "movies anticipated", handler: MoviesAnticipatedHandler{}, module: "movies", action: "anticipated", path: "/movies/anticipated", body: movies},
		{name: "shows trending", handler: ShowsTrendingHandler{}, module: "shows", action: "trending", path: "/shows/trending", body: shows},
		{name: "shows popular", handler: ShowsPopularHandler{}, module: "shows", action: "popular", path: "/shows/popular", body: shows},
		{name: "shows anticipated", handler: ShowsAnticipatedHandler{}, module: "shows", action: "anticipated", path: "/shows/anticipated", body: shows},
		{name: "media trending", handler: MediaTrendingHandler{}, module: "media", action: "trending", path: "/media/trending", body: media},
		{name: "media popular", handler: MediaPopularHandler{}, module: "media", action: "popular", path: "/media/popular", body: media},
		{name: "media anticipated", handler: MediaAnticipatedHandler{}, module: "media", action: "anticipated", path: "/media/anticipated", body: media},
		{name: "recommendations movies", handler: RecommendationsMoviesHandler{}, module: "recommendations", action: "movies", path: "/recommendations/movies", body: media},
		{name: "recommendations shows", handler: RecommendationsShowsHandler{}, module: "recommendations", action: "shows", path: "/recommendations/shows", body: media},
		{name: "movies watched", handler: MoviesWatchedHandler{}, module: "movies", action: "watched", path: "/movies/watched/weekly", body: movies, noTheme: true},
		{name: "shows watched", handler: ShowsWatchedHandler{}, module: "shows", action: "watched", path: "/shows/watched/weekly", body: shows, noTheme: true},
	}
	cases := []struct {
		name    string
		options str.Options
		query   string
	}{
		{name: "theme", options: str.Options{Theme: "halloween"}, query: "page=1&theme=halloween"},
		{name: "theme with a filter", options: str.Options{Theme: "christmas", Genres: "comedy"}, query: "genres=comedy&page=1&theme=christmas"},
		{name: "no theme", query: "page=1"},
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
					test.SafeFprint(w, h.body)
				})

				options := tt.options
				options.Module = h.module
				options.Action = h.action
				options.Period = "weekly"
				options.Output = filepath.Join(t.TempDir(), "out.json")
				test.AssertNilError(t, h.handler.Handle(&options, s.Client))

				want := tt.query
				if h.noTheme {
					q := "page=1"
					if len(tt.options.Genres) > 0 {
						q = "genres=" + tt.options.Genres + "&page=1"
					}
					want = q
				}
				assert.Equal(t, []string{want}, queries)
			})
		}
	}
}
