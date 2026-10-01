// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/uri"
)

func TestCalendarsServiceRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	const startDate = "2026-10-01"
	const days = 7
	const (
		movies = `[{"released":"2026-10-02","movie":{"title":"TRON: Legacy","year":2010}}]`
		shows  = `[{"episode":{"season":1,"number":1,"title":"Pilot"},"show":{"title":"Breaking Bad","year":2008}}]`
	)
	cases := []serviceCase{
		{name: "GetDVDReleases", method: http.MethodGet, path: "/calendars/all/dvd/2026-10-01/7", query: "limit=10", body: movies,
			call: func(c *Client) (any, error) {
				r, _, err := c.Calendars.GetDVDReleases(ctx, "all", startDate, days, opts)
				return r, err
			}},
		{name: "GetMovies", method: http.MethodGet, path: "/calendars/my/movies/2026-10-01/7", query: "limit=10", body: movies,
			call: func(c *Client) (any, error) {
				r, _, err := c.Calendars.GetMovies(ctx, "my", startDate, days, opts)
				return r, err
			}},
		{name: "GetSeasonPremieres", method: http.MethodGet, path: "/calendars/all/shows/premieres/2026-10-01/7", query: "limit=10", body: shows,
			call: func(c *Client) (any, error) {
				r, _, err := c.Calendars.GetSeasonPremieres(ctx, "all", startDate, days, opts)
				return r, err
			}},
		{name: "GetShows", method: http.MethodGet, path: "/calendars/my/shows/2026-10-01/7", query: "limit=10", body: shows,
			call: func(c *Client) (any, error) {
				r, _, err := c.Calendars.GetShows(ctx, "my", startDate, days, opts)
				return r, err
			}},
		{name: "GetNewShows", method: http.MethodGet, path: "/calendars/all/shows/new/2026-10-01/7", query: "limit=10", body: shows,
			call: func(c *Client) (any, error) {
				r, _, err := c.Calendars.GetNewShows(ctx, "all", startDate, days, opts)
				return r, err
			}},
		{name: "GetFinales", method: http.MethodGet, path: "/calendars/my/shows/finales/2026-10-01/7", query: "limit=10", body: shows,
			call: func(c *Client) (any, error) {
				r, _, err := c.Calendars.GetFinales(ctx, "my", startDate, days, opts)
				return r, err
			}},
	}
	runServiceCases(t, cases)
}
