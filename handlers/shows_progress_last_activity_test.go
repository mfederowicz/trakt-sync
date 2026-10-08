// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestShowsProgressLastActivity(t *testing.T) {
	tests := []struct {
		name      string
		handler   Handler
		action    string
		value     string
		wantPath  string
		wantQuery string
		wantErr   string
	}{
		{name: "collection without the flag", handler: ShowsCollectionProgressHandler{}, action: "collection_progress", wantPath: "/shows/55/progress/collection"},
		{name: "collection collected", handler: ShowsCollectionProgressHandler{}, action: "collection_progress", value: "collected",
			wantPath: "/shows/55/progress/collection", wantQuery: "last_activity=collected"},
		{name: "collection aired", handler: ShowsCollectionProgressHandler{}, action: "collection_progress", value: "aired",
			wantPath: "/shows/55/progress/collection", wantQuery: "last_activity=aired"},
		{name: "collection watched is not valid", handler: ShowsCollectionProgressHandler{}, action: "collection_progress", value: "watched",
			wantErr: "last_activity 'watched' is not valid"},
		{name: "watched without the flag", handler: ShowsWatchedProgressHandler{}, action: "watched_progress", wantPath: "/shows/55/progress/watched"},
		{name: "watched watched", handler: ShowsWatchedProgressHandler{}, action: "watched_progress", value: "watched",
			wantPath: "/shows/55/progress/watched", wantQuery: "last_activity=watched"},
		{name: "watched collected is not valid", handler: ShowsWatchedProgressHandler{}, action: "watched_progress", value: "collected",
			wantErr: "last_activity 'collected' is not valid"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/shows/55/progress/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				if r.URL.Path != tt.wantPath {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.wantPath)
				}
				if got := r.URL.RawQuery; got != tt.wantQuery {
					t.Errorf("query is %q, want %q", got, tt.wantQuery)
				}
				test.SafeFprint(w, `{"aired":8,"completed":6,"last_episode":{"season":1,"number":6}}`)
			})

			options := &str.Options{Module: "shows", Action: tt.action, InternalID: "55", LastActivity: tt.value, Output: filepath.Join(t.TempDir(), "out.json")}
			err := tt.handler.Handle(options, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				if calls != 0 {
					t.Error("API was called with an invalid last_activity")
				}
				return
			}
			test.AssertNilError(t, err)
			if calls != 1 {
				t.Errorf("API was called %d times, want 1", calls)
			}
		})
	}
}
