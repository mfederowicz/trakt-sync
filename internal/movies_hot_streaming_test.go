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

func TestMoviesServiceHotAndStreaming(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		call func(s *MoviesService, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error)
		want []*str.MoviesItem
	}{
		{
			name: "hot",
			path: "/movies/hot",
			body: `[{"list_count":120,"movie":{"title":"Tron: Ares"}}]`,
			call: func(s *MoviesService, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
				return s.GetHotMovies(context.Background(), opts)
			},
			want: []*str.MoviesItem{{ListCount: test.Ptr(120), Movie: &str.Movie{Title: str.String("Tron: Ares")}}},
		},
		{
			name: "streaming",
			path: "/movies/streaming/daily",
			body: `[{"rank":1,"delta":-2,"movie":{"title":"Weapons"}}]`,
			call: func(s *MoviesService, opts *uri.ListOptions) ([]*str.MoviesItem, *str.Response, error) {
				return s.GetStreamingMovies(context.Background(), str.String("daily"), opts)
			},
			want: []*str.MoviesItem{{Rank: test.Ptr(1), Delta: test.Ptr(-2), Movie: &str.Movie{Title: str.String("Weapons")}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				if got, want := r.URL.Query().Get("page"), "2"; got != want {
					t.Errorf("page query is %q, want %q", got, want)
				}
				test.SafeFprint(w, tt.body)
			})

			got, _, err := tt.call(setup.Client.Movies, &uri.ListOptions{Page: 2})
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, tt.want, got)
		})
	}
}
