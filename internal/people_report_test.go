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

func TestPeopleServiceReportPerson(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/people/john-wayne/report", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.PersonReport)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.PersonReport{Reason: str.String("metadata"), Message: str.String("wrong birthday")}, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.People.ReportPerson(context.Background(), str.String("john-wayne"), &str.PersonReport{Reason: str.String("metadata"), Message: str.String("wrong birthday")})
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}

func TestPeopleServiceReportPersonConflict(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/people/john-wayne/report", `{}`)

	_, err := setup.Client.People.ReportPerson(context.Background(), str.String("john-wayne"), &str.PersonReport{Reason: str.String("spam")})
	if got, want := fmt.Sprint(err), fmt.Sprintf(consts.PersonReportPending, "john-wayne"); got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}
