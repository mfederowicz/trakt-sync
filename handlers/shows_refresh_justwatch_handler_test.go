// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

// The 426 (VIP only) path is covered by cli.HandleVIPResponse tests, where the browser can be stubbed.
func TestShowsRefreshJustwatchHandler(t *testing.T) {
	path := "/shows/the-sopranos/refresh/justwatch"
	tests := []struct {
		name    string
		id      string
		status  int
		wantErr string
	}{
		{name: "refresh justwatch", id: "the-sopranos", status: http.StatusCreated},
		{name: "refresh justwatch not found", id: "the-sopranos", status: http.StatusNotFound, wantErr: "not found show for:the-sopranos"},
		{name: "refresh justwatch without id", wantErr: consts.EmptyShowIDMsg},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/shows/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.Method+" "+r.URL.Path]++
				w.WriteHeader(tt.status)
			})

			err := ShowsRefreshJustwatchHandler{}.Handle(&str.Options{Module: consts.Shows, InternalID: tt.id}, s.Client)

			wantCalls := map[string]int{}
			if tt.id != "" {
				wantCalls[http.MethodPost+" "+path] = 1
			}
			test.AssertNoDiff(t, wantCalls, calls)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			test.AssertNilError(t, err)
		})
	}
}
