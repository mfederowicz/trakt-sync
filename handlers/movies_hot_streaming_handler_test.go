// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestMoviesHotStreamingHandlers(t *testing.T) {
	tests := []struct {
		name    string
		handler Handler
		period  string
		path    string
		wantErr string
	}{
		{name: "hot", handler: MoviesHotHandler{}, path: "/movies/hot"},
		{name: "streaming default period", handler: MoviesStreamingHandler{}, period: cfg.DefaultConfig().MoviesPeriod, path: "/movies/streaming/weekly"},
		{name: "streaming daily", handler: MoviesStreamingHandler{}, period: "daily", path: "/movies/streaming/daily"},
		{name: "streaming all is not a streaming period", handler: MoviesStreamingHandler{}, period: "all", wantErr: "period 'all' is not valid for streaming"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.URL.Path]++
				test.SafeFprint(w, `[{"rank":1,"movie":{"title":"Weapons"}}]`)
			})

			options := &str.Options{Module: "movies", Period: tt.period, Output: filepath.Join(t.TempDir(), "out.json")}
			err := tt.handler.Handle(options, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				test.AssertNoDiff(t, map[string]int{}, calls)
				return
			}
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, map[string]int{tt.path: 1}, calls)
		})
	}
}
