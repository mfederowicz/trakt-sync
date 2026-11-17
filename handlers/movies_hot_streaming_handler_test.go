// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/trakt"
)

func TestMoviesHotStreamingHandlers(t *testing.T) {
	const locked = " is locked since 2026-11-17: Trakt removed this endpoint from its API, the live API no longer serves it; no request was sent"
	tests := []struct {
		name    string
		handler Handler
		period  string
		wantErr string
		removed bool
	}{
		{name: "hot", handler: MoviesHotHandler{}, wantErr: "GET movies/hot" + locked, removed: true},
		{name: "streaming default period", handler: MoviesStreamingHandler{}, period: cfg.DefaultConfig().MoviesPeriod, wantErr: "GET movies/streaming/{period}" + locked, removed: true},
		{name: "streaming daily", handler: MoviesStreamingHandler{}, period: "daily", wantErr: "GET movies/streaming/{period}" + locked, removed: true},
		{name: "streaming all is not a streaming period", handler: MoviesStreamingHandler{}, period: "all", wantErr: "period 'all' is not valid for streaming"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.URL.Path]++
				w.WriteHeader(http.StatusNotFound)
			})

			options := &str.Options{Module: "movies", Period: tt.period, Output: filepath.Join(t.TempDir(), "out.json")}
			err := tt.handler.Handle(options, s.Client)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
			}
			removed := new(trakt.EndpointRemovedError)
			if got := errors.As(err, &removed); got != tt.removed {
				t.Errorf("error is *trakt.EndpointRemovedError: %v, want %v", got, tt.removed)
			}
			test.AssertNoDiff(t, map[string]int{}, calls)
		})
	}
}
