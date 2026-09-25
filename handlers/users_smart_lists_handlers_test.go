// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersSmartListsHandlers(t *testing.T) {
	create := `{"name":"Sci-Fi","source":"popular","media_type":"movies","filters":{"genres":["science-fiction"]},"privacy":"public"}`
	tests := []struct {
		name       string
		handler    Handler
		options    str.Options
		input      string
		method     string
		path       string
		status     int
		body       string
		wantErr    string
		wantNoAPI  bool
		wantOutput bool
		vipUser    bool
	}{
		{name: "list", handler: UsersSmartListsHandler{}, options: str.Options{Action: consts.SmartLists},
			method: http.MethodGet, path: "/users/sean/smart-lists", status: http.StatusOK, body: `[{"name":"Sci-Fi"}]`, wantOutput: true},
		{name: "list empty", handler: UsersSmartListsHandler{}, options: str.Options{Action: consts.SmartLists},
			method: http.MethodGet, path: "/users/sean/smart-lists", status: http.StatusOK, body: `[]`, wantErr: consts.EmptyResult},
		{name: "get", handler: UsersSmartListHandler{}, options: str.Options{Action: consts.SmartList, ID: "sci-fi"},
			method: http.MethodGet, path: "/users/sean/smart-lists/sci-fi", status: http.StatusOK, body: `{"name":"Sci-Fi"}`, wantOutput: true},
		{name: "get empty object", handler: UsersSmartListHandler{}, options: str.Options{Action: consts.SmartList, ID: "nope"},
			method: http.MethodGet, path: "/users/sean/smart-lists/nope", status: http.StatusOK, body: `{}`, wantErr: "not found smart list for:nope"},
		{name: "get not found", handler: UsersSmartListHandler{}, options: str.Options{Action: consts.SmartList, ID: "nope"},
			method: http.MethodGet, path: "/users/sean/smart-lists/nope", status: http.StatusNotFound, body: `{}`, wantErr: "not found smart list for:nope"},
		{name: "get without id", handler: UsersSmartListHandler{}, options: str.Options{Action: consts.SmartList},
			wantErr: consts.EmptySmartListIDMsg, wantNoAPI: true},
		{name: "add", handler: UsersAddSmartListHandler{}, options: str.Options{Action: consts.AddSmartList}, input: create,
			method: http.MethodPost, path: "/users/sean/smart-lists", status: http.StatusCreated, body: `{"ids":{"trakt":7,"slug":"sci-fi"}}`, wantOutput: true},
		{name: "add missing source", handler: UsersAddSmartListHandler{}, options: str.Options{Action: consts.AddSmartList},
			input: `{"name":"Sci-Fi","media_type":"movies"}`, wantErr: "needs name, source and media_type", wantNoAPI: true},
		{name: "add invalid media_type", handler: UsersAddSmartListHandler{}, options: str.Options{Action: consts.AddSmartList},
			input: `{"name":"Sci-Fi","source":"popular","media_type":"episodes"}`, wantErr: "media_type 'episodes' is not valid", wantNoAPI: true},
		{name: "add empty source", handler: UsersAddSmartListHandler{}, options: str.Options{Action: consts.AddSmartList},
			input: `{"name":"Sci-Fi","source":"","media_type":"movies"}`, wantErr: "source '' is not valid", wantNoAPI: true},
		{name: "add empty media_type", handler: UsersAddSmartListHandler{}, options: str.Options{Action: consts.AddSmartList},
			input: `{"name":"Sci-Fi","source":"popular","media_type":""}`, wantErr: "media_type '' is not valid", wantNoAPI: true},
		{name: "add unknown key", handler: UsersAddSmartListHandler{}, options: str.Options{Action: consts.AddSmartList},
			input: `{"name":"Sci-Fi","source":"popular","media_type":"movies","genre":["x"]}`, wantErr: "invalid smart list JSON", wantNoAPI: true},
		// X-VIP-User: true only reports the limit; without it a 420 opens the upgrade page in a real browser
		{name: "add over account limit", handler: UsersAddSmartListHandler{}, options: str.Options{Action: consts.AddSmartList}, input: create,
			method: http.MethodPost, path: "/users/sean/smart-lists", status: 420, body: `{}`, vipUser: true, wantErr: "account limit exceeded"},
		{name: "update", handler: UsersUpdateSmartListHandler{}, options: str.Options{Action: consts.UpdateSmartList, ID: "sci-fi"},
			input: `{"privacy":"private"}`, method: http.MethodPut, path: "/users/sean/smart-lists/sci-fi", status: http.StatusOK,
			body: `{"name":"Sci-Fi","privacy":"private"}`, wantOutput: true},
		{name: "update empty body", handler: UsersUpdateSmartListHandler{}, options: str.Options{Action: consts.UpdateSmartList, ID: "sci-fi"},
			input: `{}`, wantErr: "needs at least one of", wantNoAPI: true},
		{name: "update invalid privacy", handler: UsersUpdateSmartListHandler{}, options: str.Options{Action: consts.UpdateSmartList, ID: "sci-fi"},
			input: `{"privacy":"link"}`, wantErr: "privacy 'link' is not valid", wantNoAPI: true},
		{name: "update empty privacy", handler: UsersUpdateSmartListHandler{}, options: str.Options{Action: consts.UpdateSmartList, ID: "sci-fi"},
			input: `{"privacy":""}`, wantErr: "privacy '' is not valid", wantNoAPI: true},
		{name: "update empty name", handler: UsersUpdateSmartListHandler{}, options: str.Options{Action: consts.UpdateSmartList, ID: "sci-fi"},
			input: `{"name":""}`, wantErr: "smart list name must not be empty", wantNoAPI: true},
		{name: "update without id", handler: UsersUpdateSmartListHandler{}, options: str.Options{Action: consts.UpdateSmartList},
			input: `{"privacy":"private"}`, wantErr: consts.EmptySmartListIDMsg, wantNoAPI: true},
		{name: "update not found", handler: UsersUpdateSmartListHandler{}, options: str.Options{Action: consts.UpdateSmartList, ID: "nope"},
			input: `{"privacy":"private"}`, method: http.MethodPut, path: "/users/sean/smart-lists/nope", status: http.StatusNotFound, body: `{}`,
			wantErr: "not found smart list for:nope"},
		{name: "delete", handler: UsersDeleteSmartListHandler{}, options: str.Options{Action: consts.DeleteSmartList, ID: "sci-fi"},
			method: http.MethodDelete, path: "/users/sean/smart-lists/sci-fi", status: http.StatusNoContent},
		{name: "delete not found", handler: UsersDeleteSmartListHandler{}, options: str.Options{Action: consts.DeleteSmartList, ID: "nope"},
			method: http.MethodDelete, path: "/users/sean/smart-lists/nope", status: http.StatusNotFound, body: `{}`, wantErr: "not found smart list for:nope"},
		{name: "delete without id", handler: UsersDeleteSmartListHandler{}, options: str.Options{Action: consts.DeleteSmartList},
			wantErr: consts.EmptySmartListIDMsg, wantNoAPI: true},
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
				if tt.vipUser {
					w.Header().Set(internal.HeaderVIPUser, "true")
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			dir := t.TempDir()
			options := tt.options
			options.UserName = "sean"
			options.Output = filepath.Join(dir, "out.json")
			if tt.input != "" {
				options.Items = filepath.Join(dir, "smart_list.json")
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
				if _, statErr := os.Stat(options.Output); !os.IsNotExist(statErr) {
					t.Errorf("output file written on error: %v", statErr)
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
