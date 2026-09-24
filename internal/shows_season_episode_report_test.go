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
)

func TestShowsServiceReportSeason(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/seasons/0/report", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.SeasonReport)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.SeasonReport{Reason: str.String("metadata"), Message: str.String("wrong specials")}, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Shows.ReportSeason(context.Background(), str.String("the-sopranos"), test.Ptr(0), &str.SeasonReport{Reason: str.String("metadata"), Message: str.String("wrong specials")})
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}

func TestShowsServiceReportSeasonConflict(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/shows/the-sopranos/seasons/1/report", `{}`)

	_, err := setup.Client.Shows.ReportSeason(context.Background(), str.String("the-sopranos"), test.Ptr(1), &str.SeasonReport{Reason: str.String("spam")})
	if got, want := fmt.Sprint(err), fmt.Sprintf(consts.SeasonReportPending, 1, "the-sopranos"); got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}

func TestShowsServiceReportEpisode(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/seasons/1/episodes/2/report", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.EpisodeReport)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.EpisodeReport{Reason: str.String("runtime")}, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Shows.ReportEpisode(context.Background(), str.String("the-sopranos"), test.Ptr(1), test.Ptr(2), &str.EpisodeReport{Reason: str.String("runtime")})
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}

func TestShowsServiceReportEpisodeConflict(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/shows/the-sopranos/seasons/1/episodes/2/report", `{}`)

	_, err := setup.Client.Shows.ReportEpisode(context.Background(), str.String("the-sopranos"), test.Ptr(1), test.Ptr(2), &str.EpisodeReport{Reason: str.String("spam")})
	if got, want := fmt.Sprint(err), fmt.Sprintf(consts.EpisodeReportPending, 1, 2, "the-sopranos"); got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}
