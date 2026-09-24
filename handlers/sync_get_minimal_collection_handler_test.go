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

func TestSyncGetMinimalCollectionHandler(t *testing.T) {
	tests := []struct {
		name        string
		typ         string
		availableOn string
		body        string
		wantCall    string
		wantFile    string
		wantErr     string
	}{
		{name: "movies", typ: "movies", body: `{"12601":"2026-09-01T10:20:30.000Z"}`, wantCall: "/sync/collection/minimal/movies?", wantFile: `"12601": "2026-09-01T10:20:30Z"`},
		{name: "episodes on plex", typ: "episodes", availableOn: "plex", body: `{"73640":"2026-09-01T10:20:30.000Z"}`, wantCall: "/sync/collection/minimal/episodes?available_on=plex", wantFile: `"73640"`},
		{name: "shows nested", typ: "shows", body: `{"1390":{"1":{"1":"2026-09-01T10:20:30.000Z"}}}`, wantCall: "/sync/collection/minimal/shows?", wantFile: `"1390": {`},
		{name: "type not in contract", typ: "seasons", wantErr: "not found type for module 'sync'"},
		{name: "available_on not in contract", typ: "movies", availableOn: "kodi", wantErr: "available_on 'kodi' is not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := []string{}
			s.Mux.HandleFunc("/sync/collection/minimal/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				calls = append(calls, r.URL.Path+"?"+r.URL.RawQuery)
				test.SafeFprint(w, tt.body)
			})

			options := str.Options{
				Module:      consts.Sync,
				Action:      consts.GetMinimalCollection,
				Type:        tt.typ,
				AvailableOn: tt.availableOn,
				Output:      filepath.Join(t.TempDir(), "out.json"),
			}
			err := SyncGetMinimalCollectionHandler{}.Handle(&options, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				test.AssertNoDiff(t, []string{}, calls)
				return
			}
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, []string{tt.wantCall}, calls)
			data, readErr := os.ReadFile(options.Output)
			test.AssertNilError(t, readErr)
			if !strings.Contains(string(data), tt.wantFile) {
				t.Errorf("output %s does not contain %s", data, tt.wantFile)
			}
		})
	}
}
