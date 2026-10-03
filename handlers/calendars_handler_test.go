// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

func TestCalendarTarget(t *testing.T) {
	tests := []struct {
		action string
		want   string
	}{
		{action: "my_media", want: "my"},
		{action: "all_media", want: "all"},
		{action: "my_streaming", want: "my"},
		{action: "all_streaming", want: "all"},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			if got := calendarTarget(tt.action); got != tt.want {
				t.Errorf("calendarTarget(%q) = %q, want %q", tt.action, got, tt.want)
			}
		})
	}
}

// The target comes from the action of each call: an all_* call must not change a later my_* call.
func TestCalendarsHandlersTargetPerCall(t *testing.T) {
	tests := []struct {
		name    string
		handler Handler
		all     string
		my      string
		route   string
	}{
		{name: "dvd", handler: CalendarsDvdHandler{}, all: consts.AllDvd, my: consts.MyDvd, route: "dvd"},
		{name: "finales", handler: CalendarsFinalesHandler{}, all: consts.AllFinales, my: consts.MyFinales, route: "shows/finales"},
		{name: "movies", handler: CalendarsMoviesHandler{}, all: consts.AllMovies, my: consts.MyMovies, route: "movies"},
		{name: "new shows", handler: CalendarsNewShowsHandler{}, all: consts.AllNewShows, my: consts.MyNewShows, route: "shows/new"},
		{name: "season premieres", handler: CalendarsSeasonPremieresHandler{}, all: consts.AllSeasonPremieres, my: consts.MySeasonPremieres, route: "shows/premieres"},
		{name: "shows", handler: CalendarsShowsHandler{}, all: consts.AllShows, my: consts.MyShows, route: "shows"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			paths := []string{}
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				paths = append(paths, r.URL.Path)
				test.SafeFprint(w, `[{}]`)
			})

			for _, action := range []string{tt.all, tt.my} {
				options := str.Options{Action: action, StartDate: "2026-10-01", Days: 7, Output: filepath.Join(t.TempDir(), "out.json")}
				test.AssertNilError(t, tt.handler.Handle(&options, s.Client))
			}

			assert.Equal(t, []string{
				"/calendars/all/" + tt.route + "/2026-10-01/7",
				"/calendars/my/" + tt.route + "/2026-10-01/7",
			}, paths)
		})
	}
}

// every calendar action sends the media filter flags; the start_date and end_date filters are left out,
// the calendar's own start date and days are in the path.
func TestCalendarsHandlersMediaFilters(t *testing.T) {
	const filtered = "certifications=pg-13&countries=us&genres=action%2Cdrama&languages=en%2Cpl&ratings=75-100&runtimes=90-150&subgenres=space&watchnow=free&years=2020-2026"
	filters := str.Options{
		WatchNow: "free", Genres: "action,drama", Subgenres: "space", Years: "2020-2026", Ratings: "75-100", Runtimes: "90-150", Countries: "us",
		Certifications: "pg-13", MediaStartDate: "2026-01-01", MediaEndDate: "2026-12-31", Languages: "en,pl",
	}
	handlers := []struct {
		action  string
		handler Handler
		path    string
	}{
		{action: consts.AllShows, handler: CalendarsShowsHandler{}, path: "/calendars/all/shows/2026-10-01/7"},
		{action: consts.AllNewShows, handler: CalendarsNewShowsHandler{}, path: "/calendars/all/shows/new/2026-10-01/7"},
		{action: consts.AllSeasonPremieres, handler: CalendarsSeasonPremieresHandler{}, path: "/calendars/all/shows/premieres/2026-10-01/7"},
		{action: consts.AllFinales, handler: CalendarsFinalesHandler{}, path: "/calendars/all/shows/finales/2026-10-01/7"},
		{action: consts.AllMovies, handler: CalendarsMoviesHandler{}, path: "/calendars/all/movies/2026-10-01/7"},
		{action: consts.AllStreaming, handler: CalendarsStreamingHandler{}, path: "/calendars/all/streaming/2026-10-01/7"},
		{action: consts.AllDvd, handler: CalendarsDvdHandler{}, path: "/calendars/all/dvd/2026-10-01/7"},
		{action: consts.AllMedia, handler: CalendarsMediaHandler{}, path: "/calendars/all/media/2026-10-01/7"},
		{action: consts.HotReleases, handler: CalendarsHotReleasesHandler{}, path: "/calendars/releases/hot/2026-10-01/7"},
		{action: consts.HotPremieres, handler: CalendarsHotPremieresHandler{}, path: "/calendars/releases/hot/premieres/2026-10-01/7"},
		{action: consts.HotFinales, handler: CalendarsHotFinalesHandler{}, path: "/calendars/releases/hot/finales/2026-10-01/7"},
		{action: consts.HotNewShows, handler: CalendarsHotNewShowsHandler{}, path: "/calendars/releases/hot/new/2026-10-01/7"},
	}
	cases := []struct {
		name    string
		options str.Options
		query   string
		wantErr string
	}{
		{name: "all filters", options: filters, query: filtered},
		{name: "no filters"},
		{name: "unknown watchnow", options: str.Options{WatchNow: "cinema"}, wantErr: "watchnow 'cinema' is not valid"},
	}

	for _, h := range handlers {
		for _, tt := range cases {
			h, tt := h, tt
			t.Run(h.action+" "+tt.name, func(t *testing.T) {
				s := setup(t)
				defer s.Teardown()
				requests := []string{}
				s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					test.AssertMethod(t, r, http.MethodGet)
					requests = append(requests, r.URL.Path+"?"+r.URL.RawQuery)
					test.SafeFprint(w, `[{}]`)
				})

				options := tt.options
				options.Action = h.action
				options.StartDate = "2026-10-01"
				options.Days = 7
				options.Output = filepath.Join(t.TempDir(), "out.json")
				err := h.handler.Handle(&options, s.Client)
				if tt.wantErr != "" {
					assert.ErrorContains(t, err, tt.wantErr)
					assert.Empty(t, requests, "no request is sent")
					return
				}
				test.AssertNilError(t, err)
				assert.Equal(t, []string{h.path + "?" + tt.query}, requests)
			})
		}
	}
}
