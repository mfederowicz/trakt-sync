// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestSyncServiceGetMinimalCollection(t *testing.T) {
	for _, strType := range []string{"movies", "episodes"} {
		t.Run(strType, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc("/sync/collection/minimal/"+strType, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				test.AssertNoDiff(t, "available_on=digital", r.URL.RawQuery)
				test.SafeFprint(w, `{"12601":"2026-09-01T10:20:30.000Z|plex","12602":"2026-09-02T10:20:30.000Z|"}`)
			})

			got, _, err := setup.Client.Sync.GetMinimalCollection(context.Background(), strType, &uri.ListOptions{AvailableOn: "digital"})
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, str.MinimalCollection{
				"12601": {CollectedAt: &str.Timestamp{Time: time.Date(2026, time.September, 1, 10, 20, 30, 0, time.UTC)}, Source: str.String("plex")},
				"12602": {CollectedAt: &str.Timestamp{Time: time.Date(2026, time.September, 2, 10, 20, 30, 0, time.UTC)}, Source: str.String("")},
			}, got)
		})
	}
}

func TestSyncServiceGetMinimalShowCollection(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/sync/collection/minimal/shows", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.AssertNoDiff(t, "", r.URL.RawQuery)
		test.SafeFprint(w, `{"1390":{"1":{"1":"2026-09-01T10:20:30.000Z","2":"2026-09-02T10:20:30.000Z"}}}`)
	})

	got, _, err := setup.Client.Sync.GetMinimalShowCollection(context.Background(), &uri.ListOptions{})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, str.MinimalShowCollection{
		"1390": {"1": {
			"1": {CollectedAt: &str.Timestamp{Time: time.Date(2026, time.September, 1, 10, 20, 30, 0, time.UTC)}},
			"2": {CollectedAt: &str.Timestamp{Time: time.Date(2026, time.September, 2, 10, 20, 30, 0, time.UTC)}},
		}},
	}, got)
}
