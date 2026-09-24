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

func TestMoviesServiceReportMovie(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/movies/tron-legacy-2010/report", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.MovieReport)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.MovieReport{Reason: str.String("runtime"), Message: str.String("wrong runtime")}, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Movies.ReportMovie(context.Background(), str.String("tron-legacy-2010"), &str.MovieReport{Reason: str.String("runtime"), Message: str.String("wrong runtime")})
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}

func TestMoviesServiceReportMovieConflict(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/movies/tron-legacy-2010/report", `{}`)

	_, err := setup.Client.Movies.ReportMovie(context.Background(), str.String("tron-legacy-2010"), &str.MovieReport{Reason: str.String("spam")})
	if got, want := fmt.Sprint(err), fmt.Sprintf(consts.MovieReportPending, "tron-legacy-2010"); got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}

func TestMoviesServiceRefreshMovieJustwatch(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/movies/tron-legacy-2010/refresh/justwatch", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Movies.RefreshMovieJustwatch(context.Background(), str.String("tron-legacy-2010"))
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}
