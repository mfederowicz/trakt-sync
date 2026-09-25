// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestSocialRecommendationsService(t *testing.T) {
	tests := []struct {
		name string
		path string
		call func(s *SocialRecommendationsService, opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error)
	}{
		{
			name: "movies",
			path: "/social_recommendations/movies",
			call: func(s *SocialRecommendationsService, opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
				return s.GetSocialMovieRecommendations(context.Background(), opts)
			},
		},
		{
			name: "shows",
			path: "/social_recommendations/shows",
			call: func(s *SocialRecommendationsService, opts *uri.ListOptions) ([]*str.Recommendation, *str.Response, error) {
				return s.GetSocialShowRecommendations(context.Background(), opts)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				want := "extended=full&ignore_collected=true&ignore_watched=true&ignore_watchlisted=false&limit=10&watch_window=30"
				if got := r.URL.RawQuery; got != want {
					t.Errorf("query is %q, want %q", got, want)
				}
				test.SafeFprint(w, `[{"title":"Andor","year":2022,"ids":{"trakt":1},"recommended_by":[{"user":{"username":"sean"},"notes":"watch it"}]}]`)
			})

			opts := &uri.ListOptions{Limit: 10, Extended: "full", IgnoreCollected: "true", IgnoreWatched: "true", IgnoreWatchlisted: "false", WatchWindow: 30}
			got, _, err := tt.call(setup.Client.SocialRecommendations, opts)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, []*str.Recommendation{{
				Title:         str.String("Andor"),
				Year:          test.Ptr(2022),
				IDs:           &str.IDs{Trakt: test.Ptr(int64(1))},
				RecommendedBy: &[]str.UserNotes{{User: &str.UserProfile{Username: str.String("sean")}, Notes: str.String("watch it")}},
			}}, got)
		})
	}
}

func TestSocialRecommendationsServiceError(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/social_recommendations/movies", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	got, resp, err := setup.Client.SocialRecommendations.GetSocialMovieRecommendations(context.Background(), &uri.ListOptions{})
	if err == nil {
		t.Fatal("expected an error")
	}
	if got != nil {
		t.Errorf("list is %v, want nil", got)
	}
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("response is %v, want status %d", resp, http.StatusUnauthorized)
	}
}
