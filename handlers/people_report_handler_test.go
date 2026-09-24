// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestPeopleReportHandler(t *testing.T) {
	report := str.Options{ID: "john-wayne", Reason: "metadata", Msg: "wrong birthday"}
	path := "/people/john-wayne/report"
	tests := []struct {
		name    string
		options str.Options
		path    string
		status  int
		wantErr string
	}{
		{name: "report", options: report, path: path, status: http.StatusCreated},
		{name: "report pending", options: report, path: path, status: http.StatusConflict, wantErr: "person john-wayne already has a pending report"},
		{name: "report without id", options: str.Options{Reason: "spam"}, wantErr: consts.EmptyPersonIDMsg},
		{name: "report default reason", options: str.Options{ID: "john-wayne", Reason: cfg.DefaultConfig().Reason}, wantErr: consts.EmptyReasonMsg},
		{name: "report runtime is not a people reason", options: str.Options{ID: "john-wayne", Reason: "runtime"}, wantErr: "reason 'runtime' is not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/people/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.Method+" "+r.URL.Path]++
				w.WriteHeader(tt.status)
				test.SafeFprint(w, `{}`)
			})

			options := tt.options
			options.Module = consts.People
			options.Action = consts.Report
			err := PeopleReportHandler{}.Handle(&options, s.Client)

			wantCalls := map[string]int{}
			if tt.path != "" {
				wantCalls[http.MethodPost+" "+tt.path] = 1
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
