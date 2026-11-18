// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersPlexHandlers(t *testing.T) {
	const locked = " is locked since 2026-11-18: Trakt removed this endpoint from its API, the live API no longer serves it; no request was sent"
	tests := []struct {
		name       string
		handler    Handler
		options    str.Options
		input      string
		method     string
		path       string
		status     int
		body       string
		wantBody   string
		wantErr    string
		wantNoAPI  bool
		wantOutput bool
	}{
		{name: "settings", handler: UsersPlexSettingsHandler{}, options: str.Options{Action: consts.PlexSettings},
			method: http.MethodGet, path: "/users/settings/plex", status: http.StatusOK, body: `{"connection":{"connected":false}}`, wantOutput: true},
		{name: "settings not open to API apps", handler: UsersPlexSettingsHandler{}, options: str.Options{Action: consts.PlexSettings},
			method: http.MethodGet, path: "/users/settings/plex", status: http.StatusUnauthorized, body: `{}`, wantErr: "plex_settings: Trakt answered 401"},
		{name: "update settings", handler: UsersUpdatePlexSettingsHandler{}, options: str.Options{Action: consts.UpdatePlexSettings},
			input:  `{"scrobbler":{"toggles":{"movie":{"watched":true}}},"webhook":{"home_users":"sean"}}`,
			method: http.MethodPut, path: "/users/settings/plex", status: http.StatusNoContent,
			wantBody: `{"scrobbler":{"toggles":{"movie":{"watched":true}}},"webhook":{"home_users":"sean"}}`},
		{name: "update settings empty", handler: UsersUpdatePlexSettingsHandler{}, options: str.Options{Action: consts.UpdatePlexSettings},
			input: `{}`, wantErr: "needs at least one of sync, scrobbler, webhook, trigger_sync", wantNoAPI: true},
		{name: "update settings unknown key", handler: UsersUpdatePlexSettingsHandler{}, options: str.Options{Action: consts.UpdatePlexSettings},
			input: `{"scrobbler":{"toggles":{"movie":{"seen":true}}}}`, wantErr: "invalid update_plex_settings JSON", wantNoAPI: true},
		{name: "connect is locked", handler: UsersPlexConnectHandler{}, options: str.Options{Action: consts.PlexConnect},
			wantErr: "POST users/settings/plex/connect" + locked, wantNoAPI: true},
		{name: "disconnect", handler: UsersPlexDisconnectHandler{}, options: str.Options{Action: consts.PlexDisconnect},
			method: http.MethodDelete, path: "/users/settings/plex/connect", status: http.StatusNoContent},
		{name: "servers is locked", handler: UsersPlexServersHandler{}, options: str.Options{Action: consts.PlexServers},
			wantErr: "GET users/settings/plex/servers" + locked, wantNoAPI: true},
		{name: "server is locked", handler: UsersPlexServerHandler{}, options: str.Options{Action: consts.PlexServer, ID: "abc"},
			wantErr: "GET users/settings/plex/servers/{server_id}" + locked, wantNoAPI: true},
		{name: "server without id", handler: UsersPlexServerHandler{}, options: str.Options{Action: consts.PlexServer},
			wantErr: "set Plex server id", wantNoAPI: true},
		{name: "sync is locked", handler: UsersPlexSyncHandler{}, options: str.Options{Action: consts.PlexSync, ID: "abc", AllData: true},
			wantErr: "POST users/settings/plex/sync" + locked, wantNoAPI: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, tt.method)
				if r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				if tt.wantBody != "" {
					var got, want any
					test.AssertNilError(t, json.NewDecoder(r.Body).Decode(&got))
					test.AssertNilError(t, json.Unmarshal([]byte(tt.wantBody), &want))
					test.AssertNoDiff(t, want, got)
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			dir := t.TempDir()
			options := tt.options
			options.Output = filepath.Join(dir, "out.json")
			if tt.input != "" {
				options.Items = filepath.Join(dir, "input.json")
				test.AssertNilError(t, os.WriteFile(options.Items, []byte(tt.input), consts.X644))
			}
			err := tt.handler.Handle(&options, s.Client)
			if tt.wantNoAPI && calls != 0 {
				t.Error("API was called with invalid options")
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			test.AssertNilError(t, err)
			_, statErr := os.Stat(options.Output)
			if tt.wantOutput && statErr != nil {
				t.Errorf("output file was not written: %v", statErr)
			}
			if !tt.wantOutput && !os.IsNotExist(statErr) {
				t.Errorf("unexpected output file: %v", statErr)
			}
		})
	}
}
