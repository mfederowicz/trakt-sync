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

func TestShowsServiceReportShow(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/report", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.ShowReport)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.ShowReport{Reason: str.String("metadata"), Message: str.String("wrong overview")}, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Shows.ReportShow(context.Background(), str.String("the-sopranos"), &str.ShowReport{Reason: str.String("metadata"), Message: str.String("wrong overview")})
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}

func TestShowsServiceReportShowConflict(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/shows/the-sopranos/report", `{}`)

	_, err := setup.Client.Shows.ReportShow(context.Background(), str.String("the-sopranos"), &str.ShowReport{Reason: str.String("spam")})
	if got, want := fmt.Sprint(err), fmt.Sprintf(consts.ShowReportPending, "the-sopranos"); got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}

func TestShowsServiceGetShowSentiments(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/sentiments", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `{"good":[{"sentiment":"great acting","comment_ids":[1,2]}],"bad":[{"sentiment":"slow start"}],"comment_count":2}`)
	})

	got, _, err := setup.Client.Shows.GetShowSentiments(context.Background(), str.String("the-sopranos"))
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, &str.Sentiments{
		Good:         []*str.Sentiment{{Sentiment: str.String("great acting"), CommentIDs: &[]int{1, 2}}},
		Bad:          []*str.Sentiment{{Sentiment: str.String("slow start")}},
		CommentCount: test.Ptr(2),
	}, got)
}
