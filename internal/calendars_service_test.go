// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestCalendarsServiceNewRoutes(t *testing.T) {
	startDate := "2026-09-24"
	days := 7
	opts := &uri.ListOptions{Extended: "full"}

	tests := []struct {
		name string
		path string
		body string
		call func(s *CalendarsService) ([]*str.CalendarList, *str.Response, error)
		want []*str.CalendarList
	}{
		{
			name: "my media",
			path: "/calendars/my/media/2026-09-24/7",
			body: `[{"released":"2026-09-25","movie":{"title":"Tron: Ares"}}]`,
			call: func(s *CalendarsService) ([]*str.CalendarList, *str.Response, error) {
				return s.GetMedia(context.Background(), str.String("my"), &startDate, &days, opts)
			},
			want: []*str.CalendarList{{Released: str.String("2026-09-25"), Movie: &str.Movie{Title: str.String("Tron: Ares")}}},
		},
		{
			name: "all streaming",
			path: "/calendars/all/streaming/2026-09-24/7",
			body: `[{"released":"2026-09-26","movie":{"title":"Weapons"}}]`,
			call: func(s *CalendarsService) ([]*str.CalendarList, *str.Response, error) {
				return s.GetStreamingReleases(context.Background(), str.String("all"), &startDate, &days, opts)
			},
			want: []*str.CalendarList{{Released: str.String("2026-09-26"), Movie: &str.Movie{Title: str.String("Weapons")}}},
		},
		{
			name: "hot releases",
			path: "/calendars/releases/hot/2026-09-24/7",
			body: `[{"show":{"title":"Andor"}}]`,
			call: func(s *CalendarsService) ([]*str.CalendarList, *str.Response, error) {
				return s.GetHotReleases(context.Background(), &startDate, &days, opts)
			},
			want: []*str.CalendarList{{Show: &str.Show{Title: str.String("Andor")}}},
		},
		{
			name: "hot premieres",
			path: "/calendars/releases/hot/premieres/2026-09-24/7",
			body: `[{"show":{"title":"Severance"}}]`,
			call: func(s *CalendarsService) ([]*str.CalendarList, *str.Response, error) {
				return s.GetHotPremieres(context.Background(), &startDate, &days, opts)
			},
			want: []*str.CalendarList{{Show: &str.Show{Title: str.String("Severance")}}},
		},
		{
			name: "hot new shows",
			path: "/calendars/releases/hot/new/2026-09-24/7",
			body: `[{"show":{"title":"Pluribus"}}]`,
			call: func(s *CalendarsService) ([]*str.CalendarList, *str.Response, error) {
				return s.GetHotNewShows(context.Background(), &startDate, &days, opts)
			},
			want: []*str.CalendarList{{Show: &str.Show{Title: str.String("Pluribus")}}},
		},
		{
			name: "hot finales",
			path: "/calendars/releases/hot/finales/2026-09-24/7",
			body: `[{"show":{"title":"The Bear"}}]`,
			call: func(s *CalendarsService) ([]*str.CalendarList, *str.Response, error) {
				return s.GetHotFinales(context.Background(), &startDate, &days, opts)
			},
			want: []*str.CalendarList{{Show: &str.Show{Title: str.String("The Bear")}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				if got, want := r.URL.Query().Get("extended"), "full"; got != want {
					t.Errorf("extended query is %q, want %q", got, want)
				}
				test.SafeFprint(w, tt.body)
			})

			got, _, err := tt.call(setup.Client.Calendars)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, tt.want, got)
		})
	}
}
