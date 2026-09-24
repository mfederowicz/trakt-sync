// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestShowsServiceRefreshShowJustwatch(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/shows/the-sopranos/refresh/justwatch", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Shows.RefreshShowJustwatch(context.Background(), str.String("the-sopranos"))
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}
