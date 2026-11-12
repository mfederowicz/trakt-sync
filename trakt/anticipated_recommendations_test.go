package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestGetAnticipatedRecommendations(t *testing.T) {
	const body = `[{"score":0.92,"released":"2026-12-18","movie":{"title":"Avatar: Fire and Ash","ids":{"trakt":1}}},` +
		`{"score":0.81,"released":null,"show":{"title":"Andor","ids":{"trakt":2}}}]`
	want := []*str.AnticipatedRecommendation{
		{Score: test.Ptr(0.92), Released: str.String("2026-12-18"), Movie: &str.Movie{Title: str.String("Avatar: Fire and Ash"), IDs: &str.IDs{Trakt: test.Ptr(int64(1))}}},
		{Score: test.Ptr(0.81), Show: &str.Show{Title: str.String("Andor"), IDs: &str.IDs{Trakt: test.Ptr(int64(2))}}},
	}
	ctx := context.Background()
	opts := &uri.ListOptions{Page: 2, Limit: 50}
	cases := []struct {
		name string
		path string
		call func(c *Client) ([]*str.AnticipatedRecommendation, *str.Response, error)
	}{
		{name: "movies", path: "/movies/recommendations/anticipated", call: func(c *Client) ([]*str.AnticipatedRecommendation, *str.Response, error) {
			return c.Movies.GetAnticipatedRecommendations(ctx, opts)
		}},
		{name: "shows", path: "/shows/recommendations/anticipated", call: func(c *Client) ([]*str.AnticipatedRecommendation, *str.Response, error) {
			return c.Shows.GetAnticipatedRecommendations(ctx, opts)
		}},
		{name: "media", path: "/media/recommendations/anticipated", call: func(c *Client) ([]*str.AnticipatedRecommendation, *str.Response, error) {
			return c.Media.GetAnticipatedRecommendations(ctx, opts)
		}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()
			setup.Mux.HandleFunc(tc.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				if got := r.URL.RawQuery; got != "limit=50&page=2" {
					t.Errorf("query is %q, want %q", got, "limit=50&page=2")
				}
				test.SafeFprint(w, body)
			})

			got, _, err := tc.call(setup.Client)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, want, got)
		})
	}
}

func TestGetAnticipatedRecommendationsUnauthorized(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	setup.Mux.HandleFunc("/movies/recommendations/anticipated", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		test.SafeFprint(w, `{"error":"unauthorized"}`)
	})

	_, resp, err := setup.Client.Movies.GetAnticipatedRecommendations(context.Background(), nil)
	if err == nil {
		t.Fatal("want an error for 401")
	}
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("want status 401, got %v", resp)
	}
}
