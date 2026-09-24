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

func TestSearchServiceGetExactTextQueryResults(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/search/movie/exact", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.AssertNoDiff(t, "limit=10&page=2&query=tron", r.URL.RawQuery)
		test.SafeFprint(w, `[{"type":"movie","score":1000,"movie":{"title":"TRON","ids":{"trakt":1}}}]`)
	})

	got, _, err := setup.Client.Search.GetExactTextQueryResults(context.Background(), str.String("movie"), &uri.ListOptions{Page: 2, Limit: 10, Query: "tron"})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.SearchListItem{{
		Type:  str.String("movie"),
		Score: test.Ptr(float32(1000)),
		Movie: &str.Movie{Title: str.String("TRON"), IDs: &str.IDs{Trakt: test.Ptr(int64(1))}},
	}}, got)
}

func TestSearchServiceGetTrendingSearches(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/search/recent_by_id/global/people", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.AssertNoDiff(t, "page=1&query=keanu", r.URL.RawQuery)
		test.SafeFprint(w, `[{"id":12345,"count":7,"type":"person","person":{"name":"Keanu Reeves"}}]`)
	})

	got, _, err := setup.Client.Search.GetTrendingSearches(context.Background(), str.String("people"), &uri.ListOptions{Page: 1, Query: "keanu"})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.SearchTrendingItem{{
		ID:     test.Ptr(int64(12345)),
		Count:  test.Ptr(7),
		Type:   str.String("person"),
		Person: &str.Person{Name: str.String("Keanu Reeves")},
	}}, got)
}

func TestSearchServiceErrors(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/search/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, _, err := setup.Client.Search.GetExactTextQueryResults(context.Background(), str.String("show"), &uri.ListOptions{Query: "x"})
	if err == nil {
		t.Error("exact query: expected an error on 500")
	}
	_, _, err = setup.Client.Search.GetTrendingSearches(context.Background(), str.String("shows"), &uri.ListOptions{})
	if err == nil {
		t.Error("trending: expected an error on 500")
	}
}
