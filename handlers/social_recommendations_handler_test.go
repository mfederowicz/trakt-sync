// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
)

// Both social recommendations routes are locked: the handlers end with the notice, send nothing and write no file.
func TestSocialRecommendationsHandlers(t *testing.T) {
	const locked = " is locked since 2026-11-20: Trakt removed this endpoint from its API, the live API no longer serves it; no request was sent"
	tests := []struct {
		name    string
		handler Handler
		action  string
		route   string
	}{
		{name: "movies", handler: SocialRecommendationsMoviesHandler{}, action: consts.Movies, route: "social_recommendations/movies"},
		{name: "shows", handler: SocialRecommendationsShowsHandler{}, action: consts.Shows, route: "social_recommendations/shows"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/social_recommendations/", func(w http.ResponseWriter, _ *http.Request) {
				calls++
				w.WriteHeader(http.StatusNotFound)
			})

			output := filepath.Join(t.TempDir(), "export_social_recommendations_"+tt.action+".json")
			options := &str.Options{Action: tt.action, PerPage: 10, WatchWindow: 7, Output: output}
			err := tt.handler.Handle(options, s.Client)
			want := "GET " + tt.route + locked
			if err == nil || err.Error() != "fetch social recommendations "+tt.action+" error: "+want {
				t.Fatalf("error is %v, want it to end with %q", err, want)
			}
			if calls != 0 {
				t.Errorf("%d requests sent, want 0", calls)
			}
			if _, statErr := os.Stat(output); !os.IsNotExist(statErr) {
				t.Errorf("output file written: %v", statErr)
			}
		})
	}
}
