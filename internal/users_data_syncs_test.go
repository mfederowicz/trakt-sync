// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestUsersServiceGetDataSyncs(t *testing.T) {
	row := `{"id":157,"created_at":"2026-09-01T10:00:00.000Z","kind":"younify","source":"netflix","application":null,"undone":false,"undone_at":null,` +
		`"items":{"history":{"movies":3,"episodes":12}},"paused_count":1,"skipped_count":2}`
	want := []*str.DataSync{{
		ID: test.Ptr(int64(157)), CreatedAt: &str.Timestamp{Time: time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)}, Kind: str.String("younify"),
		Source: str.String("netflix"), Undone: test.Ptr(false),
		Items:       &str.DataSyncItems{History: &str.DataSyncCounts{Movies: test.Ptr(3), Episodes: test.Ptr(12)}},
		PausedCount: test.Ptr(1), SkippedCount: test.Ptr(2),
	}}
	for _, tt := range []struct{ name, syncType, path string }{
		{name: "all", path: "/users/syncs"},
		{name: "by type", syncType: "plex", path: "/users/syncs/plex"},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				if got := r.URL.RawQuery; got != "limit=10&page=1" {
					t.Errorf("query is %q, want %q", got, "limit=10&page=1")
				}
				test.SafeFprint(w, "["+row+"]")
			})

			got, _, err := setup.Client.Users.GetDataSyncs(context.Background(), &tt.syncType, &uri.ListOptions{Page: 1, Limit: 10})
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, want, got)
		})
	}
}

func TestUsersServiceDataSyncByID(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		status int
		body   string
		call   func(u *UsersService) (any, *str.Response, error)
		want   any
	}{
		{name: "details", method: http.MethodGet, path: "/users/syncs/157", status: http.StatusOK, body: `{"id":157,"kind":"plex","paused_count":0}`,
			call: func(u *UsersService) (any, *str.Response, error) { return u.GetDataSync(context.Background(), 157) },
			want: &str.DataSync{ID: test.Ptr(int64(157)), Kind: str.String("plex"), PausedCount: test.Ptr(0)}},
		{name: "paused", method: http.MethodGet, path: "/users/syncs/157/paused", status: http.StatusOK, body: `[{"kind":"history","type":"episode","progress":42.5}]`,
			call: func(u *UsersService) (any, *str.Response, error) {
				return u.GetDataSyncItems(context.Background(), 157, str.String("paused"), &uri.ListOptions{})
			},
			want: []*str.SyncItem{{Kind: str.String("history"), Type: str.String("episode"), Progress: test.Ptr(42.5)}}},
		{name: "skipped", method: http.MethodGet, path: "/users/syncs/157/skipped", status: http.StatusOK, body: `[{"kind":"rating","type":null,"trakt_item":null,"rating_value":7}]`,
			call: func(u *UsersService) (any, *str.Response, error) {
				return u.GetDataSyncItems(context.Background(), 157, str.String("skipped"), &uri.ListOptions{})
			},
			want: []*str.SyncItem{{Kind: str.String("rating"), RatingValue: test.Ptr(7)}}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, tt.method)
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			got, _, err := tt.call(setup.Client.Users)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, tt.want, got)
		})
	}
}

func TestUsersServiceUndoDataSync(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/users/syncs/157", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := setup.Client.Users.UndoDataSync(context.Background(), 157)
	test.AssertNilError(t, err)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		t.Errorf("response is %v, want status %d", resp, http.StatusNoContent)
	}
}
