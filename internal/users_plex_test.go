// Package internal used for client and services
package internal

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
		{name: "connect", method: http.MethodPost, path: "/users/settings/plex/connect", status: http.StatusOK, body: `{"url":"https://plex.example/auth"}`,
			wantBody: `{"return_url":"trakt://settings"}`,
			call: func(u *UsersService) (any, *str.Response, error) {
				return u.ConnectPlex(context.Background(), &str.PlexConnect{ReturnURL: str.String("trakt://settings")})
			},
			want: &str.PlexConnectResult{URL: str.String("https://plex.example/auth")}},
		{name: "disconnect", method: http.MethodDelete, path: "/users/settings/plex/connect", status: http.StatusNoContent,
			call: func(u *UsersService) (any, *str.Response, error) {
				resp, err := u.DisconnectPlex(context.Background())
				return nil, resp, err
			}},
		{name: "servers", method: http.MethodGet, path: "/users/settings/plex/servers", status: http.StatusOK,
			body: `{"servers":[{"id":"abc","name":"Home","connection_count":2,"connection_timeout":5,"ports":[32400],"owned":true,"url":null}]}`,
			call: func(u *UsersService) (any, *str.Response, error) { return u.GetPlexServers(context.Background()) },
			want: &str.PlexServers{Servers: []*str.PlexServer{{ID: str.String("abc"), Name: str.String("Home"), ConnectionCount: test.Ptr(2),
				ConnectionTimeout: test.Ptr(5), Ports: []int{32400}, Owned: test.Ptr(true)}}}},
		{name: "server accounts", method: http.MethodGet, path: "/users/settings/plex/servers/abc", status: http.StatusOK,
			body: `{"accounts":[{"id":1,"name":"sean"}],"libraries":[{"id":3,"uuid":"u1","type":"movie","title":"Movies","agent":"a","scanner":"s","selected":true,"url":"/lib/3"}]}`,
			call: func(u *UsersService) (any, *str.Response, error) {
				return u.GetPlexServerAccounts(context.Background(), str.String("abc"))
			},
			want: &str.PlexServerAccounts{Accounts: []*str.PlexAccount{{ID: test.Ptr(1), Name: str.String("sean")}},
				Libraries: []*str.PlexLibrary{{ID: test.Ptr(3), UUID: str.String("u1"), Type: str.String("movie"), Title: str.String("Movies"),
					Agent: str.String("a"), Scanner: str.String("s"), Selected: test.Ptr(true), URL: str.String("/lib/3")}}}},
		{name: "sync", method: http.MethodPost, path: "/users/settings/plex/sync", status: http.StatusCreated,
			wantBody: `{"server_id":"abc","all_data":true}`,
			call: func(u *UsersService) (any, *str.Response, error) {
				resp, err := u.SyncPlex(context.Background(), &str.PlexSyncRequest{ServerID: str.String("abc"), AllData: test.Ptr(true)})
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
