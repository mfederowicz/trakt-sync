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

func TestYounifyServiceGetConnections(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/younify/connections", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `[{"id":"netflix","name":"Netflix","vip":false,"color":"#e50914","images":{"logo":"netflix.png"},`+
			`"connectable":true,"connected":true,"active":true,"profile":"Sean","last_synced_at":"2026-09-25T10:00:00.000Z"},`+
			`{"id":"hulu","name":"Hulu","vip":true,"color":"#1ce783","images":{"logo":null},"connectable":false,"connected":false,"active":false,"profile":null,"last_synced_at":null}]`)
	})

	got, _, err := setup.Client.Younify.GetConnections(context.Background())
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.YounifyConnection{
		{
			ID: str.String("netflix"), Name: str.String("Netflix"), Vip: test.Ptr(false), Color: str.String("#e50914"),
			Images:      &str.WatchNowSourceImages{Logo: str.String("netflix.png")},
			Connectable: test.Ptr(true), Connected: test.Ptr(true), Active: test.Ptr(true), Profile: str.String("Sean"),
			LastSyncedAt: &str.Timestamp{Time: time.Date(2026, 9, 25, 10, 0, 0, 0, time.UTC)},
		},
		{
			ID: str.String("hulu"), Name: str.String("Hulu"), Vip: test.Ptr(true), Color: str.String("#1ce783"),
			Images:      &str.WatchNowSourceImages{},
			Connectable: test.Ptr(false), Connected: test.Ptr(false), Active: test.Ptr(false),
		},
	}, got)
}

func TestYounifyServiceConnect(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/younify/connect", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.YounifyConnect)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.YounifyConnect{ServiceID: str.String("netflix"), ReturnURL: str.String("https://trakt.tv")}, got)
		test.SafeFprint(w, `{"url":"https://auth.example/netflix?sig=abc"}`)
	})

	got, _, err := setup.Client.Younify.Connect(context.Background(), &str.YounifyConnect{ServiceID: str.String("netflix"), ReturnURL: str.String("https://trakt.tv")})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, &str.YounifyConnectResult{URL: str.String("https://auth.example/netflix?sig=abc")}, got)
}

func TestYounifyServiceRefreshAndDisconnect(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		call   func(s *YounifyService) (*str.Response, error)
	}{
		{name: "refresh", method: http.MethodPost, path: "/younify/users/refresh/netflix",
			call: func(s *YounifyService) (*str.Response, error) {
				return s.RefreshService(context.Background(), str.String("netflix"), false)
			}},
		{name: "refresh all data", method: http.MethodPost, path: "/younify/users/refresh/netflix/all_data",
			call: func(s *YounifyService) (*str.Response, error) {
				return s.RefreshService(context.Background(), str.String("netflix"), true)
			}},
		{name: "disconnect", method: http.MethodDelete, path: "/younify/users/services/netflix",
			call: func(s *YounifyService) (*str.Response, error) {
				return s.DisconnectService(context.Background(), str.String("netflix"))
			}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			calls := 0
			setup.Mux.HandleFunc("/younify/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, tt.method)
				if r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				w.WriteHeader(http.StatusNoContent)
			})

			resp, err := tt.call(setup.Client.Younify)
			test.AssertNilError(t, err)
			if resp == nil || resp.StatusCode != http.StatusNoContent {
				t.Errorf("response is %v, want status %d", resp, http.StatusNoContent)
			}
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
		})
	}
}
