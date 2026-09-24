// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestShowsServiceGetShowWatchNow(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/watchnow/us", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		if got, want := r.URL.RawQuery, "extended=streaming_ranks&links=tvos%2Cwebos"; got != want {
			t.Errorf("query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `{"us":{"cable":[],"free":[],"cinema":[],"purchase":[],
			"subscription":[{"source":"disney_plus","link":"watchnow.trakt.tv/watchnow/1","uhd":true,"currency":"usd","prices":{"rent":"3.99"},
				"link_tvos":"https://tv.example/1","link_webos":{"id":"com.example","params":{"contentTarget":"x"}}}],
			"streaming_ranks":{"rank":5787,"delta":-586,"link":"https://www.justwatch.com/us/tv-show/the-sopranos"}}}`)
	})

	got, _, err := setup.Client.Shows.GetShowWatchNow(context.Background(), str.String("the-sopranos"), str.String("us"), &uri.ListOptions{Extended: "streaming_ranks", Links: "tvos,webos"})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, map[string]*str.WatchNowSources{"us": {
		Cable:    []*str.WatchNowService{},
		Free:     []*str.WatchNowService{},
		Cinema:   []*str.WatchNowService{},
		Purchase: []*str.WatchNowService{},
		Subscription: []*str.WatchNowService{{
			Source:    str.String("disney_plus"),
			Link:      str.String("watchnow.trakt.tv/watchnow/1"),
			UHD:       test.Ptr(true),
			Currency:  str.String("usd"),
			Prices:    &str.WatchNowPrices{Rent: str.String("3.99")},
			LinkTvos:  str.String("https://tv.example/1"),
			LinkWebos: &str.WatchNowWebosLink{ID: str.String("com.example"), Params: &str.WatchNowWebosParams{ContentTarget: str.String("x")}},
		}},
		StreamingRanks: &str.WatchNowRank{Rank: test.Ptr(5787), Delta: test.Ptr(-586), Link: str.String("https://www.justwatch.com/us/tv-show/the-sopranos")},
	}}, got)
}

func TestShowsServiceGetShowJustwatchLinks(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/watchnow/justwatch_links/pl", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `{"pl":"justwatch.com/pl/serial/rodzina-soprano"}`)
	})

	got, _, err := setup.Client.Shows.GetShowJustwatchLinks(context.Background(), str.String("the-sopranos"), str.String("pl"))
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, map[string]string{"pl": "justwatch.com/pl/serial/rodzina-soprano"}, got)
}

func TestShowsServiceGetShowWatchNowForbidden(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/watchnow/justwatch_links/us", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	_, resp, err := setup.Client.Shows.GetShowJustwatchLinks(context.Background(), str.String("the-sopranos"), str.String("us"))
	var forbidden *ForbiddenError
	if !errors.As(err, &forbidden) {
		t.Fatalf("error is %v, want *ForbiddenError", err)
	}
	if got, want := resp.StatusCode, http.StatusForbidden; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}
