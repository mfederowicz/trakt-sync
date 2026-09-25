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

func TestWatchNowSourcesHandler(t *testing.T) {
	body := `[{"us":[{"source":"netflix","name":"Netflix"}]}]`
	tests := []struct {
		name    string
		country string
		path    string
		status  int
		body    string
		wantErr string
	}{
		{name: "all", path: "/watchnow/sources", status: http.StatusOK, body: body},
		{name: "country", country: "us", path: "/watchnow/sources/us", status: http.StatusOK, body: body},
		{name: "limited access", path: "/watchnow/sources", status: http.StatusForbidden, body: `{}`, wantErr: "watchnow sources is limited access on Trakt"},
		{name: "country not found", country: "xx", path: "/watchnow/sources/xx", status: http.StatusNotFound, body: `{}`, wantErr: "not found sources for:xx"},
		{name: "empty list", path: "/watchnow/sources", status: http.StatusOK, body: `[]`, wantErr: consts.EmptyResult},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			s.Mux.HandleFunc("/watchnow/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				if r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			output := filepath.Join(t.TempDir(), "out.json")
			options := &str.Options{Action: consts.Sources, Country: tt.country, Output: output}
			err := WatchNowSourcesHandler{}.Handle(options, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				if _, statErr := os.Stat(output); !os.IsNotExist(statErr) {
					t.Errorf("output file written on error: %v", statErr)
				}
				return
			}
			test.AssertNilError(t, err)

			data, err := os.ReadFile(output)
			test.AssertNilError(t, err)
			got := []map[string][]*str.WatchNowSource{}
			test.AssertNilError(t, json.Unmarshal(data, &got))
			test.AssertNoDiff(t, []map[string][]*str.WatchNowSource{{"us": {{Source: str.String("netflix"), Name: str.String("Netflix")}}}}, got)
		})
	}
}
