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

func TestShowsWatchNowHandlers(t *testing.T) {
	sopranos := str.Options{InternalID: "the-sopranos", Country: "us"}
	sopranosLinks := str.Options{InternalID: "the-sopranos", Country: "us", Links: "direct", ExtendedInfo: "streaming_ranks"}
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
		{name: "watchnow", handler: ShowsWatchNowHandler{}, options: sopranos, path: "/shows/the-sopranos/watchnow/us", status: http.StatusOK},
		{name: "watchnow links and extended", handler: ShowsWatchNowHandler{}, options: sopranosLinks, path: "/shows/the-sopranos/watchnow/us", query: "extended=streaming_ranks&links=direct", status: http.StatusOK},
		{name: "watchnow not found", handler: ShowsWatchNowHandler{}, options: sopranos, status: http.StatusNotFound, wantErr: "not found show for:the-sopranos"},
		{name: "watchnow limited access", handler: ShowsWatchNowHandler{}, options: sopranos, status: http.StatusForbidden, wantErr: "watchnow is limited access on Trakt"},
		{name: "watchnow without id", handler: ShowsWatchNowHandler{}, options: str.Options{Country: "us"}, wantErr: consts.EmptyShowIDMsg, wantNoAPI: true},
		{name: "watchnow default country", handler: ShowsWatchNowHandler{}, options: str.Options{InternalID: "the-sopranos", Country: cfg.DefaultConfig().ShowsCountry}, wantErr: consts.EmptyCountryMsg, wantNoAPI: true},
		{name: "justwatch links", handler: ShowsJustwatchLinksHandler{}, options: sopranos, path: "/shows/the-sopranos/watchnow/justwatch_links/us", status: http.StatusOK},
		{name: "justwatch links limited access", handler: ShowsJustwatchLinksHandler{}, options: sopranos, status: http.StatusForbidden, wantErr: "justwatch_links is limited access on Trakt"},
		{name: "justwatch links without country", handler: ShowsJustwatchLinksHandler{}, options: str.Options{InternalID: "the-sopranos"}, wantErr: consts.EmptyCountryMsg, wantNoAPI: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/shows/", func(w http.ResponseWriter, r *http.Request) {
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
