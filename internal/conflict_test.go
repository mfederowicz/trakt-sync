// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

// conflictMux registers path to answer with 409 and body.
func conflictMux(t *testing.T, mux *http.ServeMux, path string, body string) {
	t.Helper()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusConflict)
		test.SafeFprint(w, body)
	})
}

func assertConflict(t *testing.T, resp *str.Response, err error) {
	t.Helper()
	var conflict *ConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("error is %v, want *ConflictError", err)
	}
	if resp == nil || resp.StatusCode != http.StatusConflict {
		t.Fatalf("response is %v, want status 409", resp)
	}
}

func TestCheckinConflictKeepsExpiresAt(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/checkin", `{"expires_at":"2026-09-24T20:00:00.000Z"}`)

	result, resp, err := setup.Client.Checkin.CheckintoAnItem(context.Background(), &str.Checkin{})
	assertConflict(t, resp, err)
	if result == nil || result.Expires == nil {
		t.Fatal("expires_at from the 409 body was not decoded")
	}
}

func TestUsersFollowConflictReturnsResponse(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/users/sean/follow", `{}`)

	_, resp, err := setup.Client.Users.Follow(context.Background(), str.String("sean"))
	assertConflict(t, resp, err)
}

func TestUsersBlockConflictReturnsResponse(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/users/sean/block", `{}`)

	resp, err := setup.Client.Users.Block(context.Background(), str.String("sean"))
	assertConflict(t, resp, err)
}

func TestUsersReportConflictKeepsMessage(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/users/sean/report", `{"message":"report already pending"}`)

	result, resp, err := setup.Client.Users.Report(context.Background(), str.String("sean"), &str.UserReport{Reason: str.String("spam")})
	assertConflict(t, resp, err)
	test.AssertNoDiff(t, &str.UserReportResult{Message: str.String("report already pending")}, result)
}

func TestUsersListReportConflictKeepsMessage(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/users/sean/lists/star-wars/report", `{"message":"report already pending"}`)

	result, resp, err := setup.Client.Users.ListReport(context.Background(), str.String("sean"), str.String("star-wars"), &str.ListReport{Reason: str.String("spam")})
	assertConflict(t, resp, err)
	test.AssertNoDiff(t, &str.ListReportResult{Message: str.String("report already pending")}, result)
}
