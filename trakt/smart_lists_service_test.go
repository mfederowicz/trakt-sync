// Package trakt is a Go client for the Trakt API
package trakt

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
			`"filters":{"genres":["science-fiction"],"genres_operator":"and","studios":["pixar","-dreamworks-animation"],"people":["steven-spielberg"],`+
			`"people_operator":"or","people_role":"directing","years":[2000,2026],"released_within_days":365,"imdb_ratings":[7.5,10],"tmdb_ratings":[7,10],`+
			`"metascores":[70,100],"votes":[1000],"imdb_votes":[5000,900000],"ignore_watched":true}}`)
	})

	got, _, err := setup.Client.SmartLists.GetSmartList(context.Background(), "top-sci-fi")
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
			Genres:             []string{"science-fiction"},
			GenresOperator:     str.String("and"),
			Studios:            []string{"pixar", "-dreamworks-animation"},
			People:             []string{"steven-spielberg"},
			PeopleOperator:     str.String("or"),
			PeopleRole:         str.String("directing"),
			Years:              []int{2000, 2026},
			ReleasedWithinDays: test.Ptr(365),
			ImdbRatings:        []float64{7.5, 10},
			TmdbRatings:        []float64{7, 10},
			Metascores:         []int{70, 100},
			Votes:              []int{1000},
			ImdbVotes:          []int{5000, 900000},
			IgnoreWatched:      test.Ptr(true),
		},
	}, got)
}

func TestSmartListsServiceGetSmartListItems(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/smart-lists/top-sci-fi/items", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		want := "certifications=pg-13&countries=us&extended=full&genres=action&ignore_watched=true&ignore_watchlisted=false&imdb_votes=5000-&limit=10" +
			"&metascores=70-100&page=2&people=steven-spielberg%2C-tom-hanks&people_operator=and&people_role=directing&ratings=75-100&released_within_days=365" +
			"&runtimes=90-150&studios=pixar%2C-dreamworks-animation&subgenres=space&tmdb_ratings=7.5-10&votes=1000-&watchnow=free&years=2020-2026"
		if got := r.URL.RawQuery; got != want {
			t.Errorf("query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `[{"rank":1,"type":"movie","movie":{"title":"Arrival","year":2016}}]`)
	})

	opts := &uri.SmartListItemsOptions{
		Page: 2, Limit: 10, Extended: "full", WatchNow: "free", Genres: "action", Subgenres: "space", Years: "2020-2026",
		Ratings: "75-100", Runtimes: "90-150", Countries: "us", Certifications: "pg-13", IgnoreWatched: "true", IgnoreWatchlisted: "false",
		Studios: "pixar,-dreamworks-animation", People: "steven-spielberg,-tom-hanks", PeopleOperator: "and", PeopleRole: "directing", ReleasedWithinDays: 365,
		TmdbRatings: "7.5-10", Metascores: "70-100", Votes: "1000-", ImdbVotes: "5000-",
	}
	got, _, err := setup.Client.SmartLists.GetSmartListItems(context.Background(), "top-sci-fi", opts)
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.UserListItem{
		{Rank: test.Ptr(1), Type: str.String("movie"), Movie: &str.Movie{Title: str.String("Arrival"), Year: test.Ptr(2016)}},
	}, got)
}

func TestSmartListsServiceGetSmartListItemsByTypeAndSort(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/smart-lists/top-sci-fi/items/movies/imdb_rating/desc", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		want := "limit=10&page=2&parental_alcohol=0-3&parental_frightening=1-2&parental_include_unrated=true&parental_nudity=0-1" +
			"&parental_profanity=0-1&parental_violence=0-2&watchnow=free&watchnow_country=pl"
		if got := r.URL.RawQuery; got != want {
			t.Errorf("query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `[{"rank":1,"type":"movie","movie":{"title":"Arrival","year":2016}}]`)
	})

	opts := &uri.SmartListItemsOptions{
		Page: 2, Limit: 10, WatchNow: "free", WatchNowCountry: "pl", ParentalNudity: "0-1", ParentalViolence: "0-2", ParentalProfanity: "0-1",
		ParentalAlcohol: "0-3", ParentalFrightening: "1-2", ParentalIncludeUnrated: true,
	}
	got, _, err := setup.Client.SmartLists.GetSmartListItemsByTypeAndSort(context.Background(), "top-sci-fi", "movies", "imdb_rating", "desc", opts)
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

	got, resp, err := setup.Client.SmartLists.GetSmartList(context.Background(), "private-list")
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

func TestSmartListsServiceGetSmartListTemplates(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/smart-lists/templates", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.AssertNoDiff(t, "extended=images", r.URL.RawQuery)
		test.SafeFprint(w, `[{"id":"trending-sci-fi","source":"trending","media_type":"movies","filters":{"genres":["science-fiction"]},"personal":false,"total":120,`+
			`"items":[{"rank":1,"type":"movie","movie":{"title":"Arrival","year":2016}}]},`+
			`{"id":"up-next","source":"watchlist","media_type":"shows","filters":{"ignore_watched":true},"personal":true,"total":null,"items":[]}]`)
	})

	got, _, err := setup.Client.SmartLists.GetSmartListTemplates(context.Background(), &uri.ListOptions{Extended: "images"})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.SmartListTemplate{
		{
			ID: str.String("trending-sci-fi"), Source: str.String("trending"), MediaType: str.String("movies"),
			Filters:  &str.SmartListFilters{Genres: []string{"science-fiction"}},
			Personal: test.Ptr(false), Total: test.Ptr(120),
			Items: []*str.UserListItem{
				{Rank: test.Ptr(1), Type: str.String("movie"), Movie: &str.Movie{Title: str.String("Arrival"), Year: test.Ptr(2016)}},
			},
		},
		{
			ID: str.String("up-next"), Source: str.String("watchlist"), MediaType: str.String("shows"),
			Filters:  &str.SmartListFilters{IgnoreWatched: test.Ptr(true)},
			Personal: test.Ptr(true), Items: []*str.UserListItem{},
		},
	}, got)
}
