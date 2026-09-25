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
	badAuth := `{"error_code":"bad_auth","message":"Plex rejected the authorization.","guidance":"Reconnect Plex."}`
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
		{name: "connect default return_url", handler: UsersPlexConnectHandler{}, options: str.Options{Action: consts.PlexConnect},
			method: http.MethodPost, path: "/users/settings/plex/connect", status: http.StatusOK, body: `{"url":"https://plex.example/auth"}`,
			wantBody: `{"return_url":"https://trakt.tv"}`, wantOutput: true},
		{name: "connect bad return_url", handler: UsersPlexConnectHandler{}, options: str.Options{Action: consts.PlexConnect, ReturnURL: "https://example.com"},
			method: http.MethodPost, path: "/users/settings/plex/connect", status: http.StatusBadRequest, body: `{}`, wantErr: "plex_connect rejected return_url https://example.com"},
		{name: "disconnect", handler: UsersPlexDisconnectHandler{}, options: str.Options{Action: consts.PlexDisconnect},
			method: http.MethodDelete, path: "/users/settings/plex/connect", status: http.StatusNoContent},
		{name: "servers", handler: UsersPlexServersHandler{}, options: str.Options{Action: consts.PlexServers},
			method: http.MethodGet, path: "/users/settings/plex/servers", status: http.StatusOK, body: `{"servers":[{"id":"abc"}]}`, wantOutput: true},
		{name: "servers empty", handler: UsersPlexServersHandler{}, options: str.Options{Action: consts.PlexServers},
			method: http.MethodGet, path: "/users/settings/plex/servers", status: http.StatusOK, body: `{"servers":[]}`, wantErr: consts.EmptyResult},
		{name: "servers plex bad_auth", handler: UsersPlexServersHandler{}, options: str.Options{Action: consts.PlexServers},
			method: http.MethodGet, path: "/users/settings/plex/servers", status: http.StatusUnauthorized, body: badAuth,
			wantErr: "plex_servers: Plex bad_auth: Plex rejected the authorization. Reconnect Plex."},
		{name: "servers plex timeout", handler: UsersPlexServersHandler{}, options: str.Options{Action: consts.PlexServers},
			method: http.MethodGet, path: "/users/settings/plex/servers", status: http.StatusGatewayTimeout,
			body: `{"error_code":"plex_timeout","message":"Plex did not answer.","guidance":"Try again."}`, wantErr: "plex_servers: Plex request failed with status 504"},
		{name: "server", handler: UsersPlexServerHandler{}, options: str.Options{Action: consts.PlexServer, ID: "abc"},
			method: http.MethodGet, path: "/users/settings/plex/servers/abc", status: http.StatusOK, body: `{"accounts":[],"libraries":[]}`, wantOutput: true},
		{name: "server not found", handler: UsersPlexServerHandler{}, options: str.Options{Action: consts.PlexServer, ID: "nope"},
			method: http.MethodGet, path: "/users/settings/plex/servers/nope", status: http.StatusNotFound,
			body: `{"error_code":"plex_not_found","message":"x","guidance":"y"}`, wantErr: "not found Plex server:nope"},
		{name: "server without id", handler: UsersPlexServerHandler{}, options: str.Options{Action: consts.PlexServer},
			wantErr: "set Plex server id", wantNoAPI: true},
		{name: "sync all servers", handler: UsersPlexSyncHandler{}, options: str.Options{Action: consts.PlexSync},
			method: http.MethodPost, path: "/users/settings/plex/sync", status: http.StatusCreated, wantBody: `{}`},
		{name: "sync one server full", handler: UsersPlexSyncHandler{}, options: str.Options{Action: consts.PlexSync, ID: "abc", AllData: true},
			method: http.MethodPost, path: "/users/settings/plex/sync", status: http.StatusCreated, wantBody: `{"server_id":"abc","all_data":true}`},
		{name: "sync without server", handler: UsersPlexSyncHandler{}, options: str.Options{Action: consts.PlexSync},
			method: http.MethodPost, path: "/users/settings/plex/sync", status: http.StatusUnprocessableEntity, body: `{}`, wantErr: "no Plex server to sync"},
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
