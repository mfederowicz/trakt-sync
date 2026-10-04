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
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/stretchr/testify/assert"
)

// watchlist, history, collection and list_items send the media filter flags; an unknown -watchnow stops before the request.
func TestUsersHandlersMediaFilters(t *testing.T) {
	filtered := str.Options{
		WatchNow: "free", Genres: "action,drama", Subgenres: "space", Years: "2020-2026", Ratings: "75-100", Runtimes: "90-150", Countries: "us",
		Certifications: "pg-13", MediaStartDate: "2026-01-01", MediaEndDate: "2026-12-31", Languages: "en,pl",
	}
	handlers := []struct {
		action  string
		handler Handler
	}{
		{action: "watchlist", handler: UsersWatchlistHandler{}},
		{action: "history", handler: UsersHistoryHandler{}},
		{action: "collection", handler: UsersCollectionHandler{}},
		{action: "list_items", handler: UsersListItemsHandler{}},
	}

	for _, h := range handlers {
		h := h
		t.Run(h.action+" all filters", func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			queries := []string{}
			s.Mux.HandleFunc("/users/sean/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				queries = append(queries, r.URL.RawQuery)
				test.SafeFprint(w, `[{"id":1,"type":"movie","movie":{"title":"Tron","ids":{"trakt":1}}}]`)
			})

			options := filtered
			options.Module, options.Action, options.UserName, options.Type, options.ID = "users", h.action, "sean", "movies", "55"
			options.Output = filepath.Join(t.TempDir(), "out.json")
			test.AssertNilError(t, h.handler.Handle(&options, s.Client))
			if assert.Len(t, queries, 1) {
				got := queries[0]
				for _, param := range []string{"certifications=pg-13", "countries=us", "end_date=2026-12-31", "genres=action%2Cdrama", "languages=en%2Cpl", "ratings=75-100",
					"runtimes=90-150", "start_date=2026-01-01", "subgenres=space", "watchnow=free", "years=2020-2026"} {
					assert.Contains(t, got, param)
				}
			}
		})
		t.Run(h.action+" unknown watchnow", func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			requests := 0
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				requests++
				test.SafeFprint(w, `[]`)
			})

			options := str.Options{Module: "users", Action: h.action, UserName: "sean", Type: "movies", ID: "55", WatchNow: "cinema", Output: filepath.Join(t.TempDir(), "out.json")}
			assert.ErrorContains(t, h.handler.Handle(&options, s.Client), "watchnow 'cinema' is not valid")
			assert.Zero(t, requests, "no request is sent")
		})
	}
}

// The page count header of a personal list ignores the media filters, so an empty page ends the paging.
func TestUsersListItemsHandlerStopsOnEmptyPage(t *testing.T) {
	tests := []struct {
		name      string
		pages     map[string]string
		wantCalls int
		wantErr   string
	}{
		{name: "filter matches nothing", pages: map[string]string{}, wantCalls: 1, wantErr: "empty list items"},
		{name: "filter matches one page", pages: map[string]string{"1": `[{"rank":1}]`}, wantCalls: 2},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/users/me/lists/55/items/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				page := r.URL.Query().Get("page")
				w.Header().Set(trakt.HeaderPaginationPage, page)
				w.Header().Set(trakt.HeaderPaginationPageCount, "5")
				body, ok := tt.pages[page]
				if !ok {
					body = `[]`
				}
				test.SafeFprint(w, body)
			})

			options := &str.Options{Module: "users", Action: "list_items", UserName: "me", ID: "55", Type: "movies", Languages: "pl", Output: filepath.Join(t.TempDir(), "out.json")}
			err := UsersListItemsHandler{}.Handle(options, s.Client)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				test.AssertNilError(t, err)
			}
			assert.Equal(t, tt.wantCalls, calls, "API calls")
		})
	}
}

// The page count header of the watchlist ignores the media filters, so an empty page ends the paging.
func TestUsersWatchlistHandlerStopsOnEmptyPage(t *testing.T) {
	const item = `[{"id":1,"type":"movie","movie":{"title":"Tron","ids":{"trakt":1}}}]`
	tests := []struct {
		name      string
		sortPath  string
		pages     map[string]string
		wantCalls int
	}{
		{name: "filter matches nothing", pages: map[string]string{}, wantCalls: 1},
		{name: "filter matches one page", pages: map[string]string{"1": item}, wantCalls: 2},
		{name: "sort route, filter matches one page", sortPath: "added", pages: map[string]string{"1": item}, wantCalls: 2},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/users/sean/watchlist/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				page := r.URL.Query().Get("page")
				w.Header().Set(trakt.HeaderPaginationPage, page)
				w.Header().Set(trakt.HeaderPaginationPageCount, "5")
				body, ok := tt.pages[page]
				if !ok {
					body = `[]`
				}
				test.SafeFprint(w, body)
			})

			options := &str.Options{Module: "users", Action: "watchlist", UserName: "sean", Type: "movies", SortPath: tt.sortPath, Languages: "pl",
				Output: filepath.Join(t.TempDir(), "out.json")}
			err := UsersWatchlistHandler{}.Handle(options, s.Client)
			if len(tt.pages) > 0 {
				test.AssertNilError(t, err)
			} else {
				assert.EqualError(t, err, consts.EmptyResult)
				_, statErr := os.Stat(options.Output)
				assert.True(t, os.IsNotExist(statErr), "no output file is written")
			}
			assert.Equal(t, tt.wantCalls, calls, "API calls")
		})
	}
}

// history, favorites and ratings end with an error and write no file when the API returns no items, like watchlist.
func TestUsersHandlersEmptyResult(t *testing.T) {
	handlers := []struct {
		action  string
		handler Handler
	}{
		{action: "history", handler: UsersHistoryHandler{}},
		{action: "favorites", handler: UsersFavoritesHandler{}},
		{action: "ratings", handler: UsersRatingsHandler{}},
	}

	for _, h := range handlers {
		h := h
		t.Run(h.action, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			s.Mux.HandleFunc("/users/sean/"+h.action+"/", func(w http.ResponseWriter, _ *http.Request) {
				test.SafeFprint(w, `[]`)
			})

			options := &str.Options{Module: "users", Action: h.action, UserName: "sean", Type: "movies", Output: filepath.Join(t.TempDir(), "out.json")}
			assert.EqualError(t, h.handler.Handle(options, s.Client), consts.EmptyResult)
			_, err := os.Stat(options.Output)
			assert.True(t, os.IsNotExist(err), "no output file is written")
		})
	}
}

// watchlist sends -hide on both of its routes; an unknown value stops before the request.
func TestUsersWatchlistHandlerHide(t *testing.T) {
	tests := []struct {
		name      string
		options   str.Options
		want      string
		wantErr   string
		wantCalls int
	}{
		{name: "typed route", options: str.Options{HideItems: "rated"}, want: "hide=rated", wantCalls: 1},
		{name: "sorted route", options: str.Options{HideItems: "airing", SortPath: "rank"}, want: "hide=airing", wantCalls: 1},
		{name: "not set", options: str.Options{}, wantCalls: 1},
		{name: "unknown value", options: str.Options{HideItems: "boring"}, wantErr: "hide 'boring' is not valid"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			queries := []string{}
			s.Mux.HandleFunc("/users/sean/watchlist/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				queries = append(queries, r.URL.RawQuery)
				test.SafeFprint(w, `[{"id":1,"type":"movie","movie":{"title":"Tron","ids":{"trakt":1}}}]`)
			})

			options := tt.options
			options.Module, options.Action, options.UserName, options.Type = "users", "watchlist", "sean", "movies"
			options.Output = filepath.Join(t.TempDir(), "out.json")
			err := UsersWatchlistHandler{}.Handle(&options, s.Client)
			if tt.wantErr != consts.EmptyString {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				test.AssertNilError(t, err)
			}
			if !assert.Len(t, queries, tt.wantCalls) || tt.wantCalls == consts.ZeroValue {
				return
			}
			if tt.want == consts.EmptyString {
				assert.NotContains(t, queries[0], "hide")
				return
			}
			assert.Contains(t, queries[0], tt.want)
		})
	}
}
