// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersDataSyncsHandlers(t *testing.T) {
	tests := []struct {
		name       string
		handler    Handler
		options    str.Options
		method     string
		path       string
		status     int
		body       string
		wantErr    string
		wantNoAPI  bool
		wantOutput bool
	}{
		{name: "all syncs", handler: UsersDataSyncsHandler{}, options: str.Options{Action: consts.DataSyncs},
			method: http.MethodGet, path: "/users/syncs", status: http.StatusOK, body: `[{"id":157}]`, wantOutput: true},
		{name: "syncs by type", handler: UsersDataSyncsHandler{}, options: str.Options{Action: consts.DataSyncs, Type: "import"},
			method: http.MethodGet, path: "/users/syncs/import", status: http.StatusOK, body: `[{"id":158}]`, wantOutput: true},
		{name: "syncs invalid type", handler: UsersDataSyncsHandler{}, options: str.Options{Action: consts.DataSyncs, Type: "trakt"},
			wantErr: "set -t to one of [younify plex import]", wantNoAPI: true},
		{name: "syncs not open to API apps", handler: UsersDataSyncsHandler{}, options: str.Options{Action: consts.DataSyncs},
			method: http.MethodGet, path: "/users/syncs", status: http.StatusUnauthorized, body: `{}`, wantErr: "data_syncs: Trakt answered 401"},
		{name: "sync not open to API apps", handler: UsersDataSyncHandler{}, options: str.Options{Action: consts.DataSync, ID: "157"},
			method: http.MethodGet, path: "/users/syncs/157", status: http.StatusUnauthorized, body: `{}`, wantErr: "this route is not open to API apps yet"},
		{name: "skipped not open to API apps", handler: UsersDataSyncItemsHandler{}, options: str.Options{Action: consts.DataSyncSkipped, ID: "157"},
			method: http.MethodGet, path: "/users/syncs/157/skipped", status: http.StatusUnauthorized, body: `{}`, wantErr: "data_sync_skipped: Trakt answered 401"},
		{name: "undo not open to API apps", handler: UsersUndoDataSyncHandler{}, options: str.Options{Action: consts.UndoDataSync, ID: "157"},
			method: http.MethodDelete, path: "/users/syncs/157", status: http.StatusUnauthorized, body: `{}`, wantErr: "undo_data_sync: Trakt answered 401"},
		{name: "syncs empty", handler: UsersDataSyncsHandler{}, options: str.Options{Action: consts.DataSyncs},
			method: http.MethodGet, path: "/users/syncs", status: http.StatusOK, body: `[]`, wantErr: consts.EmptyResult},
		{name: "sync", handler: UsersDataSyncHandler{}, options: str.Options{Action: consts.DataSync, ID: "157"},
			method: http.MethodGet, path: "/users/syncs/157", status: http.StatusOK, body: `{"id":157}`, wantOutput: true},
		{name: "sync of another user", handler: UsersDataSyncHandler{}, options: str.Options{Action: consts.DataSync, ID: "9"},
			method: http.MethodGet, path: "/users/syncs/9", status: http.StatusNotFound, body: `{}`, wantErr: "not found data sync:9"},
		{name: "sync without id", handler: UsersDataSyncHandler{}, options: str.Options{Action: consts.DataSync},
			wantErr: "set data sync id", wantNoAPI: true},
		{name: "paused", handler: UsersDataSyncItemsHandler{}, options: str.Options{Action: consts.DataSyncPaused, ID: "157"},
			method: http.MethodGet, path: "/users/syncs/157/paused", status: http.StatusOK, body: `[{"kind":"history","raw_field":1}]`, wantOutput: true},
		{name: "skipped", handler: UsersDataSyncItemsHandler{}, options: str.Options{Action: consts.DataSyncSkipped, ID: "157"},
			method: http.MethodGet, path: "/users/syncs/157/skipped", status: http.StatusOK, body: `[{"kind":"rating"}]`, wantOutput: true},
		{name: "skipped empty", handler: UsersDataSyncItemsHandler{}, options: str.Options{Action: consts.DataSyncSkipped, ID: "157"},
			method: http.MethodGet, path: "/users/syncs/157/skipped", status: http.StatusOK, body: `[]`, wantErr: "no skipped items in data sync 157"},
		{name: "paused of another user", handler: UsersDataSyncItemsHandler{}, options: str.Options{Action: consts.DataSyncPaused, ID: "9"},
			method: http.MethodGet, path: "/users/syncs/9/paused", status: http.StatusNotFound, body: `{}`, wantErr: "not found data sync:9"},
		{name: "undo", handler: UsersUndoDataSyncHandler{}, options: str.Options{Action: consts.UndoDataSync, ID: "157"},
			method: http.MethodDelete, path: "/users/syncs/157", status: http.StatusNoContent},
		{name: "undo slug id", handler: UsersUndoDataSyncHandler{}, options: str.Options{Action: consts.UndoDataSync, ID: "plex"},
			wantErr: "set data sync id", wantNoAPI: true},
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
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			options := tt.options
			options.Output = filepath.Join(t.TempDir(), "out.json")
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

// TestUsersDataSyncPausedKeepsRawFields checks that the export keeps the passed-through keys of each item.
func TestUsersDataSyncPausedKeepsRawFields(t *testing.T) {
	s := setup(t)
	defer s.Teardown()
	s.Mux.HandleFunc("/users/syncs/157/paused", func(w http.ResponseWriter, _ *http.Request) {
		test.SafeFprint(w, `[{"kind":"history","netflix_title":"Arrival"}]`)
	})

	options := &str.Options{Action: consts.DataSyncPaused, ID: "157", Output: filepath.Join(t.TempDir(), "out.json")}
	test.AssertNilError(t, UsersDataSyncItemsHandler{}.Handle(options, s.Client))
	data, err := os.ReadFile(options.Output)
	test.AssertNilError(t, err)
	if !strings.Contains(string(data), `"netflix_title": "Arrival"`) {
		t.Errorf("export lost the passed-through key:\n%s", data)
	}
}
