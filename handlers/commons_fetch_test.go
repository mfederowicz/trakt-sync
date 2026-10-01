// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/stretchr/testify/assert"
)

// commonCase is one CommonLogic wrapper check: the wrapper must send method+path(+query) built from the
// options and hand back the decoded body.
type commonCase struct {
	name    string
	options str.Options
	method  string
	path    string
	query   string
	status  int
	body    string
	call    func(c *CommonLogic, client *trakt.Client, options *str.Options) (any, error)
}

// runCommonCases runs each case against the mock server and compares the result with body as JSON.
func runCommonCases(t *testing.T, cases []commonCase) {
	t.Helper()
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			testSetup := setup(t)
			defer testSetup.Teardown()

			calls := 0
			testSetup.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, tc.method)
				assert.Equal(t, tc.path, r.URL.Path)
				assert.Equal(t, tc.query, r.URL.RawQuery)
				if tc.status != 0 {
					w.WriteHeader(tc.status)
				}
				if tc.body != "" {
					test.SafeFprint(w, tc.body)
				}
			})

			options := tc.options
			got, err := tc.call(&CommonLogic{}, testSetup.Client, &options)
			test.AssertNilError(t, err)
			assert.Equal(t, 1, calls, "API calls")
			if tc.body == "" {
				return
			}
			encoded, err := json.Marshal(got)
			test.AssertNilError(t, err)
			assert.JSONEq(t, tc.body, string(encoded))
		})
	}
}

func TestCommonFetchSingleObjects(t *testing.T) {
	const (
		comment = `{"id":417,"comment":"Agreed, this show is awesome.","spoiler":false}`
		note    = `{"id":190,"notes":"Watch it in IMAX.","privacy":"private"}`
		item    = `{"type":"movie","movie":{"title":"TRON: Legacy","year":2010}}`
	)
	cases := []commonCase{
		{name: "FetchPerson", options: str.Options{InternalID: "bryan-cranston", ExtendedInfo: "full"},
			method: http.MethodGet, path: "/people/bryan-cranston", query: "extended=full", body: `{"name":"Bryan Cranston"}`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchPerson(client, o)
			}},
		{name: "FetchList", options: str.Options{InternalID: "2143363"},
			method: http.MethodGet, path: "/lists/2143363", body: `{"name":"Star Wars","privacy":"public"}`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) { return c.FetchList(client, o) }},
		{name: "FetchComment", options: str.Options{CommentID: 417},
			method: http.MethodGet, path: "/comments/417", body: comment,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchComment(client, o)
			}},
		{name: "FetchCommentItem", options: str.Options{CommentID: 417, ExtendedInfo: "full"},
			method: http.MethodGet, path: "/comments/417/item", query: "extended=full", body: item,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchCommentItem(client, o)
			}},
		{name: "FetchNotes", options: str.Options{InternalID: "190"},
			method: http.MethodGet, path: "/notes/190", body: note,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchNotes(client, o)
			}},
		{name: "FetchNotesItem", options: str.Options{InternalID: "190"},
			method: http.MethodGet, path: "/notes/190/item", body: item,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchNotesItem(client, o)
			}},
	}
	runCommonCases(t, cases)
}

func TestCommonFetchLists(t *testing.T) {
	const (
		comments = `[{"type":"movie","movie":{"title":"TRON: Legacy","year":2010},"comment":{"id":417,"comment":"Great movie!"}}]`
		history  = `[{"id":101,"type":"movie","movie":{"title":"TRON: Legacy","year":2010}}]`
		ratings  = `[{"rating":9,"type":"movie","movie":{"title":"TRON: Legacy","year":2010}}]`
	)
	paged := str.Options{PerPage: 10, ExtendedInfo: "full"}
	filtered := str.Options{PerPage: 10, CommentType: "reviews", Type: "movies", IncludeReplies: "true"}
	cases := []commonCase{
		{name: "FetchCommentUserLikes", options: str.Options{CommentID: 417, PerPage: 10},
			method: http.MethodGet, path: "/comments/417/likes", query: "limit=10&page=1", body: `[{"user":{"username":"sean","private":false}}]`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchCommentUserLikes(client, o, consts.DefaultPage)
			}},
		{name: "FetchTrendingComments", options: filtered,
			method: http.MethodGet, path: "/comments/trending/reviews/movies", query: "include_replies=true&limit=10&page=1", body: comments,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchTrendingComments(client, o, consts.DefaultPage)
			}},
		{name: "FetchRecentComments", options: filtered,
			method: http.MethodGet, path: "/comments/recent/reviews/movies", query: "include_replies=true&limit=10&page=1", body: comments,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchRecentComments(client, o, consts.DefaultPage)
			}},
		{name: "FetchUpdatedComments", options: filtered,
			method: http.MethodGet, path: "/comments/updates/reviews/movies", query: "include_replies=true&limit=10&page=2", body: comments,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUpdatedComments(client, o, 2)
			}},
		{name: "FetchMovieRecommendations", options: paged,
			method: http.MethodGet, path: "/recommendations/movies", query: "extended=full&limit=10&page=1", body: `[{"title":"Blackfish","year":2013}]`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchMovieRecommendations(client, o, consts.DefaultPage)
			}},
		{name: "FetchShowRecommendations", options: paged,
			method: http.MethodGet, path: "/recommendations/shows", query: "extended=full&limit=10&page=1", body: `[{"title":"Dark","year":2017}]`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchShowRecommendations(client, o, consts.DefaultPage)
			}},
		{name: "FetchHistoryList", options: str.Options{PerPage: 10, Type: "movies", StartDate: "2026-09-01T00:00:00Z", EndDate: "2026-10-01T00:00:00Z"},
			method: http.MethodGet, path: "/sync/history/movies", query: "end_at=2026-10-01T00%3A00%3A00Z&limit=10&page=1&start_at=2026-09-01T00%3A00%3A00Z", body: history,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchHistoryList(client, o, consts.DefaultPage)
			}},
		{name: "FetchHistoryList of one item", options: str.Options{PerPage: 10, Type: "movies", TraktID: 12601},
			method: http.MethodGet, path: "/sync/history/movies/12601", query: "limit=10&page=1", body: history,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchHistoryList(client, o, consts.DefaultPage)
			}},
		{name: "FetchRatings", options: str.Options{PerPage: 10, Type: "movies"},
			method: http.MethodGet, path: "/sync/ratings/movies", query: "limit=10&page=1", body: ratings,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchRatings(client, o, consts.DefaultPage)
			}},
		{name: "FetchRatings filtered by rating", options: str.Options{PerPage: 10, Type: "movies", Rating: str.SliceInt{9, 10}},
			method: http.MethodGet, path: "/sync/ratings/movies/9,10", query: "limit=10&page=1", body: ratings,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchRatings(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersRatings", options: str.Options{PerPage: 10, Type: "shows", UserName: "sean", Rating: str.SliceInt{8}},
			method: http.MethodGet, path: "/users/sean/ratings/shows/8", query: "limit=10&page=1", body: ratings,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersRatings(client, o, consts.DefaultPage)
			}},
	}
	runCommonCases(t, cases)
}

func TestCommonWrites(t *testing.T) {
	const (
		comment  = `{"id":417,"comment":"Agreed, this show is awesome.","spoiler":false}`
		note     = `{"id":190,"notes":"Watch it in IMAX.","privacy":"private"}`
		scrobble = `{"id":3373536619,"action":"start","progress":10.5,"movie":{"title":"TRON: Legacy","year":2010}}`
	)
	cases := []commonCase{
		{name: "UpdateComment", options: str.Options{CommentID: 417},
			method: http.MethodPut, path: "/comments/417", body: comment,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.UpdateComment(client, o, new(str.Comment))
				return r, err
			}},
		{name: "DeleteComment", options: str.Options{CommentID: 417},
			method: http.MethodDelete, path: "/comments/417", status: http.StatusNoContent,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				_, err := c.DeleteComment(client, o)
				return nil, err
			}},
		{name: "Reply", options: str.Options{},
			method: http.MethodPost, path: "/comments/417/replies", status: http.StatusCreated, body: comment,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.Reply(client, Ptr(417), new(str.Comment), o)
				return r, err
			}},
		{name: "Notes", options: str.Options{},
			method: http.MethodPost, path: "/notes", status: http.StatusCreated, body: note,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.Notes(client, new(str.Notes), o)
				return r, err
			}},
		{name: "UpdateNotes", options: str.Options{InternalID: "190"},
			method: http.MethodPut, path: "/notes/190", body: note,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.UpdateNotes(client, o, new(str.Notes))
				return r, err
			}},
		{name: "DeleteNotes", options: str.Options{InternalID: "190"},
			method: http.MethodDelete, path: "/notes/190", status: http.StatusNoContent,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				_, err := c.DeleteNotes(client, o)
				return nil, err
			}},
		{name: "HideMovieRecommendation", options: str.Options{InternalID: "tron-legacy-2010"},
			method: http.MethodDelete, path: "/recommendations/movies/tron-legacy-2010", status: http.StatusNoContent,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				_, err := c.HideMovieRecommendation(client, o)
				return nil, err
			}},
		{name: "HideShowRecommendation", options: str.Options{InternalID: "dark"},
			method: http.MethodDelete, path: "/recommendations/shows/dark", status: http.StatusNoContent,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				_, err := c.HideShowRecommendation(client, o)
				return nil, err
			}},
		{name: "StartScrobble", options: str.Options{},
			method: http.MethodPost, path: "/scrobble/start", status: http.StatusCreated, body: scrobble,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.StartScrobble(client, new(str.Scrobble), o)
				return r, err
			}},
		{name: "PauseScrobble", options: str.Options{},
			method: http.MethodPost, path: "/scrobble/pause", status: http.StatusCreated, body: scrobble,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.PauseScrobble(client, new(str.Scrobble), o)
				return r, err
			}},
		{name: "StopScrobble", options: str.Options{},
			method: http.MethodPost, path: "/scrobble/stop", status: http.StatusCreated, body: scrobble,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.StopScrobble(client, new(str.Scrobble), o)
				return r, err
			}},
	}
	runCommonCases(t, cases)
}

func TestCommonFetchListError(t *testing.T) {
	testSetup := setup(t)
	defer testSetup.Teardown()
	testSetup.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	c := &CommonLogic{}
	options := &str.Options{PerPage: 10, Type: "movies", CommentType: "reviews", CommentID: 417}
	fetches := map[string]func() (any, error){
		"FetchCommentUserLikes": func() (any, error) { return c.FetchCommentUserLikes(testSetup.Client, options, consts.DefaultPage) },
		"FetchTrendingComments": func() (any, error) { return c.FetchTrendingComments(testSetup.Client, options, consts.DefaultPage) },
		"FetchRecentComments":   func() (any, error) { return c.FetchRecentComments(testSetup.Client, options, consts.DefaultPage) },
		"FetchUpdatedComments":  func() (any, error) { return c.FetchUpdatedComments(testSetup.Client, options, consts.DefaultPage) },
		"FetchHistoryList":      func() (any, error) { return c.FetchHistoryList(testSetup.Client, options, consts.DefaultPage) },
		"FetchRatings":          func() (any, error) { return c.FetchRatings(testSetup.Client, options, consts.DefaultPage) },
		"FetchUsersRatings":     func() (any, error) { return c.FetchUsersRatings(testSetup.Client, options, consts.DefaultPage) },
	}
	for name, fetch := range fetches {
		_, err := fetch()
		assert.Error(t, err, name)
	}
}

func TestCreateScrobble(t *testing.T) {
	const (
		movie = `{"title":"TRON: Legacy","year":2010,"ids":{"trakt":12601,"slug":"tron-legacy-2010"}}`
		show  = `{"title":"Dark","year":2017,"ids":{"trakt":119172,"slug":"dark"}}`
	)
	c := &CommonLogic{}

	t.Run("movie with progress", func(t *testing.T) {
		testSetup := setup(t)
		defer testSetup.Teardown()
		testSetup.Mux.HandleFunc("/movies/tron-legacy-2010", func(w http.ResponseWriter, r *http.Request) {
			test.AssertMethod(t, r, http.MethodGet)
			test.SafeFprint(w, movie)
		})

		got, err := c.CreateScrobble(testSetup.Client, &str.Options{Type: consts.Movie, InternalID: "tron-legacy-2010", Progress: 42.5})
		test.AssertNilError(t, err)
		assert.Equal(t, int64(12601), *got.Movie.IDs.Trakt)
		assert.Equal(t, 42.5, *got.Progress)
		assert.Nil(t, got.Episode)
		assert.Nil(t, got.Show)
	})

	t.Run("episode by trakt id", func(t *testing.T) {
		testSetup := setup(t)
		defer testSetup.Teardown()
		testSetup.Mux.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		})

		got, err := c.CreateScrobble(testSetup.Client, &str.Options{Type: consts.Episode, InternalID: "3950"})
		test.AssertNilError(t, err)
		assert.Equal(t, int64(3950), *got.Episode.IDs.Trakt)
		assert.Nil(t, got.Progress, "progress 0 is not sent")
	})

	t.Run("episode with a slug", func(t *testing.T) {
		testSetup := setup(t)
		defer testSetup.Teardown()

		got, err := c.CreateScrobble(testSetup.Client, &str.Options{Type: consts.Episode, InternalID: "pilot"})
		assert.Nil(t, got)
		assert.ErrorContains(t, err, `trakt id must be a positive number, got "pilot"`)
	})

	t.Run("show episode by code and absolute number", func(t *testing.T) {
		testSetup := setup(t)
		defer testSetup.Teardown()
		testSetup.Mux.HandleFunc("/shows/dark", func(w http.ResponseWriter, r *http.Request) {
			test.AssertMethod(t, r, http.MethodGet)
			test.SafeFprint(w, show)
		})

		got, err := c.CreateScrobble(testSetup.Client, &str.Options{Type: consts.ShowEpisode, InternalID: "dark", EpisodeCode: "2x05", EpisodeAbs: 15})
		test.AssertNilError(t, err)
		assert.Equal(t, int64(119172), *got.Show.IDs.Trakt)
		assert.Equal(t, 2, *got.Episode.Season)
		assert.Equal(t, 5, *got.Episode.Number)
		assert.Equal(t, 15, *got.Episode.NumberAbs)
	})

	t.Run("show episode without code", func(t *testing.T) {
		testSetup := setup(t)
		defer testSetup.Teardown()
		testSetup.Mux.HandleFunc("/shows/dark", func(w http.ResponseWriter, _ *http.Request) {
			test.SafeFprint(w, show)
		})

		got, err := c.CreateScrobble(testSetup.Client, &str.Options{Type: consts.ShowEpisode, InternalID: "dark"})
		test.AssertNilError(t, err)
		assert.Nil(t, got.Episode.Season)
		assert.Nil(t, got.Episode.Number)
		assert.Nil(t, got.Episode.NumberAbs)
	})

	t.Run("show episode with a bad code", func(t *testing.T) {
		testSetup := setup(t)
		defer testSetup.Teardown()
		testSetup.Mux.HandleFunc("/shows/dark", func(w http.ResponseWriter, _ *http.Request) {
			test.SafeFprint(w, show)
		})

		got, err := c.CreateScrobble(testSetup.Client, &str.Options{Type: consts.ShowEpisode, InternalID: "dark", EpisodeCode: "s02e05"})
		assert.Nil(t, got)
		assert.ErrorContains(t, err, "invalid format")
	})

	t.Run("show episode of an unknown show", func(t *testing.T) {
		testSetup := setup(t)
		defer testSetup.Teardown()
		testSetup.Mux.HandleFunc("/shows/dark", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		got, err := c.CreateScrobble(testSetup.Client, &str.Options{Type: consts.ShowEpisode, InternalID: "dark", EpisodeCode: "2x05"})
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("unknown type", func(t *testing.T) {
		testSetup := setup(t)
		defer testSetup.Teardown()

		got, err := c.CreateScrobble(testSetup.Client, &str.Options{Type: "season", Progress: 5})
		test.AssertNilError(t, err)
		assert.Nil(t, got.Movie)
		assert.Nil(t, got.Episode)
		assert.Equal(t, 5.0, *got.Progress)
	})
}
