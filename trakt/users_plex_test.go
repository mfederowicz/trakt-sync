// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersServiceGetPlexSettings(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/users/settings/plex", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `{"connection":{"connected":true,"username":null},`+
			`"webhook":{"url":null,"last_event_at":"2026-09-20T18:00:00.000Z","event_count":12,"home_users":"sean"},`+
			`"sync":{"configured":true,"error":false,"server_limit":1,"selection":{"server_ids":["abc"],"library_ids":[{"server_id":"abc","uuid":"u1"}],"user_ids":[]},`+
			`"toggles":{"movie":{"watching":true,"watched":true,"rated":false,"collected":true,"watchlist":false},"show":{"rated":true,"watchlist":true}}},`+
			`"scrobbler":{"toggles":{"episode":{"watching":true,"watched":true,"rated":false,"collected":false}}}}`)
	})

	got, _, err := setup.Client.Users.GetPlexSettings(context.Background())
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, &str.PlexSettings{
		Connection: &str.PlexConnection{Connected: test.Ptr(true)},
		Webhook: &str.PlexWebhook{LastEventAt: &str.Timestamp{Time: time.Date(2026, 9, 20, 18, 0, 0, 0, time.UTC)}, EventCount: test.Ptr(12),
			HomeUsers: str.String("sean")},
		Sync: &str.PlexSync{Configured: test.Ptr(true), Error: test.Ptr(false), ServerLimit: test.Ptr(1),
			Selection: &str.PlexSelection{ServerIDs: []string{"abc"}, LibraryIDs: []*str.PlexLibrarySelection{{ServerID: str.String("abc"), UUID: str.String("u1")}},
				UserIDs: []string{}},
			Toggles: &str.PlexToggleGroups{
				Movie: &str.PlexToggles{Watching: test.Ptr(true), Watched: test.Ptr(true), Rated: test.Ptr(false), Collected: test.Ptr(true), Watchlist: test.Ptr(false)},
				Show:  &str.PlexToggles{Rated: test.Ptr(true), Watchlist: test.Ptr(true)},
			}},
		Scrobbler: &str.PlexScrobbler{Toggles: &str.PlexToggleGroups{
			Episode: &str.PlexToggles{Watching: test.Ptr(true), Watched: test.Ptr(true), Rated: test.Ptr(false), Collected: test.Ptr(false)},
		}},
	}, got)
}

func TestUsersServicePlexRequests(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		status   int
		body     string
		wantBody string
		call     func(u *UsersService) (any, *str.Response, error)
		want     any
	}{
		{name: "update settings", method: http.MethodPut, path: "/users/settings/plex", status: http.StatusNoContent,
			wantBody: `{"sync":{"toggles":{"movie":{"watched":true}}},"trigger_sync":{"watched_all_data":true}}`,
			call: func(u *UsersService) (any, *str.Response, error) {
				resp, err := u.UpdatePlexSettings(context.Background(), &str.PlexSettingsUpdate{
					Sync:        &str.PlexSyncUpdate{Toggles: &str.PlexToggleGroups{Movie: &str.PlexToggles{Watched: test.Ptr(true)}}},
					TriggerSync: &str.PlexTriggerSync{WatchedAllData: test.Ptr(true)},
				})
				return nil, resp, err
			}},
		{name: "disconnect", method: http.MethodDelete, path: "/users/settings/plex/connect", status: http.StatusNoContent,
			call: func(u *UsersService) (any, *str.Response, error) {
				resp, err := u.DisconnectPlex(context.Background())
				return nil, resp, err
			}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, tt.method)
				if tt.wantBody != "" {
					var got, want any
					test.AssertNilError(t, json.NewDecoder(r.Body).Decode(&got))
					test.AssertNilError(t, json.Unmarshal([]byte(tt.wantBody), &want))
					test.AssertNoDiff(t, want, got)
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			got, resp, err := tt.call(setup.Client.Users)
			test.AssertNilError(t, err)
			if resp == nil || resp.StatusCode != tt.status {
				t.Errorf("response is %v, want status %d", resp, tt.status)
			}
			if tt.want != nil {
				test.AssertNoDiff(t, tt.want, got)
			}
		})
	}
}
