// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"slices"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestNotesServiceRequests(t *testing.T) {
	ctx := context.Background()
	const id = "190"
	const note = `{"id":190,"notes":"Watch it in IMAX.","privacy":"private","spoiler":false}`
	cases := []serviceCase{
		{name: "DeleteNotes", method: http.MethodDelete, path: "/notes/190", status: http.StatusNoContent,
			call: func(c *Client) (any, error) { _, err := c.Notes.DeleteNotes(ctx, id); return nil, err }},
		{name: "UpdateNotes", method: http.MethodPut, path: "/notes/190", body: note,
			call: func(c *Client) (any, error) { r, _, err := c.Notes.UpdateNotes(ctx, id, new(str.Notes)); return r, err }},
		{name: "GetNotes", method: http.MethodGet, path: "/notes/190", body: note,
			call: func(c *Client) (any, error) { r, _, err := c.Notes.GetNotes(ctx, id); return r, err }},
		{name: "GetNotesItem", method: http.MethodGet, path: "/notes/190/item", body: `{"type":"movie","movie":{"title":"TRON: Legacy","year":2010}}`,
			call: func(c *Client) (any, error) { r, _, err := c.Notes.GetNotesItem(ctx, id); return r, err }},
	}
	runServiceCases(t, cases)
}

func TestScrobbleServiceRequests(t *testing.T) {
	ctx := context.Background()
	const scrobble = `{"id":3373536619,"action":"start","progress":10.5,"movie":{"title":"TRON: Legacy","year":2010}}`
	cases := []serviceCase{
		{name: "StartScrobble", method: http.MethodPost, path: "/scrobble/start", status: http.StatusCreated, body: scrobble,
			call: func(c *Client) (any, error) {
				r, _, err := c.Scrobble.StartScrobble(ctx, new(str.Scrobble))
				return r, err
			}},
		{name: "PauseScrobble", method: http.MethodPost, path: "/scrobble/pause", status: http.StatusCreated, body: scrobble,
			call: func(c *Client) (any, error) {
				r, _, err := c.Scrobble.PauseScrobble(ctx, new(str.Scrobble))
				return r, err
			}},
		{name: "StopScrobble", method: http.MethodPost, path: "/scrobble/stop", status: http.StatusCreated, body: scrobble,
			call: func(c *Client) (any, error) {
				r, _, err := c.Scrobble.StopScrobble(ctx, new(str.Scrobble))
				return r, err
			}},
	}
	runServiceCases(t, cases)
}

func TestRecommendationsServiceRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	const recommendations = `[{"title":"Blackfish","year":2013,"ids":{"trakt":58}}]`
	cases := []serviceCase{
		{name: "HideShowRecommendation", method: http.MethodDelete, path: "/recommendations/shows/922", status: http.StatusNoContent,
			call: func(c *Client) (any, error) {
				_, err := c.Recommendations.HideShowRecommendation(ctx, "922")
				return nil, err
			}},
		{name: "GetMovieRecommendations", method: http.MethodGet, path: "/recommendations/movies", query: "limit=10", body: recommendations,
			call: func(c *Client) (any, error) {
				r, _, err := c.Recommendations.GetMovieRecommendations(ctx, opts)
				return r, err
			}},
		{name: "GetShowRecommendations", method: http.MethodGet, path: "/recommendations/shows", query: "limit=10", body: recommendations,
			call: func(c *Client) (any, error) {
				r, _, err := c.Recommendations.GetShowRecommendations(ctx, opts)
				return r, err
			}},
	}
	runServiceCases(t, cases)
}

func TestOauthServiceRequests(t *testing.T) {
	ctx := context.Background()
	const token = `{"access_token":"access","token_type":"bearer","expires_in":7200,"refresh_token":"refresh","scope":"public","created_at":1487889741}`
	cases := []serviceCase{
		{name: "GenerateNewDeviceCodes", method: http.MethodPost, path: "/oauth/device/code",
			body: `{"device_code":"d9c126a7","user_code":"5055CC52","verification_url":"https://trakt.tv/activate","expires_in":600,"interval":5}`,
			call: func(c *Client) (any, error) {
				r, _, err := c.Oauth.GenerateNewDeviceCodes(ctx, new(str.NewDeviceCode))
				return r, err
			}},
		{name: "PollForAccessToken", method: http.MethodPost, path: "/oauth/device/token", body: token,
			call: func(c *Client) (any, error) {
				r, _, err := c.Oauth.PollForAccessToken(ctx, new(str.NewDeviceToken))
				return r, err
			}},
		{name: "ExchangeRefreshTokenForAccessToken", method: http.MethodPost, path: "/oauth/token", body: token,
			call: func(c *Client) (any, error) {
				r, _, err := c.Oauth.ExchangeRefreshTokenForAccessToken(ctx, new(str.CurrentDeviceToken))
				return r, err
			}},
	}
	runServiceCases(t, cases)
}

func TestSearchServiceRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	const results = `[{"type":"movie","score":26.5,"movie":{"title":"TRON: Legacy","year":2010}}]`
	cases := []serviceCase{
		{name: "GetTextQueryResults", method: http.MethodGet, path: "/search/movie,show", query: "limit=10", body: results,
			call: func(c *Client) (any, error) {
				r, _, err := c.Search.GetTextQueryResults(ctx, "movie,show", opts)
				return r, err
			}},
		{name: "GetIDLookupResults", method: http.MethodGet, path: "/search/imdb/tt1104001", query: "limit=10", body: results,
			call: func(c *Client) (any, error) {
				r, _, err := c.Search.GetIDLookupResults(ctx, "imdb", "tt1104001", opts)
				return r, err
			}},
	}
	runServiceCases(t, cases)
}

func TestSmallServicesRequests(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{Limit: 10}
	cases := []serviceCase{
		{name: "Checkin.DeleteAnyActiveCheckins", method: http.MethodDelete, path: "/checkin", status: http.StatusNoContent,
			call: func(c *Client) (any, error) { _, err := c.Checkin.DeleteAnyActiveCheckins(ctx); return nil, err }},
		{name: "Genres.GetGenres", method: http.MethodGet, path: "/genres/movies", body: `[{"name":"Action","slug":"action"}]`,
			call: func(c *Client) (any, error) { r, _, err := c.Genres.GetGenres(ctx, "movies"); return r, err }},
		{name: "Languages.GetLanguages", method: http.MethodGet, path: "/languages/shows", body: `[{"name":"Polish","code":"pl"}]`,
			call: func(c *Client) (any, error) { r, _, err := c.Languages.GetLanguages(ctx, "shows"); return r, err }},
		{name: "Sync.GetLastActivity", method: http.MethodGet, path: "/sync/last_activities", body: `{"movies":{},"episodes":{}}`,
			call: func(c *Client) (any, error) { r, _, err := c.Sync.GetLastActivity(ctx); return r, err }},
		{name: "Sync.GetWatched", method: http.MethodGet, path: "/sync/watched/movies", query: "limit=10", body: `[{"plays":4,"movie":{"title":"TRON: Legacy","year":2010}}]`,
			call: func(c *Client) (any, error) { r, _, err := c.Sync.GetWatched(ctx, "movies", opts); return r, err }},
		{name: "Sync.GetCollectedSeasons empty collection", method: http.MethodGet, path: "/sync/collection/shows", query: "limit=10", body: `[]`,
			call: func(c *Client) (any, error) { r, _, err := c.Sync.GetCollectedSeasons(ctx, opts); return r, err }},
	}
	runServiceCases(t, cases)
}

func TestSmallServicesOwnErrorMessages(t *testing.T) {
	ctx := context.Background()
	const id = "190"
	const (
		notesNotFound = "notes not found with Id:190"
		notesUser     = "invalid user for notes Id:190"
	)
	cases := []struct {
		name   string
		status int
		want   string
		call   func(c *Client) error
	}{
		{name: "DeleteNotes not found", status: http.StatusNotFound, want: notesNotFound,
			call: func(c *Client) error { _, err := c.Notes.DeleteNotes(ctx, id); return err }},
		{name: "DeleteNotes other user", status: http.StatusUnauthorized, want: notesUser,
			call: func(c *Client) error { _, err := c.Notes.DeleteNotes(ctx, id); return err }},
		{name: "GetNotes not found", status: http.StatusNotFound, want: notesNotFound,
			call: func(c *Client) error { _, _, err := c.Notes.GetNotes(ctx, id); return err }},
		{name: "GetNotes other user", status: http.StatusUnauthorized, want: notesUser,
			call: func(c *Client) error { _, _, err := c.Notes.GetNotes(ctx, id); return err }},
		{name: "GetNotesItem not found", status: http.StatusNotFound, want: notesNotFound,
			call: func(c *Client) error { _, _, err := c.Notes.GetNotesItem(ctx, id); return err }},
		{name: "GetNotesItem other user", status: http.StatusUnauthorized, want: notesUser,
			call: func(c *Client) error { _, _, err := c.Notes.GetNotesItem(ctx, id); return err }},
		{name: "HideShowRecommendation not found", status: http.StatusNotFound, want: "recommendation not found with Id:190",
			call: func(c *Client) error { _, err := c.Recommendations.HideShowRecommendation(ctx, id); return err }},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tc.status)
			})

			err := tc.call(setup.Client)
			if err == nil {
				t.Fatal("error is nil, want an error")
			}
			if err.Error() != tc.want {
				t.Errorf("error is %q, want %q", err.Error(), tc.want)
			}
		})
	}
}

// Every collected season is returned with its own ids (all items used to point at the last season),
// together with the response of the collection request, which callers page on (it used to be nil).
func TestSyncGetCollectedSeasons(t *testing.T) {
	delay := collectedSeasonsDelay
	collectedSeasonsDelay = 0
	t.Cleanup(func() { collectedSeasonsDelay = delay })

	setup := Setup()
	defer setup.Teardown()
	setup.Mux.HandleFunc("/sync/collection/shows", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `[{"show":{"title":"Tron","ids":{"slug":"tron"}},"seasons":[{"number":1},{"number":2}]}]`)
	})
	setup.Mux.HandleFunc("/shows/tron/seasons", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `[{"number":1,"ids":{"trakt":11}},{"number":2,"ids":{"trakt":12}},{"number":3,"ids":{"trakt":13}}]`)
	})

	list, resp, err := setup.Client.Sync.GetCollectedSeasons(context.Background(), &uri.ListOptions{})
	test.AssertNilError(t, err)
	got := []int64{}
	for _, item := range list {
		got = append(got, *item.Season.IDs.Trakt)
	}
	if want := []int64{11, 12}; !slices.Equal(got, want) {
		t.Errorf("collected season ids = %v, want %v", got, want)
	}
	if resp == nil {
		t.Error("GetCollectedSeasons returned a nil response")
	}
}
