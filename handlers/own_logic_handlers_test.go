// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// A show without a scheduled episode answers 204: the handler reports it and writes no file.
func TestShowsNextEpisodeNoContent(t *testing.T) {
	s := setup(t)
	defer s.Teardown()
	s.Mux.HandleFunc("/shows/55/next_episode", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusNoContent)
	})

	options := str.Options{InternalID: "55", Output: filepath.Join(t.TempDir(), "out.json")}
	test.AssertNilError(t, ShowsNextEpisodeHandler{}.Handle(&options, s.Client))
	_, err := os.Stat(options.Output)
	assert.True(t, os.IsNotExist(err), "nothing is written without a next episode")
}

// inTempDir runs the test in a temporary working directory, for handlers that write a file with a fixed name.
func inTempDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	previous, err := os.Getwd()
	test.AssertNilError(t, err)
	test.AssertNilError(t, os.Chdir(dir))
	t.Cleanup(func() {
		test.AssertNilError(t, os.Chdir(previous))
	})
	return dir
}

// add_to_history first removes the items from history, then adds them; each step writes its own result file.
func TestSyncAddToHistoryCleansThenAdds(t *testing.T) {
	const items = `[{"type":"movie","watched_at":"2026-10-01T10:00:00.000Z","movie":{"title":"Tron","ids":{"trakt":1}}}]`
	s := setup(t)
	defer s.Teardown()
	requests := []string{}
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.Path)
		body, err := io.ReadAll(r.Body)
		test.AssertNilError(t, err)
		assert.Contains(t, string(body), `"ids":{"trakt":1}`)
		test.SafeFprint(w, `{}`)
	})

	dir := inTempDir(t)
	options := str.Options{Module: "sync", Action: "add_to_history", Type: "movies", Output: filepath.Join(dir, "out.json"), Items: filepath.Join(dir, "items.json")}
	test.AssertNilError(t, os.WriteFile(options.Items, []byte(items), 0o600))

	test.AssertNilError(t, SyncAddToHistoryHandler{}.Handle(&options, s.Client))
	assert.Equal(t, []string{"POST /sync/history/remove", "POST /sync/history"}, requests)
	for _, name := range []string{"sync_remove_from_history_results.json", "out.json"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		test.AssertNilError(t, err)
		assert.True(t, json.Valid(data), "%s is JSON", name)
	}
}

func TestSyncAddToHistoryFailedCleanup(t *testing.T) {
	s := setup(t)
	defer s.Teardown()
	requests := []string{}
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	})

	dir := inTempDir(t)
	options := str.Options{Module: "sync", Action: "add_to_history", Type: "movies", Output: filepath.Join(dir, "out.json"), Items: filepath.Join(dir, "items.json")}
	test.AssertNilError(t, os.WriteFile(options.Items, []byte(`[]`), 0o600))

	assert.Error(t, SyncAddToHistoryHandler{}.Handle(&options, s.Client))
	assert.Equal(t, []string{"POST /sync/history/remove"}, requests, "nothing is added when the cleanup fails")
	entries, err := os.ReadDir(dir)
	test.AssertNilError(t, err)
	assert.Len(t, entries, 1, "only the items file is in the directory")
}

// items without watched_at / rated_at are sent too, they used to be dropped while the items were read.
func TestSyncItemsHandlersItemWithoutDates(t *testing.T) {
	const items = `[{"type":"movie","rating":8,"movie":{"title":"Tron","ids":{"trakt":1}}}]`
	cases := map[string]Handler{
		consts.AddToHistory:      SyncAddToHistoryHandler{},
		consts.RemoveFromHistory: SyncRemoveFromHistoryHandler{},
		consts.AddToRatings:      SyncAddToRatingsHandler{},
		consts.RemoveFromRatings: SyncRemoveFromRatingsHandler{},
	}
	for action, handler := range cases {
		action, handler := action, handler
		t.Run(action, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			requests := consts.ZeroValue
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests++
				body, err := io.ReadAll(r.Body)
				test.AssertNilError(t, err)
				assert.Contains(t, string(body), `"ids":{"trakt":1}`, r.URL.Path)
				test.SafeFprint(w, `{}`)
			})

			dir := inTempDir(t)
			options := str.Options{Module: "sync", Action: action, Type: consts.Movies, Output: filepath.Join(dir, "out.json"), Items: filepath.Join(dir, "items.json")}
			test.AssertNilError(t, os.WriteFile(options.Items, []byte(items), 0o600))

			test.AssertNilError(t, handler.Handle(&options, s.Client))
			assert.NotZero(t, requests)
		})
	}
}

// an unknown -t stops before any request, it used to panic while the items were read.
func TestSyncItemsHandlersUnknownType(t *testing.T) {
	const items = `[{"type":"movie","watched_at":"2026-10-01T10:00:00.000Z","rated_at":"2026-10-01T10:00:00.000Z","rating":8,"movie":{"title":"Tron","ids":{"trakt":1}}}]`
	cases := map[string]Handler{
		consts.AddToHistory:      SyncAddToHistoryHandler{},
		consts.RemoveFromHistory: SyncRemoveFromHistoryHandler{},
		consts.AddToRatings:      SyncAddToRatingsHandler{},
		consts.RemoveFromRatings: SyncRemoveFromRatingsHandler{},
	}
	for action, handler := range cases {
		action, handler := action, handler
		t.Run(action, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			requests := []string{}
			s.Mux.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
				requests = append(requests, r.Method+" "+r.URL.Path)
			})

			dir := t.TempDir()
			options := str.Options{Module: "sync", Action: action, Type: "movie", Output: filepath.Join(dir, "out.json"), Items: filepath.Join(dir, "items.json")}
			test.AssertNilError(t, os.WriteFile(options.Items, []byte(items), 0o600))

			err := handler.Handle(&options, s.Client)
			assert.EqualError(t, err, "type 'movie' is not valid for action '"+action+"', available types:[all movies shows seasons episodes]")
			assert.Empty(t, requests)
		})
	}
}

// an item without a trakt id stops before any request, it used to panic while the items were read.
func TestSyncItemsHandlersItemWithoutTraktID(t *testing.T) {
	const items = `[{"type":"movie","watched_at":"2026-10-01T10:00:00.000Z","rated_at":"2026-10-01T10:00:00.000Z","rating":8,"movie":{"title":"Tron","ids":{"imdb":"tt0084827"}}}]`
	cases := map[string]Handler{
		consts.AddToHistory:      SyncAddToHistoryHandler{},
		consts.RemoveFromHistory: SyncRemoveFromHistoryHandler{},
		consts.AddToRatings:      SyncAddToRatingsHandler{},
		consts.RemoveFromRatings: SyncRemoveFromRatingsHandler{},
	}
	for action, handler := range cases {
		action, handler := action, handler
		t.Run(action, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			requests := []string{}
			s.Mux.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
				requests = append(requests, r.Method+" "+r.URL.Path)
			})

			dir := t.TempDir()
			options := str.Options{Module: "sync", Action: action, Type: consts.Movies, Output: filepath.Join(dir, "out.json"), Items: filepath.Join(dir, "items.json")}
			test.AssertNilError(t, os.WriteFile(options.Items, []byte(items), 0o600))

			err := handler.Handle(&options, s.Client)
			assert.EqualError(t, err, "item at index 0: movie has no trakt id")
			assert.Empty(t, requests)
		})
	}
}

// a reorder item without its id stops before any request, it used to panic while the request was built.
func TestReorderHandlersItemWithoutID(t *testing.T) {
	const (
		items = `[{"id":7,"type":"movie","movie":{"title":"Tron","ids":{"trakt":1}}},{"type":"movie","movie":{"title":"Heat","ids":{"trakt":2}}}]`
		lists = `[{"name":"Favorites","ids":{"trakt":1}},{"name":"Rewatch"}]`
	)
	cases := []struct {
		handler            Handler
		module, action, in string
		want               string
	}{
		{SyncReorderWatchlistHandler{}, "sync", consts.ReorderWatchlist, items, "item at index 1: has no id"},
		{SyncReorderFavoritesHandler{}, "sync", consts.ReorderFavorites, items, "item at index 1: has no id"},
		{UsersReorderListItemsHandler{}, "users", consts.ReorderListItems, items, "item at index 1: has no id"},
		{UsersReorderListsHandler{}, "users", consts.ReorderLists, lists, "list at index 1: has no trakt id"},
		{SyncReorderWatchlistHandler{}, "sync", consts.ReorderWatchlist, `[null]`, "item at index 0: is empty"},
		{UsersReorderListsHandler{}, "users", consts.ReorderLists, `null`, "lists are empty"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.action+" "+tc.want, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			requests := []string{}
			s.Mux.HandleFunc("/", func(_ http.ResponseWriter, r *http.Request) {
				requests = append(requests, r.Method+" "+r.URL.Path)
			})

			dir := t.TempDir()
			options := str.Options{Module: tc.module, Action: tc.action, Type: consts.Movies, UserName: "me", ID: "55", Output: filepath.Join(dir, "out.json"), Items: filepath.Join(dir, "items.json")}
			test.AssertNilError(t, os.WriteFile(options.Items, []byte(tc.in), 0o600))

			err := tc.handler.Handle(&options, s.Client)
			assert.EqualError(t, err, tc.want)
			assert.Empty(t, requests)
		})
	}
}

// An episode title is nullish in the API (upcoming episodes often have none): the episode is still written, it used to panic.
func TestShowsEpisodeHandlersWithoutTitle(t *testing.T) {
	cases := map[string]struct {
		handler Handler
		path    string
	}{
		"last episode": {handler: ShowsLastEpisodeHandler{}, path: "/shows/55/last_episode"},
		"next episode": {handler: ShowsNextEpisodeHandler{}, path: "/shows/55/next_episode"},
	}
	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			s.Mux.HandleFunc(tc.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				test.SafeFprint(w, `{"season":2,"number":5,"title":null,"ids":{"trakt":9}}`)
			})

			options := str.Options{InternalID: "55", Output: filepath.Join(t.TempDir(), "out.json")}
			test.AssertNilError(t, tc.handler.Handle(&options, s.Client))
			written, err := os.ReadFile(options.Output)
			test.AssertNilError(t, err)
			assert.JSONEq(t, `{"season":2,"number":5,"ids":{"trakt":9}}`, string(written))
		})
	}
}
