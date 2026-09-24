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

func TestMoviesSentimentsHandler(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		status  int
		wantErr string
	}{
		{name: "sentiments", id: "tron-legacy-2010", status: http.StatusOK},
		{name: "not found", id: "tron-legacy-2010", status: http.StatusNotFound, wantErr: "not found movie for:tron-legacy-2010"},
		{name: "without id", wantErr: consts.EmptyMovieIDMsg},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				w.WriteHeader(tt.status)
				test.SafeFprint(w, `{"good":[],"bad":[],"comment_count":0}`)
			})

			output := filepath.Join(t.TempDir(), "out.json")
			err := MoviesSentimentsHandler{}.Handle(&str.Options{InternalID: tt.id, Output: output}, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				if tt.id == "" && calls != 0 {
					t.Error("API was called without an id")
				}
				return
			}
			test.AssertNilError(t, err)
			if _, statErr := os.Stat(output); statErr != nil {
				t.Errorf("output file was not written: %v", statErr)
			}
		})
	}
}
