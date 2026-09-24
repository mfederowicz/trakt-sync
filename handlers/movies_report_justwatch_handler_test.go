// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
)

func TestMoviesReportAndJustwatchHandlers(t *testing.T) {
	tronOptions := str.Options{InternalID: "tron-legacy-2010"}
	tests := []struct {
		name    string
		handler Handler
		options str.Options
		path    string
		status  int
		wantErr string
	}{
		{name: "report", handler: MoviesReportHandler{}, options: str.Options{InternalID: "tron-legacy-2010", Reason: "spam"}, path: "/movies/tron-legacy-2010/report", status: http.StatusCreated},
		{name: "report without id", handler: MoviesReportHandler{}, options: str.Options{Reason: "spam"}, wantErr: consts.EmptyMovieIDMsg},
		{name: "report without reason", handler: MoviesReportHandler{}, options: tronOptions, wantErr: consts.EmptyReasonMsg},
		{name: "report invalid reason", handler: MoviesReportHandler{}, options: str.Options{InternalID: "tron-legacy-2010", Reason: "boring"}, wantErr: "reason 'boring' is not valid"},
		{name: "refresh justwatch", handler: MoviesRefreshJustwatchHandler{}, options: tronOptions, path: "/movies/tron-legacy-2010/refresh/justwatch", status: http.StatusCreated},
		{name: "refresh justwatch not found", handler: MoviesRefreshJustwatchHandler{}, options: tronOptions, path: "/movies/tron-legacy-2010/refresh/justwatch", status: http.StatusNotFound, wantErr: "not found movie for:tron-legacy-2010"},
		{name: "refresh justwatch without id", handler: MoviesRefreshJustwatchHandler{}, wantErr: consts.EmptyMovieIDMsg},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.Method+" "+r.URL.Path]++
				w.WriteHeader(tt.status)
			})

			options := tt.options
			options.Module = consts.Movies
			options.Action = consts.Report
			err := tt.handler.Handle(&options, s.Client)

			wantCalls := map[string]int{}
			if tt.path != "" {
				wantCalls[http.MethodPost+" "+tt.path] = 1
			}
			if len(calls) != len(wantCalls) || calls[http.MethodPost+" "+tt.path] != wantCalls[http.MethodPost+" "+tt.path] {
				t.Errorf("calls are %v, want %v", calls, wantCalls)
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
