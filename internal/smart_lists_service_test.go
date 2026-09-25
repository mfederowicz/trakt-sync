// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestSmartListsServiceGetSmartList(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/smart-lists/top-sci-fi", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `{"name":"Top Sci-Fi","privacy":"public","created_at":"2026-09-01T10:00:00.000Z","updated_at":"2026-09-02T10:00:00.000Z",`+
			`"ids":{"trakt":42,"slug":"top-sci-fi"},"images":{"posters":["p1.jpg"]},"source":"trakt","media_type":"movie",`+
			`"filters":{"genres":["science-fiction"],"genres_operator":"and","years":[2000,2026],"imdb_ratings":[7.5,10],"ignore_watched":true}}`)
	})

	got, _, err := setup.Client.SmartLists.GetSmartList(context.Background(), str.String("top-sci-fi"))
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, &str.SmartList{
		Name:      str.String("Top Sci-Fi"),
		Privacy:   str.String("public"),
		CreatedAt: &str.Timestamp{Time: time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)},
		UpdatedAt: &str.Timestamp{Time: time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC)},
		IDs:       &str.IDs{Trakt: test.Ptr(int64(42)), Slug: str.String("top-sci-fi")},
		Images:    &str.SmartListImages{Posters: []string{"p1.jpg"}},
		Source:    str.String("trakt"),
		MediaType: str.String("movie"),
		Filters: &str.SmartListFilters{
			Genres:         []string{"science-fiction"},
			GenresOperator: str.String("and"),
			Years:          []int{2000, 2026},
			ImdbRatings:    []float64{7.5, 10},
			IgnoreWatched:  test.Ptr(true),
		},
	}, got)
}

func TestSmartListsServiceGetSmartListItems(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/smart-lists/top-sci-fi/items", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		want := "certifications=pg-13&countries=us&extended=full&genres=action&ignore_watched=true&ignore_watchlisted=false&limit=10&page=2" +
			"&ratings=75-100&runtimes=90-150&subgenres=space&watchnow=free&years=2020-2026"
		if got := r.URL.RawQuery; got != want {
			t.Errorf("query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `[{"rank":1,"type":"movie","movie":{"title":"Arrival","year":2016}}]`)
	})

	opts := &uri.SmartListItemsOptions{
		Page: 2, Limit: 10, Extended: "full", WatchNow: "free", Genres: "action", Subgenres: "space", Years: "2020-2026",
		Ratings: "75-100", Runtimes: "90-150", Countries: "us", Certifications: "pg-13", IgnoreWatched: "true", IgnoreWatchlisted: "false",
	}
	got, _, err := setup.Client.SmartLists.GetSmartListItems(context.Background(), str.String("top-sci-fi"), opts)
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.UserListItem{
		{Rank: test.Ptr(1), Type: str.String("movie"), Movie: &str.Movie{Title: str.String("Arrival"), Year: test.Ptr(2016)}},
	}, got)
}

func TestSmartListsServiceNotFound(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/smart-lists/private-list", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	got, resp, err := setup.Client.SmartLists.GetSmartList(context.Background(), str.String("private-list"))
	if err == nil {
		t.Fatal("expected an error")
	}
	if got != nil {
		t.Errorf("result is %v, want nil", got)
	}
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		t.Errorf("response is %v, want status %d", resp, http.StatusNotFound)
	}
}
