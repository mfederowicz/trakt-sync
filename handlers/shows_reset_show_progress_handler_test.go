// Package handlers used to handle module actions
package handlers

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// The 426 (VIP only) path is covered by cli.HandleVIPResponse tests, where the browser can be stubbed.
func TestShowsResetShowProgressHandler(t *testing.T) {
	path := "/shows/the-sopranos/progress/watched/reset"
	tests := []struct {
		name     string
		options  str.Options
		method   string
		status   int
		body     string
		wantSent string
		wantOut  string
		wantErr  string
	}{
		{name: "reset answers 201", options: str.Options{InternalID: "the-sopranos", ResetAt: "2026-10-01T20:30:00+02:00"},
			method: http.MethodPost, status: http.StatusCreated, body: `{"reset_at":"2026-10-01T18:30:00.000Z"}`,
			wantSent: `{"reset_at":"2026-10-01T18:30:00Z"}`, wantOut: "result: success, reset show progress:"},
		{name: "reset answers 200", options: str.Options{InternalID: "the-sopranos"},
			method: http.MethodPost, status: http.StatusOK, body: `{"reset_at":"2026-10-01T18:30:00.000Z"}`,
			wantSent: `{}`, wantOut: "result: success, reset show progress:"},
		{name: "undo reset", options: str.Options{InternalID: "the-sopranos", Delete: true},
			method: http.MethodDelete, status: http.StatusNoContent, wantOut: "result: success, undo reset show progress for:the-sopranos"},
		{name: "reset without id", wantErr: consts.EmptyShowIDMsg},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			sent := ""
			s.Mux.HandleFunc("/shows/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.Method+" "+r.URL.Path]++
				body, err := io.ReadAll(r.Body)
				test.AssertNilError(t, err)
				sent = string(body)
				w.WriteHeader(tt.status)
				if tt.body != "" {
					test.SafeFprint(w, tt.body)
				}
			})

			var err error
			out := captureStdout(t, func() {
				err = ShowsResetShowProgressHandler{}.Handle(&tt.options, s.Client)
			})

			wantCalls := map[string]int{}
			if tt.method != "" {
				wantCalls[tt.method+" "+path] = 1
			}
			test.AssertNoDiff(t, wantCalls, calls)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			test.AssertNilError(t, err)
			if tt.wantSent != "" {
				assert.JSONEq(t, tt.wantSent, sent)
			}
			assert.Contains(t, out, tt.wantOut)
		})
	}
}
