// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/uri"
)

func TestPeopleServiceRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	const id = "bryan-cranston"
	const (
		credits = `{"cast":[{"character":"Walter White"}]}`
		person  = `{"name":"Bryan Cranston","birthday":"1956-03-07","birthplace":"San Fernando Valley, California, USA"}`
	)
	cases := []serviceCase{
		{name: "GetListsContainingThisPerson", method: http.MethodGet, path: "/people/bryan-cranston/lists/personal/popular", query: "limit=10", body: `[{"name":"Best actors","privacy":"public","item_count":5}]`,
			call: func(c *Client) (any, error) {
				r, _, err := c.People.GetListsContainingThisPerson(ctx, id, "personal", "popular", opts)
				return r, err
			}},
		{name: "GetMovieCredits", method: http.MethodGet, path: "/people/bryan-cranston/movies", query: "limit=10", body: credits,
			call: func(c *Client) (any, error) { r, _, err := c.People.GetMovieCredits(ctx, id, opts); return r, err }},
		{name: "GetShowCredits", method: http.MethodGet, path: "/people/bryan-cranston/shows", query: "limit=10", body: credits,
			call: func(c *Client) (any, error) { r, _, err := c.People.GetShowCredits(ctx, id, opts); return r, err }},
		{name: "GetSinglePerson", method: http.MethodGet, path: "/people/bryan-cranston", query: "limit=10", body: person,
			call: func(c *Client) (any, error) { r, _, err := c.People.GetSinglePerson(ctx, id, opts); return r, err }},
		{name: "GetRecentlyUpdatedPeople", method: http.MethodGet, path: "/people/updates/2026-10-01", query: "limit=10", body: `[{"person":{"name":"Bryan Cranston"}}]`,
			call: func(c *Client) (any, error) {
				r, _, err := c.People.GetRecentlyUpdatedPeople(ctx, "2026-10-01", opts)
				return r, err
			}},
		{name: "GetRecentlyUpdatedPeopleTraktIDs", method: http.MethodGet, path: "/people/updates/id/2026-10-01", query: "limit=10", body: `[297737,1]`,
			call: func(c *Client) (any, error) {
				r, _, err := c.People.GetRecentlyUpdatedPeopleTraktIDs(ctx, "2026-10-01", opts)
				return r, err
			}},
		{name: "RefreshPersonMetadata", method: http.MethodPost, path: "/people/bryan-cranston/refresh", status: http.StatusCreated,
			call: func(c *Client) (any, error) { _, err := c.People.RefreshPersonMetadata(ctx, id); return nil, err }},
	}
	runServiceCases(t, cases)
}
