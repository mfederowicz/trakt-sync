// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestSyncPlaybackHandler(t *testing.T) {
	tests := []struct {
		name     string
		typ      string
		extended string
		wantPath string
		wantErr  string
	}{
		{name: "all is the untyped route", typ: consts.ActionTypeAll, wantPath: "/sync/playback"},
		{name: "movies", typ: "movies", wantPath: "/sync/playback/movies"},
		{name: "episodes with extended", typ: "episodes", extended: "full", wantPath: "/sync/playback/episodes"},
		{name: "type not in contract", typ: "shows", wantErr: "not found type for module 'sync'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := []string{}
			s.Mux.HandleFunc("/sync/playback", func(w http.ResponseWriter, r *http.Request) {
				calls = append(calls, r.URL.Path+"?"+r.URL.RawQuery)
				test.SafeFprint(w, `[]`)
			})
			s.Mux.HandleFunc("/sync/playback/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				calls = append(calls, r.URL.Path+"?"+r.URL.RawQuery)
				test.SafeFprint(w, `[]`)
			})

			options := str.Options{
				Module:       consts.Sync,
				Action:       consts.Playback,
				Type:         tt.typ,
				ExtendedInfo: tt.extended,
				StartDate:    "2026-09-01T00:00:00.000Z",
				EndDate:      "2026-09-20T00:00:00.000Z",
				Output:       filepath.Join(t.TempDir(), "out.json"),
			}
			err := SyncPlaybackHandler{}.Handle(&options, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				test.AssertNoDiff(t, []string{}, calls)
				return
			}
			test.AssertNilError(t, err)
			if len(calls) != 1 || !strings.HasPrefix(calls[0], tt.wantPath+"?") {
				t.Fatalf("calls are %v, want one call to %s", calls, tt.wantPath)
			}
			if tt.extended != "" && !strings.Contains(calls[0], "extended="+tt.extended) {
				t.Errorf("call %s does not send extended=%s", calls[0], tt.extended)
			}
		})
	}
}
