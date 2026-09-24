// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestSyncProgressHandlers(t *testing.T) {
	tests := []struct {
		name     string
		handler  Handler
		options  str.Options
		wantCall string
		wantErr  string
	}{
		{
			name:     "up next",
			handler:  SyncGetUpNextHandler{},
			options:  str.Options{Action: consts.GetUpNext},
			wantCall: "/sync/progress/up_next?page=1",
		},
		{
			name:     "up next with stats, sort and extended",
			handler:  SyncGetUpNextHandler{},
			options:  str.Options{Action: consts.GetUpNext, IncludeStats: true, LifetimeStats: true, SortBy: "added", SortHow: "desc", ExtendedInfo: "full"},
			wantCall: "/sync/progress/up_next?extended=full&include_stats=true&lifetime_stats=true&page=1&sort_by=added&sort_how=desc",
		},
		{
			name:     "up next ignores watched-only filters",
			handler:  SyncGetUpNextHandler{},
			options:  str.Options{Action: consts.GetUpNext, HideCompleted: true, OnlyRewatching: true},
			wantCall: "/sync/progress/up_next?page=1",
		},
		{
			name:     "watched progress with filters",
			handler:  SyncGetWatchedProgressHandler{},
			options:  str.Options{Action: consts.GetWatchedProgress, HideCompleted: true, OnlyRewatching: true, PerPage: 50},
			wantCall: "/sync/progress/watched?hide_completed=true&limit=50&only_rewatching=true&page=1",
		},
		{
			name:     "watched progress ignores include_stats",
			handler:  SyncGetWatchedProgressHandler{},
			options:  str.Options{Action: consts.GetWatchedProgress, IncludeStats: true},
			wantCall: "/sync/progress/watched?page=1",
		},
		{
			name:     "up next nitro",
			handler:  SyncGetUpNextNitroHandler{},
			options:  str.Options{Action: consts.GetUpNextNitro},
			wantCall: "/sync/progress/up_next_nitro?page=1",
		},
		{
			name:    "up next nitro with every filter",
			handler: SyncGetUpNextNitroHandler{},
			options: str.Options{
				Action: consts.GetUpNextNitro, Intent: "start", WatchNow: "free_all", Genres: "action", Subgenres: "heist",
				Years: "2020", Ratings: "75-100", MediaStartDate: "2026-01-01", MediaEndDate: "2026-12-31", Runtimes: "30-60",
				Countries: "us", Certifications: "tv-14", SortBy: "added", SortHow: "desc", PerPage: 20,
			},
			wantCall: "/sync/progress/up_next_nitro?certifications=tv-14&countries=us&end_date=2026-12-31&genres=action&intent=start&limit=20" +
				"&page=1&ratings=75-100&runtimes=30-60&sort_by=added&sort_how=desc&start_date=2026-01-01&subgenres=heist&watchnow=free_all&years=2020",
		},
		{
			name:     "up next nitro ignores progress-only flags",
			handler:  SyncGetUpNextNitroHandler{},
			options:  str.Options{Action: consts.GetUpNextNitro, IncludeStats: true, HideCompleted: true, ExtendedInfo: "full"},
			wantCall: "/sync/progress/up_next_nitro?page=1",
		},
		{
			name:    "up next nitro intent not in contract",
			handler: SyncGetUpNextNitroHandler{},
			options: str.Options{Action: consts.GetUpNextNitro, Intent: "paused"},
			wantErr: "intent 'paused' is not valid",
		},
		{
			name:    "up next nitro watchnow not in contract",
			handler: SyncGetUpNextNitroHandler{},
			options: str.Options{Action: consts.GetUpNextNitro, WatchNow: "rent"},
			wantErr: "watchnow 'rent' is not valid",
		},
		{
			name:    "watched progress with both hide filters",
			handler: SyncGetWatchedProgressHandler{},
			options: str.Options{Action: consts.GetWatchedProgress, HideCompleted: true, HideNotCompleted: true},
			wantErr: consts.HideBothProgressMsg,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := []string{}
			s.Mux.HandleFunc("/sync/progress/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				calls = append(calls, r.URL.Path+"?"+r.URL.RawQuery)
				test.SafeFprint(w, `[{"show":{"title":"Reacher"},"progress":{"aired":24,"completed":20}}]`)
			})

			options := tt.options
			options.Module = consts.Sync
			options.Output = filepath.Join(t.TempDir(), "out.json")
			err := tt.handler.Handle(&options, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				test.AssertNoDiff(t, []string{}, calls)
				return
			}
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, []string{tt.wantCall}, calls)
			data, readErr := os.ReadFile(options.Output)
			test.AssertNilError(t, readErr)
			if !strings.Contains(string(data), `"title": "Reacher"`) {
				t.Errorf("output %s does not contain the show", data)
			}
		})
	}
}

func TestSyncProgressHandlerPages(t *testing.T) {
	s := setup(t)
	defer s.Teardown()

	pages := []string{}
	s.Mux.HandleFunc("/sync/progress/up_next", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		pages = append(pages, page)
		w.Header().Set(internal.HeaderPaginationPage, page)
		w.Header().Set(internal.HeaderPaginationPageCount, "2")
		test.SafeFprint(w, `[{"show":{"title":"Reacher"}}]`)
	})

	options := str.Options{Module: consts.Sync, Action: consts.GetUpNext, Output: filepath.Join(t.TempDir(), "out.json")}
	test.AssertNilError(t, SyncGetUpNextHandler{}.Handle(&options, s.Client))
	test.AssertNoDiff(t, []string{"1", "2"}, pages)
}
