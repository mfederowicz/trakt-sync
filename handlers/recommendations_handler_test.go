// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
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

// no recommendations ends with an error and writes no file, like the other list actions.
func TestRecommendationsHandlersEmptyResult(t *testing.T) {
	tests := []struct {
		name    string
		handler Handler
		action  string
		path    string
	}{
		{name: "movies", handler: RecommendationsMoviesHandler{}, action: consts.Movies, path: "/recommendations/movies"},
		{name: "shows", handler: RecommendationsShowsHandler{}, action: consts.Shows, path: "/recommendations/shows"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			s.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, _ *http.Request) {
				test.SafeFprint(w, `[]`)
			})

			options := str.Options{Action: tt.action, PerPage: 10, Output: filepath.Join(t.TempDir(), "out.json")}
			assert.EqualError(t, tt.handler.Handle(&options, s.Client), consts.EmptyResult)
			_, err := os.Stat(options.Output)
			assert.True(t, os.IsNotExist(err), "no output file is written")
		})
	}
}

// with -ex full the export keeps the extended fields of the movie or show, not only title, year and ids.
func TestRecommendationsHandlersExportExtendedFields(t *testing.T) {
	const movie = `{"title":"Ida","year":2013,"ids":{"trakt":1},"tagline":"t","overview":"o","released":"2013-10-25","runtime":82,"country":"pl",` +
		`"trailer":"https://example.com/t","homepage":"https://example.com","status":"released","rating":7.5,"votes":10,"comment_count":2,` +
		`"language":"pl","languages":["pl","la"],"available_translations":["en","pl"],"genres":["drama"],"certification":"PG-13"}`
	const show = `{"title":"1670","year":2023,"ids":{"trakt":2},"overview":"o","airs":{"day":"Wednesday","time":"09:00","timezone":"Europe/Warsaw"},` +
		`"runtime":30,"certification":"TV-MA","network":"Netflix","country":"pl","status":"returning series","rating":8,"votes":5,` +
		`"language":"pl","languages":["pl"],"genres":["comedy"],"aired_episodes":8}`
	tests := []struct {
		name    string
		handler Handler
		action  string
		path    string
		item    string
	}{
		{name: "movies", handler: RecommendationsMoviesHandler{}, action: consts.Movies, path: "/recommendations/movies", item: movie},
		{name: "shows", handler: RecommendationsShowsHandler{}, action: consts.Shows, path: "/recommendations/shows", item: show},
		{name: "social movies", handler: SocialRecommendationsMoviesHandler{}, action: consts.Movies, path: "/social_recommendations/movies", item: movie},
		{name: "social shows", handler: SocialRecommendationsShowsHandler{}, action: consts.Shows, path: "/social_recommendations/shows", item: show},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			s.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "full", r.URL.Query().Get("extended"))
				test.SafeFprint(w, "["+tt.item+"]")
			})

			options := str.Options{Action: tt.action, PerPage: 10, ExtendedInfo: "full", Output: filepath.Join(t.TempDir(), "out.json")}
			test.AssertNilError(t, tt.handler.Handle(&options, s.Client))

			data, err := os.ReadFile(options.Output)
			test.AssertNilError(t, err)
			got := []json.RawMessage{}
			test.AssertNilError(t, json.Unmarshal(data, &got))
			if assert.Len(t, got, 1) {
				assert.JSONEq(t, tt.item, string(got[0]))
			}
		})
	}
}
