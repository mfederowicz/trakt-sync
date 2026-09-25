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

// No case mocks 426: the VIP handling would open the upgrade page in a real browser.
func TestUsersWritesHandlers(t *testing.T) {
	filters := `[{"name":"Sci-fi","url":"https://trakt.tv/movies/trending?genres=science-fiction"}]`
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
	}{
		{name: "update settings", handler: UsersUpdateSettingsHandler{}, options: str.Options{Action: consts.UpdateSettings},
			input:  `{"user":{"location":"Warsaw"},"browsing":{"spoilers":{"episodes":"hide"},"dark_knight":"auto"}}`,
			method: http.MethodPut, path: "/users/settings", status: http.StatusCreated},
		{name: "update settings empty", handler: UsersUpdateSettingsHandler{}, options: str.Options{Action: consts.UpdateSettings},
			input: `{}`, wantErr: "needs a user or browsing block", wantNoAPI: true},
		{name: "update settings invalid spoiler", handler: UsersUpdateSettingsHandler{}, options: str.Options{Action: consts.UpdateSettings},
			input: `{"browsing":{"spoilers":{"shows":"blur"}}}`, wantErr: "browsing.spoilers.shows 'blur' is not valid", wantNoAPI: true},
		{name: "update settings invalid dark_knight", handler: UsersUpdateSettingsHandler{}, options: str.Options{Action: consts.UpdateSettings},
			input: `{"browsing":{"dark_knight":""}}`, wantErr: "browsing.dark_knight '' is not valid", wantNoAPI: true},
		{name: "update settings unknown key", handler: UsersUpdateSettingsHandler{}, options: str.Options{Action: consts.UpdateSettings},
			input: `{"browsing":{"dark_night":"auto"}}`, wantErr: "invalid update_settings JSON", wantNoAPI: true},
		{name: "add saved filters", handler: UsersAddSavedFiltersHandler{}, options: str.Options{Action: consts.AddSavedFilters}, input: filters,
			method: http.MethodPost, path: "/users/saved_filters", status: http.StatusCreated, body: `{"added":[{"id":101}],"skipped":[]}`, wantOutput: true},
		{name: "add saved filters empty", handler: UsersAddSavedFiltersHandler{}, options: str.Options{Action: consts.AddSavedFilters},
			input: `[]`, wantErr: "needs a JSON array", wantNoAPI: true},
		{name: "add saved filters missing url", handler: UsersAddSavedFiltersHandler{}, options: str.Options{Action: consts.AddSavedFilters},
			input: `[{"name":"Sci-fi"}]`, wantErr: "saved filter 1 needs name and url", wantNoAPI: true},
		{name: "delete saved filter", handler: UsersDeleteSavedFilterHandler{}, options: str.Options{Action: consts.DeleteSavedFilter, ID: "101"},
			method: http.MethodDelete, path: "/users/saved_filters/101", status: http.StatusNoContent},
		{name: "delete saved filter not found", handler: UsersDeleteSavedFilterHandler{}, options: str.Options{Action: consts.DeleteSavedFilter, ID: "7"},
			method: http.MethodDelete, path: "/users/saved_filters/7", status: http.StatusNotFound, body: `{}`, wantErr: "not found saved filter:7"},
		{name: "delete saved filter without id", handler: UsersDeleteSavedFilterHandler{}, options: str.Options{Action: consts.DeleteSavedFilter},
			wantErr: "set saved filter id", wantNoAPI: true},
		{name: "delete saved filter slug id", handler: UsersDeleteSavedFilterHandler{}, options: str.Options{Action: consts.DeleteSavedFilter, ID: "sci-fi"},
			wantErr: "set saved filter id", wantNoAPI: true},
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

			dir := t.TempDir()
			options := tt.options
			options.Output = filepath.Join(dir, "out.json")
			if tt.input != "" {
				options.Items = filepath.Join(dir, "input.json")
				test.AssertNilError(t, os.WriteFile(options.Items, []byte(tt.input), consts.X644))
			}
			err := tt.handler.Handle(&options, s.Client)
			if tt.wantNoAPI && calls != 0 {
				t.Error("API was called with invalid input")
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			test.AssertNilError(t, err)
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
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
