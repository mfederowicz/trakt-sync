// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/stretchr/testify/assert"
)

const (
	testPageQuery = "limit=10&page=1"
	testFollowers = `[{"user":{"username":"justin","private":false}}]`
	testItems     = `[{"rank":1,"type":"movie","movie":{"title":"TRON: Legacy","year":2010}}]`
	testLikes     = `[{"type":"list","list":{"name":"Star Wars","privacy":"public"}}]`
	testLists     = `[{"name":"Star Wars","privacy":"public","item_count":5}]`
	testRequest   = `{"id":3,"user":{"username":"justin","private":false}}`
)

func TestCommonUsersFetchLists(t *testing.T) {
	sean := str.Options{UserName: "sean", PerPage: 10}
	seanList := str.Options{UserName: "sean", ID: "star-wars", PerPage: 10}
	cases := []commonCase{
		{name: "FetchWatchlist", options: str.Options{PerPage: 10, Type: "movies", SortBy: "rank", SortHow: "asc"},
			method: http.MethodGet, path: "/sync/watchlist/movies/rank/asc", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchWatchlist(client, o, consts.DefaultPage)
			}},
		{name: "FetchFavorites", options: str.Options{PerPage: 10, Type: "shows", SortBy: "added", SortHow: "desc"},
			method: http.MethodGet, path: "/sync/favorites/shows/added/desc", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchFavorites(client, o, consts.DefaultPage)
			}},
		{name: "FetchPendingFollowingRequests", options: str.Options{ExtendedInfo: "full"},
			method: http.MethodGet, path: "/users/requests/following", query: "extended=full", body: "[" + testRequest + "]",
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchPendingFollowingRequests(client, o)
			}},
		{name: "FetchFollowRequests", options: str.Options{ExtendedInfo: "full"},
			method: http.MethodGet, path: "/users/requests", query: "extended=full", body: "[" + testRequest + "]",
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchFollowRequests(client, o)
			}},
		{name: "FetchUsersHiddenItems", options: str.Options{PerPage: 10, Section: "calendar", Type: "movie"},
			method: http.MethodGet, path: "/users/hidden/calendar", query: "limit=10&page=1&type=movie", body: `[{"type":"movie","movie":{"title":"TRON: Legacy","year":2010}}]`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersHiddenItems(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersLikes", options: str.Options{UserName: "sean", PerPage: 10, Type: "lists"},
			method: http.MethodGet, path: "/users/sean/likes/lists", query: testPageQuery, body: testLikes,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersLikes(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersListLikes", options: seanList,
			method: http.MethodGet, path: "/users/sean/lists/star-wars/likes", query: testPageQuery, body: testFollowers,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersListLikes(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersCollection", options: str.Options{UserName: "sean", PerPage: 10, Type: "movies"},
			method: http.MethodGet, path: "/users/sean/collection/movies", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersCollection(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersComments", options: str.Options{UserName: "sean", PerPage: 10, CommentType: "reviews", Type: "movies", IncludeReplies: "true"},
			method: http.MethodGet, path: "/users/sean/comments/reviews/movies", query: "include_replies=true&limit=10&page=1",
			body: `[{"type":"movie","movie":{"title":"TRON: Legacy","year":2010},"comment":{"id":417,"comment":"Great movie!"}}]`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersComments(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersNotes", options: str.Options{UserName: "sean", PerPage: 10, Type: "movies"},
			method: http.MethodGet, path: "/users/sean/notes/movies", query: testPageQuery, body: `[{"type":"movie","movie":{"title":"TRON: Legacy","year":2010},"note":{"id":190,"notes":"Watch it in IMAX."}}]`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersNotes(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersCollaborations", options: sean,
			method: http.MethodGet, path: "/users/sean/lists/collaborations", query: testPageQuery, body: testLists,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersCollaborations(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersListItems", options: str.Options{UserName: "sean", ID: "star-wars", PerPage: 10, Type: "movies", SortBy: "rank", SortHow: "asc"},
			method: http.MethodGet, path: "/users/sean/lists/star-wars/items/movies/rank/asc", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersListItems(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersListComments", options: str.Options{UserName: "sean", ID: "star-wars", PerPage: 10, Sort: "newest"},
			method: http.MethodGet, path: "/users/sean/lists/star-wars/comments/newest", query: testPageQuery, body: `[{"id":8,"comment":"Can't wait to watch everything on this epic list!","spoiler":false}]`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersListComments(client, o, consts.DefaultPage)
			}},
		{name: "FetchBlockedUsers", options: str.Options{PerPage: 10},
			method: http.MethodGet, path: "/users/blocked", query: testPageQuery, body: testFollowers,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchBlockedUsers(client, o, consts.DefaultPage)
			}},
		{name: "FetchFollowers", options: sean,
			method: http.MethodGet, path: "/users/sean/followers", query: testPageQuery, body: testFollowers,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchFollowers(client, o, consts.DefaultPage)
			}},
		{name: "FetchFollowing", options: sean,
			method: http.MethodGet, path: "/users/sean/following", query: testPageQuery, body: testFollowers,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchFollowing(client, o, consts.DefaultPage)
			}},
		{name: "FetchFriends without a user", options: str.Options{PerPage: 10},
			method: http.MethodGet, path: "/users/me/friends", query: testPageQuery, body: testFollowers,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchFriends(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersHistory", options: str.Options{UserName: "sean", PerPage: 10, Type: "movies", StartDate: "2026-09-01T00:00:00Z", EndDate: "2026-10-01T00:00:00Z"},
			method: http.MethodGet, path: "/users/sean/history/movies", query: "end_at=2026-10-01T00%3A00%3A00Z&limit=10&page=1&start_at=2026-09-01T00%3A00%3A00Z", body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersHistory(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersHistory of one item", options: str.Options{UserName: "sean", PerPage: 10, Type: "movies", ItemID: 12601},
			method: http.MethodGet, path: "/users/sean/history/movies/12601", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersHistory(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersWatchlist", options: str.Options{UserName: "sean", PerPage: 10, Type: "movies", SortBy: "rank", SortHow: "asc"},
			method: http.MethodGet, path: "/users/sean/watchlist/movies/rank/asc", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersWatchlist(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersWatchlist by sort path", options: str.Options{UserName: "sean", PerPage: 10, Type: "all", SortPath: "added"},
			method: http.MethodGet, path: "/users/sean/watchlist/movie,show/added", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersWatchlist(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersWatchlistComments", options: str.Options{UserName: "sean", PerPage: 10, Sort: "likes"},
			method: http.MethodGet, path: "/users/sean/watchlist/comments/likes", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersWatchlistComments(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersFavorites", options: str.Options{UserName: "sean", PerPage: 10, Type: "shows", SortBy: "rank", SortHow: "asc"},
			method: http.MethodGet, path: "/users/sean/favorites/shows/rank/asc", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersFavorites(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersFavorites by sort path", options: str.Options{UserName: "sean", PerPage: 10, Type: "all", SortPath: "title"},
			method: http.MethodGet, path: "/users/sean/favorites/media/title", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersFavorites(client, o, consts.DefaultPage)
			}},
		{name: "FetchUsersFavoritesComments", options: str.Options{UserName: "sean", PerPage: 10, Sort: "newest"},
			method: http.MethodGet, path: "/users/sean/favorites/comments/newest", query: testPageQuery, body: testItems,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.FetchUsersFavoritesComments(client, o, consts.DefaultPage)
			}},
	}
	runCommonCases(t, cases)
}

func TestCommonUsersWrites(t *testing.T) {
	seanList := str.Options{UserName: "sean", ID: "star-wars"}
	cases := []commonCase{
		{name: "ApproveFollowRequest", options: str.Options{FollowerRequest: 3},
			method: http.MethodPost, path: "/users/requests/3", body: testRequest,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.ApproveFollowRequest(client, o)
				return r, err
			}},
		{name: "DenyFollowRequest", options: str.Options{FollowerRequest: 3},
			method: http.MethodDelete, path: "/users/requests/3", body: testRequest,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.DenyFollowRequest(client, o)
				return r, err
			}},
		{name: "UsersAddToHiddenItems", options: str.Options{Section: "calendar"},
			method: http.MethodPost, path: "/users/hidden/calendar", status: http.StatusCreated, body: `{"added":{"movies":1}}`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.UsersAddToHiddenItems(client, o, new(str.HistoryItems))
			}},
		{name: "UsersRemoveHiddenItems", options: str.Options{Section: "calendar"},
			method: http.MethodPost, path: "/users/hidden/calendar/remove", body: `{"deleted":{"movies":1}}`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				return c.UsersRemoveHiddenItems(client, o, new(str.HistoryItems))
			}},
		{name: "UsersAddPersonalList", options: str.Options{UserName: "sean"},
			method: http.MethodPost, path: "/users/sean/lists", status: http.StatusCreated, body: `{"name":"Star Wars","privacy":"public"}`,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				r, _, err := c.UsersAddPersonalList(client, o, new(str.PersonalList))
				return r, err
			}},
		{name: "UsersListLike", options: seanList,
			method: http.MethodPost, path: "/users/sean/lists/star-wars/like", status: http.StatusNoContent,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				_, err := c.UsersListLike(client, o)
				return nil, err
			}},
		{name: "UsersRemoveListLike", options: seanList,
			method: http.MethodDelete, path: "/users/sean/lists/star-wars/like", status: http.StatusNoContent,
			call: func(c *CommonLogic, client *trakt.Client, o *str.Options) (any, error) {
				_, err := c.UsersRemoveListLike(client, o)
				return nil, err
			}},
	}
	runCommonCases(t, cases)
}

func TestCommonUsersSortPathErrors(t *testing.T) {
	testSetup := setup(t)
	defer testSetup.Teardown()
	testSetup.Mux.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
	})
	c := &CommonLogic{}

	_, err := c.FetchUsersWatchlist(testSetup.Client, &str.Options{UserName: "sean", Type: "movies", SortPath: "random"}, consts.DefaultPage)
	assert.ErrorContains(t, err, "sort 'random' is not valid")

	_, err = c.FetchUsersFavorites(testSetup.Client, &str.Options{UserName: "sean", Type: "seasons", SortPath: "rank"}, consts.DefaultPage)
	assert.ErrorContains(t, err, "-sort works with -t all, movies or shows, not 'seasons'")
}

func TestCommonUsersErrors(t *testing.T) {
	c := &CommonLogic{}
	options := &str.Options{UserName: "sean", ID: "star-wars", PerPage: 10, Type: "movies", Section: "calendar", Sort: "newest", FollowerRequest: 3}
	calls := map[string]func(client *trakt.Client) error{
		"FetchWatchlist": func(client *trakt.Client) error {
			_, err := c.FetchWatchlist(client, options, consts.DefaultPage)
			return err
		},
		"FetchFavorites": func(client *trakt.Client) error {
			_, err := c.FetchFavorites(client, options, consts.DefaultPage)
			return err
		},
		"FetchPendingFollowingRequests": func(client *trakt.Client) error {
			_, err := c.FetchPendingFollowingRequests(client, options)
			return err
		},
		"FetchFollowRequests":  func(client *trakt.Client) error { _, err := c.FetchFollowRequests(client, options); return err },
		"ApproveFollowRequest": func(client *trakt.Client) error { _, _, err := c.ApproveFollowRequest(client, options); return err },
		"DenyFollowRequest":    func(client *trakt.Client) error { _, _, err := c.DenyFollowRequest(client, options); return err },
		"FetchUsersHiddenItems": func(client *trakt.Client) error {
			_, err := c.FetchUsersHiddenItems(client, options, consts.DefaultPage)
			return err
		},
		"UsersAddToHiddenItems": func(client *trakt.Client) error {
			_, err := c.UsersAddToHiddenItems(client, options, new(str.HistoryItems))
			return err
		},
		"UsersRemoveHiddenItems": func(client *trakt.Client) error {
			_, err := c.UsersRemoveHiddenItems(client, options, new(str.HistoryItems))
			return err
		},
		"FetchUsersLikes": func(client *trakt.Client) error {
			_, err := c.FetchUsersLikes(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersListLikes": func(client *trakt.Client) error {
			_, err := c.FetchUsersListLikes(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersCollection": func(client *trakt.Client) error {
			_, err := c.FetchUsersCollection(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersComments": func(client *trakt.Client) error {
			_, err := c.FetchUsersComments(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersNotes": func(client *trakt.Client) error {
			_, err := c.FetchUsersNotes(client, options, consts.DefaultPage)
			return err
		},
		"UsersAddPersonalList": func(client *trakt.Client) error {
			_, _, err := c.UsersAddPersonalList(client, options, new(str.PersonalList))
			return err
		},
		"FetchUsersCollaborations": func(client *trakt.Client) error {
			_, err := c.FetchUsersCollaborations(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersListItems": func(client *trakt.Client) error {
			_, err := c.FetchUsersListItems(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersListComments": func(client *trakt.Client) error {
			_, err := c.FetchUsersListComments(client, options, consts.DefaultPage)
			return err
		},
		"FetchBlockedUsers": func(client *trakt.Client) error {
			_, err := c.FetchBlockedUsers(client, options, consts.DefaultPage)
			return err
		},
		"FetchFollowers": func(client *trakt.Client) error {
			_, err := c.FetchFollowers(client, options, consts.DefaultPage)
			return err
		},
		"FetchFollowing": func(client *trakt.Client) error {
			_, err := c.FetchFollowing(client, options, consts.DefaultPage)
			return err
		},
		"FetchFriends": func(client *trakt.Client) error {
			_, err := c.FetchFriends(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersHistory": func(client *trakt.Client) error {
			_, err := c.FetchUsersHistory(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersWatchlist": func(client *trakt.Client) error {
			_, err := c.FetchUsersWatchlist(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersWatchlistComments": func(client *trakt.Client) error {
			_, err := c.FetchUsersWatchlistComments(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersFavorites": func(client *trakt.Client) error {
			_, err := c.FetchUsersFavorites(client, options, consts.DefaultPage)
			return err
		},
		"FetchUsersFavoritesComments": func(client *trakt.Client) error {
			_, err := c.FetchUsersFavoritesComments(client, options, consts.DefaultPage)
			return err
		},
	}

	testSetup := setup(t)
	defer testSetup.Teardown()
	testSetup.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	for name, call := range calls {
		assert.Error(t, call(testSetup.Client), name)
	}
}

func TestFetchUsersListCommentsNotFound(t *testing.T) {
	testSetup := setup(t)
	defer testSetup.Teardown()
	testSetup.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	c := &CommonLogic{}
	_, err := c.FetchUsersListComments(testSetup.Client, &str.Options{UserName: "sean", ID: "star-wars", Sort: "newest"}, consts.DefaultPage)
	assert.EqualError(t, err, "comments not found for:star-wars")
}
