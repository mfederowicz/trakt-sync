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
