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

func TestSeasonsEpisodesWatchNowHandlers(t *testing.T) {
	episode := str.Options{InternalID: "the-sopranos", Season: 1, Episode: 2, Country: "us"}
	episodeLinks := str.Options{InternalID: "the-sopranos", Season: 1, Episode: 2, Country: "us", Links: "direct", ExtendedInfo: "streaming_ranks"}
	season := str.Options{InternalID: "the-sopranos", Season: 1, Country: "pl"}
	episodePath := "/shows/the-sopranos/seasons/1/episodes/2/watchnow/us"
	seasonPath := "/shows/the-sopranos/seasons/1/watchnow/justwatch_links/pl"
	tests := []struct {
		name    string
		handler Handler
		options str.Options
		path    string
		query   string
		status  int
		wantErr string
	}{
		{name: "episode watchnow", handler: EpisodesWatchNowHandler{}, options: episode, path: episodePath, status: http.StatusOK},
		{name: "episode watchnow links and extended", handler: EpisodesWatchNowHandler{}, options: episodeLinks, path: episodePath, query: "extended=streaming_ranks&links=direct", status: http.StatusOK},
		{name: "episode watchnow not found", handler: EpisodesWatchNowHandler{}, options: episode, path: episodePath, status: http.StatusNotFound, wantErr: "not found episode for:the-sopranos"},
		{name: "episode watchnow limited access", handler: EpisodesWatchNowHandler{}, options: episode, path: episodePath, status: http.StatusForbidden, wantErr: "watchnow is limited access on Trakt"},
		{name: "episode watchnow without id", handler: EpisodesWatchNowHandler{}, options: str.Options{Season: 1, Episode: 2, Country: "us"}, wantErr: consts.EmptyShowIDMsg},
		{name: "episode watchnow without country", handler: EpisodesWatchNowHandler{}, options: str.Options{InternalID: "the-sopranos", Season: 1, Episode: 2}, wantErr: consts.EmptyCountryMsg},
		{name: "episode watchnow by trakt id", handler: EpisodesWatchNowHandler{}, options: str.Options{ID: "73482", Country: "us"}, path: "/episodes/73482/watchnow/us", status: http.StatusOK},
		{name: "episode watchnow by trakt id not found", handler: EpisodesWatchNowHandler{}, options: str.Options{ID: "73482", Country: "us"}, path: "/episodes/73482/watchnow/us", status: http.StatusNotFound, wantErr: "not found episode for:73482"},
		{name: "episode watchnow by trakt id without country", handler: EpisodesWatchNowHandler{}, options: str.Options{ID: "73482"}, wantErr: consts.EmptyCountryMsg},
		{name: "season justwatch links", handler: SeasonsJustwatchLinksHandler{}, options: season, path: seasonPath, status: http.StatusOK},
		{name: "season justwatch links not found", handler: SeasonsJustwatchLinksHandler{}, options: season, path: seasonPath, status: http.StatusNotFound, wantErr: "not found season for:the-sopranos"},
		{name: "season justwatch links limited access", handler: SeasonsJustwatchLinksHandler{}, options: season, path: seasonPath, status: http.StatusForbidden, wantErr: "justwatch_links is limited access on Trakt"},
		{name: "season justwatch links without country", handler: SeasonsJustwatchLinksHandler{}, options: str.Options{InternalID: "the-sopranos", Season: 1}, wantErr: consts.EmptyCountryMsg},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			handle := func(w http.ResponseWriter, r *http.Request) {
				calls[r.Method+" "+r.URL.Path]++
				if r.URL.RawQuery != tt.query {
					t.Errorf("query is %q, want %q", r.URL.RawQuery, tt.query)
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, `{}`)
			}
			for _, prefix := range []string{"/shows/", "/episodes/"} {
				s.Mux.HandleFunc(prefix, handle)
			}

			options := tt.options
			options.Output = filepath.Join(t.TempDir(), "out.json")
			err := tt.handler.Handle(&options, s.Client)

			wantCalls := map[string]int{}
			if tt.path != "" {
				wantCalls[http.MethodGet+" "+tt.path] = 1
			}
			test.AssertNoDiff(t, wantCalls, calls)
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
