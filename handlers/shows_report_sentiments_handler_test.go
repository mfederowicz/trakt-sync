// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestShowsReportSentimentsHandlers(t *testing.T) {
	sopranos := str.Options{InternalID: "the-sopranos"}
	report := str.Options{InternalID: "the-sopranos", Reason: "spam", Msg: "fake title"}
	sentimentsPath := "/shows/the-sopranos/sentiments"
	tests := []struct {
		name       string
		handler    Handler
		action     string
		options    str.Options
		method     string
		path       string
		status     int
		wantErr    string
		wantOutput bool
	}{
		{name: "report", handler: ShowsReportHandler{}, action: consts.Report, options: report, method: http.MethodPost, path: "/shows/the-sopranos/report", status: http.StatusCreated},
		{name: "report pending", handler: ShowsReportHandler{}, action: consts.Report, options: report, method: http.MethodPost, path: "/shows/the-sopranos/report", status: http.StatusConflict, wantErr: "show the-sopranos already has a pending report"},
		{name: "report without id", handler: ShowsReportHandler{}, action: consts.Report, options: str.Options{Reason: "spam"}, wantErr: consts.EmptyShowIDMsg},
		{name: "report default reason", handler: ShowsReportHandler{}, action: consts.Report, options: str.Options{InternalID: "the-sopranos", Reason: cfg.DefaultConfig().Reason}, wantErr: consts.EmptyReasonMsg},
		{name: "report invalid reason", handler: ShowsReportHandler{}, action: consts.Report, options: str.Options{InternalID: "the-sopranos", Reason: "boring"}, wantErr: "reason 'boring' is not valid"},
		{name: "sentiments", handler: ShowsSentimentsHandler{}, action: consts.Sentiments, options: sopranos, method: http.MethodGet, path: sentimentsPath, status: http.StatusOK, wantOutput: true},
		{name: "sentiments not found", handler: ShowsSentimentsHandler{}, action: consts.Sentiments, options: sopranos, method: http.MethodGet, path: sentimentsPath, status: http.StatusNotFound, wantErr: "not found show for:the-sopranos"},
		{name: "sentiments without id", handler: ShowsSentimentsHandler{}, action: consts.Sentiments, wantErr: consts.EmptyShowIDMsg},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/shows/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.Method+" "+r.URL.Path]++
				w.WriteHeader(tt.status)
				test.SafeFprint(w, `{"good":[],"bad":[],"comment_count":0}`)
			})

			options := tt.options
			options.Module = consts.Shows
			options.Action = tt.action
			options.Output = filepath.Join(t.TempDir(), "out.json")
			err := tt.handler.Handle(&options, s.Client)

			wantCalls := map[string]int{}
			if tt.path != "" {
				wantCalls[tt.method+" "+tt.path] = 1
			}
			test.AssertNoDiff(t, wantCalls, calls)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			test.AssertNilError(t, err)
			if _, statErr := os.Stat(options.Output); (statErr == nil) != tt.wantOutput {
				t.Errorf("output file written: %v, want %v", statErr == nil, tt.wantOutput)
			}
		})
	}
}
