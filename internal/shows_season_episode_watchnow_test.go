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

func TestShowsServiceGetEpisodeWatchNow(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/seasons/1/episodes/2/watchnow/us", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		if got, want := r.URL.RawQuery, "extended=streaming_ranks&links=direct"; got != want {
			t.Errorf("query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `{"us":{"subscription":[{"source":"max","link":"watchnow.trakt.tv/watchnow/1","uhd":false,"prices":{},"link_direct":"https://play.example/2"}]}}`)
	})

	got, _, err := setup.Client.Shows.GetEpisodeWatchNow(context.Background(), str.String("the-sopranos"), test.Ptr(1), test.Ptr(2), str.String("us"), &uri.ListOptions{Extended: "streaming_ranks", Links: "direct"})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, map[string]*str.WatchNowSources{"us": {
		Subscription: []*str.WatchNowService{{
			Source:     str.String("max"),
			Link:       str.String("watchnow.trakt.tv/watchnow/1"),
			UHD:        test.Ptr(false),
			Prices:     &str.WatchNowPrices{},
			LinkDirect: str.String("https://play.example/2"),
		}},
	}}, got)
}

func TestShowsServiceGetSeasonJustwatchLinks(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/seasons/0/watchnow/justwatch_links/pl", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `{"pl":"justwatch.com/pl/serial/rodzina-soprano/sezon-0"}`)
	})

	got, _, err := setup.Client.Shows.GetSeasonJustwatchLinks(context.Background(), str.String("the-sopranos"), test.Ptr(0), str.String("pl"))
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, map[string]string{"pl": "justwatch.com/pl/serial/rodzina-soprano/sezon-0"}, got)
}
