// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

// TestUsersPlexHandlers checks every Plex action stops before the API: the locked ones with their notice,
// update_plex_settings and plex_server on invalid input first.
func TestUsersPlexHandlers(t *testing.T) {
	const status = ": Trakt removed this endpoint from its API, the live API no longer serves it; no request was sent"
	const locked = " is locked since 2026-11-18" + status
	const lockedLater = " is locked since 2026-11-27" + status
	tests := []struct {
		name    string
		handler Handler
		options str.Options
		input   string
		wantErr string
	}{
		{name: "settings is locked", handler: UsersPlexSettingsHandler{}, options: str.Options{Action: consts.PlexSettings},
			wantErr: "GET users/settings/plex" + lockedLater},
		{name: "update settings is locked", handler: UsersUpdatePlexSettingsHandler{}, options: str.Options{Action: consts.UpdatePlexSettings},
			input: `{"scrobbler":{"toggles":{"movie":{"watched":true}}},"webhook":{"home_users":"sean"}}`, wantErr: "PUT users/settings/plex" + lockedLater},
		{name: "update settings empty", handler: UsersUpdatePlexSettingsHandler{}, options: str.Options{Action: consts.UpdatePlexSettings},
			input: `{}`, wantErr: "needs at least one of sync, scrobbler, webhook, trigger_sync"},
		{name: "update settings unknown key", handler: UsersUpdatePlexSettingsHandler{}, options: str.Options{Action: consts.UpdatePlexSettings},
			input: `{"scrobbler":{"toggles":{"movie":{"seen":true}}}}`, wantErr: "invalid update_plex_settings JSON"},
		{name: "connect is locked", handler: UsersPlexConnectHandler{}, options: str.Options{Action: consts.PlexConnect},
			wantErr: "POST users/settings/plex/connect" + locked},
		{name: "disconnect is locked", handler: UsersPlexDisconnectHandler{}, options: str.Options{Action: consts.PlexDisconnect},
			wantErr: "DELETE users/settings/plex/connect" + lockedLater},
		{name: "servers is locked", handler: UsersPlexServersHandler{}, options: str.Options{Action: consts.PlexServers},
			wantErr: "GET users/settings/plex/servers" + locked},
		{name: "server is locked", handler: UsersPlexServerHandler{}, options: str.Options{Action: consts.PlexServer, ID: "abc"},
			wantErr: "GET users/settings/plex/servers/{server_id}" + locked},
		{name: "server without id", handler: UsersPlexServerHandler{}, options: str.Options{Action: consts.PlexServer},
			wantErr: "set Plex server id"},
		{name: "sync is locked", handler: UsersPlexSyncHandler{}, options: str.Options{Action: consts.PlexSync, ID: "abc", AllData: true},
			wantErr: "POST users/settings/plex/sync" + locked},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			var calls atomic.Int32
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				calls.Add(1)
				w.WriteHeader(http.StatusNotFound)
			})

			dir := t.TempDir()
			options := tt.options
			options.Output = filepath.Join(dir, "out.json")
			if tt.input != "" {
				options.Items = filepath.Join(dir, "input.json")
				test.AssertNilError(t, os.WriteFile(options.Items, []byte(tt.input), consts.X644))
			}
			err := tt.handler.Handle(&options, s.Client)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
			}
			if got := calls.Load(); got != 0 {
				t.Errorf("%d requests sent, want 0", got)
			}
			if _, statErr := os.Stat(options.Output); !os.IsNotExist(statErr) {
				t.Errorf("unexpected output file: %v", statErr)
			}
		})
	}
}
