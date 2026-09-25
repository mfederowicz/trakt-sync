// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestSocialRecommendationsHandlers(t *testing.T) {
	tests := []struct {
		name    string
		handler Handler
		action  string
		path    string
	}{
		{name: "movies", handler: SocialRecommendationsMoviesHandler{}, action: consts.Movies, path: "/social_recommendations/movies"},
		{name: "shows", handler: SocialRecommendationsShowsHandler{}, action: consts.Shows, path: "/social_recommendations/shows"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			s.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				q := r.URL.Query()
				for key, want := range map[string]string{"limit": "10", "extended": "images", "ignore_watched": "true", "watch_window": "7"} {
					if got := q.Get(key); got != want {
						t.Errorf("query %s is %q, want %q", key, got, want)
					}
				}
				test.SafeFprint(w, `[{"title":"Tron: Ares","year":2025,"ids":{"trakt":7}}]`)
			})

			output := filepath.Join(t.TempDir(), "export_social_recommendations_"+tt.action+".json")
			options := &str.Options{Action: tt.action, PerPage: 10, ExtendedInfo: "images", IgnoreWatched: "true", WatchWindow: 7, Output: output}
			test.AssertNilError(t, tt.handler.Handle(options, s.Client))

			data, err := os.ReadFile(output)
			test.AssertNilError(t, err)
			got := []*str.Recommendation{}
			test.AssertNilError(t, json.Unmarshal(data, &got))
			test.AssertNoDiff(t, []*str.Recommendation{{Title: str.String("Tron: Ares"), Year: test.Ptr(2025), IDs: &str.IDs{Trakt: test.Ptr(int64(7))}}}, got)
		})
	}
}

func TestSocialRecommendationsHandlerErrors(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr string
	}{
		{name: "empty list", status: http.StatusOK, body: `[]`, wantErr: consts.EmptyResult},
		{name: "unauthorized", status: http.StatusUnauthorized, body: `{}`, wantErr: "fetch social recommendations movies error"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			s.Mux.HandleFunc("/social_recommendations/movies", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			output := filepath.Join(t.TempDir(), "export_social_recommendations_movies.json")
			options := &str.Options{Action: consts.Movies, Output: output}
			err := SocialRecommendationsMoviesHandler{}.Handle(options, s.Client)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
			}
			if _, statErr := os.Stat(output); !os.IsNotExist(statErr) {
				t.Errorf("output file written on error: %v", statErr)
			}
		})
	}
}
