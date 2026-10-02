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

// watchlist, history, collection and list_items send the media filter flags; an unknown -watchnow stops before the request.
func TestUsersHandlersMediaFilters(t *testing.T) {
	filtered := str.Options{
		WatchNow: "free", Genres: "action,drama", Subgenres: "space", Years: "2020-2026", Ratings: "75-100", Runtimes: "90-150", Countries: "us",
		Certifications: "pg-13", MediaStartDate: "2026-01-01", MediaEndDate: "2026-12-31",
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
				for _, param := range []string{"certifications=pg-13", "countries=us", "end_date=2026-12-31", "genres=action%2Cdrama", "ratings=75-100",
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
