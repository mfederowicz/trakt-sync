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

func TestSeasonsEpisodesReportHandlers(t *testing.T) {
	season := str.Options{Module: consts.Seasons, InternalID: "the-sopranos", Season: 1, Reason: "metadata", Msg: "wrong overview"}
	episode := str.Options{Module: consts.Episodes, InternalID: "the-sopranos", Season: 1, Episode: 2, Reason: "runtime"}
	seasonPath := "/shows/the-sopranos/seasons/1/report"
	episodePath := "/shows/the-sopranos/seasons/1/episodes/2/report"
	seasonByID := str.Options{Module: consts.Seasons, ID: "3950", Reason: "metadata"}
	episodeByID := str.Options{Module: consts.Episodes, ID: "73482", Reason: "runtime"}
	withOptions := func(o str.Options, change func(*str.Options)) str.Options {
		change(&o)
		return o
	}
	tests := []struct {
		name    string
		handler Handler
		options str.Options
		path    string
		status  int
		wantErr string
	}{
		{name: "season report", handler: SeasonsReportHandler{}, options: season, path: seasonPath, status: http.StatusCreated},
		{name: "season report pending", handler: SeasonsReportHandler{}, options: season, path: seasonPath, status: http.StatusConflict, wantErr: "season 1 of show the-sopranos already has a pending report"},
		{name: "season report without id", handler: SeasonsReportHandler{}, options: withOptions(season, func(o *str.Options) { o.InternalID = "" }), wantErr: consts.EmptyShowIDMsg},
		{name: "season report default reason", handler: SeasonsReportHandler{}, options: withOptions(season, func(o *str.Options) { o.Reason = cfg.DefaultConfig().Reason }), wantErr: consts.EmptyReasonMsg},
		{name: "season report invalid reason", handler: SeasonsReportHandler{}, options: withOptions(season, func(o *str.Options) { o.Reason = "boring" }), wantErr: "reason 'boring' is not valid"},
		{name: "season report by trakt id", handler: SeasonsReportHandler{}, options: seasonByID, path: "/seasons/3950/report", status: http.StatusCreated},
		{name: "season report by trakt id pending", handler: SeasonsReportHandler{}, options: seasonByID, path: "/seasons/3950/report", status: http.StatusConflict, wantErr: "season 3950 already has a pending report"},
		{name: "episode report", handler: EpisodesReportHandler{}, options: episode, path: episodePath, status: http.StatusCreated},
		{name: "episode report by trakt id", handler: EpisodesReportHandler{}, options: episodeByID, path: "/episodes/73482/report", status: http.StatusCreated},
		{name: "episode report by trakt id invalid reason", handler: EpisodesReportHandler{}, options: withOptions(episodeByID, func(o *str.Options) { o.Reason = "boring" }), wantErr: "reason 'boring' is not valid"},
		{name: "episode report pending", handler: EpisodesReportHandler{}, options: episode, path: episodePath, status: http.StatusConflict, wantErr: "episode 1x2 of show the-sopranos already has a pending report"},
		{name: "episode report default episode", handler: EpisodesReportHandler{}, options: withOptions(episode, func(o *str.Options) { o.Episode = cfg.DefaultConfig().Episode }), wantErr: consts.EmptyEpisodeMsg},
		{name: "episode report invalid reason", handler: EpisodesReportHandler{}, options: withOptions(episode, func(o *str.Options) { o.Reason = "boring" }), wantErr: "reason 'boring' is not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			handle := func(w http.ResponseWriter, r *http.Request) {
				calls[r.Method+" "+r.URL.Path]++
				w.WriteHeader(tt.status)
				test.SafeFprint(w, `{}`)
			}
			for _, prefix := range []string{"/shows/", "/seasons/", "/episodes/"} {
				s.Mux.HandleFunc(prefix, handle)
			}

			options := tt.options
			options.Action = consts.Report
			err := tt.handler.Handle(&options, s.Client)

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
