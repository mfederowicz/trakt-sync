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

func TestMoviesWatchNowHandlers(t *testing.T) {
	tron := str.Options{InternalID: "tron-legacy-2010", Country: "us"}
	tronLinks := str.Options{InternalID: "tron-legacy-2010", Country: "us", Links: "direct", ExtendedInfo: "streaming_ranks"}
	tests := []struct {
		name      string
		handler   Handler
		options   str.Options
		path      string
		query     string
		status    int
		wantErr   string
		wantNoAPI bool
	}{
		{name: "watchnow", handler: MoviesWatchNowHandler{}, options: tron, path: "/movies/tron-legacy-2010/watchnow/us", status: http.StatusOK},
		{name: "watchnow links and extended", handler: MoviesWatchNowHandler{}, options: tronLinks, path: "/movies/tron-legacy-2010/watchnow/us", query: "extended=streaming_ranks&links=direct", status: http.StatusOK},
		{name: "watchnow not found", handler: MoviesWatchNowHandler{}, options: tron, status: http.StatusNotFound, wantErr: "not found movie for:tron-legacy-2010"},
		{name: "watchnow limited access", handler: MoviesWatchNowHandler{}, options: tron, status: http.StatusForbidden, wantErr: "watchnow is limited access on Trakt"},
		{name: "watchnow without id", handler: MoviesWatchNowHandler{}, options: str.Options{Country: "us"}, wantErr: consts.EmptyMovieIDMsg, wantNoAPI: true},
		{name: "watchnow default country", handler: MoviesWatchNowHandler{}, options: str.Options{InternalID: "tron-legacy-2010", Country: cfg.DefaultConfig().MoviesCountry}, wantErr: consts.EmptyCountryMsg, wantNoAPI: true},
		{name: "justwatch links", handler: MoviesJustwatchLinksHandler{}, options: tron, path: "/movies/tron-legacy-2010/watchnow/justwatch_links/us", status: http.StatusOK},
		{name: "justwatch links limited access", handler: MoviesJustwatchLinksHandler{}, options: tron, status: http.StatusForbidden, wantErr: "justwatch_links is limited access on Trakt"},
		{name: "justwatch links without country", handler: MoviesJustwatchLinksHandler{}, options: str.Options{InternalID: "tron-legacy-2010"}, wantErr: consts.EmptyCountryMsg, wantNoAPI: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				if tt.path != "" && r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				if r.URL.RawQuery != tt.query {
					t.Errorf("query is %q, want %q", r.URL.RawQuery, tt.query)
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, `{}`)
			})

			options := tt.options
			options.Output = filepath.Join(t.TempDir(), "out.json")
			err := tt.handler.Handle(&options, s.Client)
			if tt.wantNoAPI && calls != 0 {
				t.Error("API was called with invalid options")
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			test.AssertNilError(t, err)
			if _, statErr := os.Stat(options.Output); statErr != nil {
				t.Errorf("output file was not written: %v", statErr)
			}
		})
	}
}
