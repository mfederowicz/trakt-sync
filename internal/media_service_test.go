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

func TestMediaServiceGetMediaItems(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		call func(s *MediaService, opts *uri.ListOptions) ([]*str.MediaItem, *str.Response, error)
		want []*str.MediaItem
	}{
		{
			name: "trending",
			path: "/media/trending",
			body: `[{"watchers":21,"movie":{"title":"Tron: Ares"}},{"watchers":15,"show":{"title":"Andor"}}]`,
			call: func(s *MediaService, opts *uri.ListOptions) ([]*str.MediaItem, *str.Response, error) {
				return s.GetTrendingMedia(context.Background(), opts)
			},
			want: []*str.MediaItem{
				{Watchers: test.Ptr(21), Movie: &str.Movie{Title: str.String("Tron: Ares")}},
				{Watchers: test.Ptr(15), Show: &str.Show{Title: str.String("Andor")}},
			},
		},
		{
			name: "anticipated",
			path: "/media/anticipated",
			body: `[{"list_count":120,"show":{"title":"Severance"}},{"list_count":98,"movie":{"title":"Dune: Part Three"}}]`,
			call: func(s *MediaService, opts *uri.ListOptions) ([]*str.MediaItem, *str.Response, error) {
				return s.GetAnticipatedMedia(context.Background(), opts)
			},
			want: []*str.MediaItem{
				{ListCount: test.Ptr(120), Show: &str.Show{Title: str.String("Severance")}},
				{ListCount: test.Ptr(98), Movie: &str.Movie{Title: str.String("Dune: Part Three")}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				assertMediaQuery(t, r)
				test.SafeFprint(w, tt.body)
			})

			got, _, err := tt.call(setup.Client.Media, &uri.ListOptions{Page: 2, Limit: 10, Extended: "full"})
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, tt.want, got)
		})
	}
}

func TestMediaServiceGetPopularMedia(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/media/popular", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		assertMediaQuery(t, r)
		test.SafeFprint(w, `[{"title":"Tron: Ares","year":2025,"released":"2025-10-10"},{"title":"Andor","year":2022,"aired_episodes":24}]`)
	})

	got, _, err := setup.Client.Media.GetPopularMedia(context.Background(), &uri.ListOptions{Page: 2, Limit: 10, Extended: "full"})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.Media{
		{Title: str.String("Tron: Ares"), Year: test.Ptr(2025), Released: str.String("2025-10-10")},
		{Title: str.String("Andor"), Year: test.Ptr(2022), AiredEpisodes: test.Ptr(24)},
	}, got)
}

func assertMediaQuery(t *testing.T, r *http.Request) {
	t.Helper()
	q := r.URL.Query()
	for key, want := range map[string]string{"page": "2", "limit": "10", "extended": "full"} {
		if got := q.Get(key); got != want {
			t.Errorf("query %s is %q, want %q", key, got, want)
		}
	}
}
