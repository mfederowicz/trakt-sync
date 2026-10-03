// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestListsTrendingPopularHandlersRoute(t *testing.T) {
	tests := []struct {
		name     string
		handler  Handler
		listType string
		path     string
	}{
		{name: "trending", handler: ListsTrendingHandler{}, path: "/lists/trending"},
		{name: "trending by type", handler: ListsTrendingHandler{}, listType: "personal", path: "/lists/trending/personal"},
		{name: "popular", handler: ListsPopularHandler{}, path: "/lists/popular"},
		{name: "popular by type", handler: ListsPopularHandler{}, listType: "official", path: "/lists/popular/official"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/lists/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.URL.Path]++
				test.SafeFprint(w, `[{"like_count":1}]`)
			})

			options := &str.Options{Type: tt.listType, Output: filepath.Join(t.TempDir(), "out.json")}
			test.AssertNilError(t, tt.handler.Handle(options, s.Client))
			test.AssertNoDiff(t, map[string]int{tt.path: 1}, calls)
		})
	}
}

func TestListsItemsHandlerDefaultsToAllItemTypes(t *testing.T) {
	tests := []struct {
		name     string
		itemType string
		path     string
	}{
		{name: "no type", path: "/lists/55/items/" + consts.ListItemsAll},
		{name: "movie", itemType: "movie", path: "/lists/55/items/movie"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			called := false
			s.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				called = true
				if got, want := r.URL.Query().Get("sort_by"), "rank"; got != want {
					t.Errorf("sort_by is %q, want %q", got, want)
				}
				test.SafeFprint(w, `[{"rank":1}]`)
			})

			options := &str.Options{InternalID: "55", Type: tt.itemType, SortBy: "rank", SortHow: "asc", Output: filepath.Join(t.TempDir(), "out.json")}
			test.AssertNilError(t, ListsItemsHandler{}.Handle(options, s.Client))
			if !called {
				t.Errorf("%s was not called", tt.path)
			}
		})
	}
}

func TestListsReportHandler(t *testing.T) {
	tests := []struct {
		name    string
		options str.Options
		wantErr string
	}{
		{name: "report", options: str.Options{InternalID: "55", Reason: "spam"}},
		{name: "missing list id", options: str.Options{Reason: "spam"}, wantErr: consts.EmptyListIDMsg},
		{name: "missing reason", options: str.Options{InternalID: "55"}, wantErr: consts.EmptyReasonMsg},
		{name: "invalid reason", options: str.Options{InternalID: "55", Reason: "boring"}, wantErr: "reason 'boring' is not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			s.Mux.HandleFunc("/lists/55/report", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodPost)
				w.WriteHeader(http.StatusCreated)
			})

			options := tt.options
			options.Module = consts.Lists
			options.Action = "report"
			err := ListsReportHandler{}.Handle(&options, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			test.AssertNilError(t, err)
		})
	}
}

// trending, popular and items send the media filter flags; an unknown -watchnow stops before the request.
func TestListsHandlersMediaFilters(t *testing.T) {
	const filtered = "certifications=pg-13&countries=us&end_date=2026-12-31&genres=action%2Cdrama&languages=en%2Cpl&page=1&ratings=75-100&runtimes=90-150&start_date=2026-01-01&subgenres=space&watchnow=free&years=2020-2026"
	filters := str.Options{
		WatchNow: "free", Genres: "action,drama", Subgenres: "space", Years: "2020-2026", Ratings: "75-100", Runtimes: "90-150", Countries: "us",
		Certifications: "pg-13", MediaStartDate: "2026-01-01", MediaEndDate: "2026-12-31", Languages: "en,pl",
	}
	handlers := []struct {
		name    string
		handler Handler
		path    string
		body    string
	}{
		{name: "trending", handler: ListsTrendingHandler{}, path: "/lists/trending", body: `[{"like_count":1,"list":{"name":"Top"}}]`},
		{name: "popular", handler: ListsPopularHandler{}, path: "/lists/popular", body: `[{"like_count":1,"list":{"name":"Top"}}]`},
		{name: "items", handler: ListsItemsHandler{}, path: "/lists/55/items/" + consts.ListItemsAll, body: `[{"rank":1}]`},
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
				requests := []string{}
				s.Mux.HandleFunc("/lists/", func(w http.ResponseWriter, r *http.Request) {
					test.AssertMethod(t, r, http.MethodGet)
					requests = append(requests, r.URL.Path+"?"+r.URL.RawQuery)
					test.SafeFprint(w, h.body)
				})

				options := tt.options
				options.InternalID = "55"
				options.Output = filepath.Join(t.TempDir(), "out.json")
				err := h.handler.Handle(&options, s.Client)
				if tt.wantErr != "" {
					if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
						t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
					}
					if len(requests) != consts.ZeroValue {
						t.Errorf("requests are %v, want none", requests)
					}
					return
				}
				test.AssertNilError(t, err)
				test.AssertNoDiff(t, []string{h.path + "?" + tt.query}, requests)
			})
		}
	}
}
