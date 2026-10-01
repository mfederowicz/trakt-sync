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
