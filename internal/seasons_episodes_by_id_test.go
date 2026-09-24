// Package internal used for client and services
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestSeasonsServiceReportSeason(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/seasons/3950/report", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.SeasonReport)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.SeasonReport{Reason: str.String("metadata")}, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Seasons.ReportSeason(context.Background(), str.String("3950"), &str.SeasonReport{Reason: str.String("metadata")})
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}

func TestSeasonsServiceReportSeasonConflict(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/seasons/3950/report", `{}`)

	_, err := setup.Client.Seasons.ReportSeason(context.Background(), str.String("3950"), &str.SeasonReport{Reason: str.String("spam")})
	if got, want := fmt.Sprint(err), fmt.Sprintf(consts.SeasonIDReportPending, "3950"); got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}

func TestEpisodesServiceReportEpisode(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/episodes/73482/report", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.EpisodeReport)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.EpisodeReport{Reason: str.String("runtime"), Message: str.String("50 min")}, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Episodes.ReportEpisode(context.Background(), str.String("73482"), &str.EpisodeReport{Reason: str.String("runtime"), Message: str.String("50 min")})
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}

func TestEpisodesServiceReportEpisodeConflict(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/episodes/73482/report", `{}`)

	_, err := setup.Client.Episodes.ReportEpisode(context.Background(), str.String("73482"), &str.EpisodeReport{Reason: str.String("spam")})
	if got, want := fmt.Sprint(err), fmt.Sprintf(consts.EpisodeIDReportPending, "73482"); got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}

func TestEpisodesServiceGetEpisodeWatchNow(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/episodes/73482/watchnow/us", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		if got, want := r.URL.RawQuery, "links=tvos"; got != want {
			t.Errorf("query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `{"us":{"subscription":[{"source":"max","link":"watchnow.trakt.tv/watchnow/1","link_tvos":"https://tv.example/1"}]}}`)
	})

	got, _, err := setup.Client.Episodes.GetEpisodeWatchNow(context.Background(), str.String("73482"), str.String("us"), &uri.ListOptions{Links: "tvos"})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, map[string]*str.WatchNowSources{"us": {
		Subscription: []*str.WatchNowService{{Source: str.String("max"), Link: str.String("watchnow.trakt.tv/watchnow/1"), LinkTvos: str.String("https://tv.example/1")}},
	}}, got)
}
