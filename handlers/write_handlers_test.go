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

// writeHandler is a handler that sends one write request, built from flags or from the -items file.
type writeHandler struct {
	name    string
	handler Handler
	options str.Options
	items   string
	sends   string
	method  string
	path    string
	status  int // 200 when not set
	writes  bool
}

func writeHandlers() []writeHandler {
	const (
		// id is the list item id, which the reorder actions send as the new rank
		movies = `[{"id":7,"type":"movie","rating":8,"rated_at":"2026-10-01T10:00:00.000Z","watched_at":"2026-10-01T10:00:00.000Z","movie":{"title":"Tron","ids":{"trakt":1}}}]`
		movie  = `"ids":{"trakt":1}`
		note   = `{"notes":"rewatch"}`
		lists  = `[{"name":"Favorites","ids":{"trakt":1}}]`
	)
	sync := func(action string) str.Options {
		return str.Options{Module: "sync", Action: action, Type: "movies"}
	}
	users := func(action string) str.Options {
		return str.Options{Module: "users", Action: action, Type: "movies", UserName: "me", ID: "55", Section: "calendar"}
	}
	undoReset := str.Options{InternalID: "55", Delete: true}
	unlike := users("list_like")
	unlike.Type = consts.EmptyString
	unlike.Delete = true
	like := unlike
	like.Delete = false
	updateItem := like
	updateItem.Action = "update_list_item"
	updateItem.ListItemID = 7
	updateItem.Notes = "rewatch"
	return []writeHandler{
		{name: "shows reset show progress", handler: ShowsResetShowProgressHandler{}, options: str.Options{InternalID: "55"}, method: http.MethodPost, path: "/shows/55/progress/watched/reset"},
		{name: "shows undo reset show progress", handler: ShowsResetShowProgressHandler{}, options: undoReset, method: http.MethodDelete, path: "/shows/55/progress/watched/reset", status: http.StatusNoContent},
		{name: "sync add to collection", handler: SyncAddToCollectionHandler{}, options: sync("add_to_collection"), items: movies, sends: movie, method: http.MethodPost, path: "/sync/collection", status: http.StatusCreated, writes: true},
		{name: "sync add to favorites", handler: SyncAddToFavoritesHandler{}, options: sync("add_to_favorites"), items: movies, sends: movie, method: http.MethodPost, path: "/sync/favorites", status: http.StatusCreated, writes: true},
		{name: "sync add to ratings", handler: SyncAddToRatingsHandler{}, options: sync("add_to_ratings"), items: movies, sends: movie, method: http.MethodPost, path: "/sync/ratings", status: http.StatusCreated, writes: true},
		{name: "sync add to watchlist", handler: SyncAddToWatchlistHandler{}, options: sync("add_to_watchlist"), items: movies, sends: movie, method: http.MethodPost, path: "/sync/watchlist", status: http.StatusCreated, writes: true},
		{name: "sync remove from collection", handler: SyncRemoveFromCollectionHandler{}, options: sync("remove_from_collection"), items: movies, sends: movie, method: http.MethodPost, path: "/sync/collection/remove", writes: true},
		{name: "sync remove from favorites", handler: SyncRemoveFromFavoritesHandler{}, options: sync("remove_from_favorites"), items: movies, sends: movie, method: http.MethodPost, path: "/sync/favorites/remove", writes: true},
		{name: "sync remove from history", handler: SyncRemoveFromHistoryHandler{}, options: sync("remove_from_history"), items: movies, sends: movie, method: http.MethodPost, path: "/sync/history/remove", writes: true},
		{name: "sync remove from ratings", handler: SyncRemoveFromRatingsHandler{}, options: sync("remove_from_ratings"), items: movies, sends: movie, method: http.MethodPost, path: "/sync/ratings/remove", writes: true},
		{name: "sync remove from watchlist", handler: SyncRemoveFromWatchlistHandler{}, options: sync("remove_from_watchlist"), items: movies, sends: movie, method: http.MethodPost, path: "/sync/watchlist/remove", writes: true},
		{name: "sync reorder favorites", handler: SyncReorderFavoritesHandler{}, options: sync("reorder_favorites"), items: movies, sends: `{"rank":[7]}`, method: http.MethodPost, path: "/sync/favorites/reorder", writes: true},
		{name: "sync reorder watchlist", handler: SyncReorderWatchlistHandler{}, options: sync("reorder_watchlist"), items: movies, sends: `{"rank":[7]}`, method: http.MethodPost, path: "/sync/watchlist/reorder", writes: true},
		{name: "sync update favorite item", handler: SyncUpdateFavoriteItemHandler{}, options: str.Options{ListItemID: 7, Notes: "rewatch"}, sends: note, method: http.MethodPut, path: "/sync/favorites/7", status: http.StatusNoContent},
		{name: "sync update favorites", handler: SyncUpdateFavoritesHandler{}, options: sync("update_favorites"), sends: `"sort_by"`, method: http.MethodPut, path: "/sync/favorites", writes: true},
		{name: "sync update watchlist", handler: SyncUpdateWatchlistHandler{}, options: sync("update_watchlist"), sends: `"sort_by"`, method: http.MethodPut, path: "/sync/watchlist", writes: true},
		{name: "sync update watchlist item", handler: SyncUpdateWatchlistItemHandler{}, options: str.Options{ListItemID: 7, Notes: "rewatch"}, sends: note, method: http.MethodPut, path: "/sync/watchlist/7", status: http.StatusNoContent},
		{name: "users add hidden items", handler: UsersAddHiddenItemsHandler{}, options: users("add_hidden_items"), items: movies, sends: movie, method: http.MethodPost, path: "/users/hidden/calendar", status: http.StatusCreated, writes: true},
		{name: "users add list items", handler: UsersAddListItemsHandler{}, options: users("add_list_items"), items: movies, sends: movie, method: http.MethodPost, path: "/users/me/lists/55/items", status: http.StatusCreated, writes: true},
		{name: "users list like", handler: UsersListLikeHandler{}, options: like, method: http.MethodPost, path: "/users/me/lists/55/like", status: http.StatusNoContent},
		{name: "users remove hidden items", handler: UsersRemoveHiddenItemsHandler{}, options: users("remove_hidden_items"), items: movies, sends: movie, method: http.MethodPost, path: "/users/hidden/calendar/remove", writes: true},
		{name: "users remove list items", handler: UsersRemoveListItemsHandler{}, options: users("remove_list_items"), items: movies, sends: movie, method: http.MethodPost, path: "/users/me/lists/55/items/remove", writes: true},
		{name: "users remove list like", handler: UsersListLikeHandler{}, options: unlike, method: http.MethodDelete, path: "/users/me/lists/55/like", status: http.StatusNoContent},
		{name: "users reorder list items", handler: UsersReorderListItemsHandler{}, options: users("reorder_list_items"), items: movies, sends: `{"rank":[7]}`, method: http.MethodPost, path: "/users/me/lists/55/items/reorder", writes: true},
		{name: "users reorder lists", handler: UsersReorderListsHandler{}, options: users("reorder_lists"), items: lists, sends: `{"rank":[1]}`, method: http.MethodPost, path: "/users/me/lists/reorder", writes: true},
		{name: "users update list item", handler: UsersUpdateListItemHandler{}, options: updateItem, sends: note, method: http.MethodPut, path: "/users/me/lists/55/items/7", status: http.StatusNoContent},
	}
}

// writeOptions returns the options of a case with -o and, when the case has items, -items in a temporary directory.
func writeOptions(t *testing.T, tc writeHandler) str.Options {
	t.Helper()
	dir := t.TempDir()
	options := tc.options
	options.Output = filepath.Join(dir, "out.json")
	if tc.items != consts.EmptyString {
		options.Items = filepath.Join(dir, "items.json")
		test.AssertNilError(t, os.WriteFile(options.Items, []byte(tc.items), 0o600))
	}
	return options
}

func TestWriteHandlersSendRequest(t *testing.T) {
	for _, tc := range writeHandlers() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			requests := []string{}
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				requests = append(requests, r.Method+" "+r.URL.Path)
				body, err := io.ReadAll(r.Body)
				test.AssertNilError(t, err)
				assert.Contains(t, string(body), tc.sends)
				if tc.status != consts.ZeroValue {
					w.WriteHeader(tc.status)
				}
				if tc.status != http.StatusNoContent {
					test.SafeFprint(w, `{}`)
				}
			})

			options := writeOptions(t, tc)
			test.AssertNilError(t, tc.handler.Handle(&options, s.Client))
			assert.Equal(t, []string{tc.method + " " + tc.path}, requests)

			data, err := os.ReadFile(options.Output)
			if !tc.writes {
				assert.True(t, os.IsNotExist(err), "no result file for this action")
				return
			}
			test.AssertNilError(t, err)
			assert.True(t, json.Valid(data), "result written to -o is JSON")
		})
	}
}

func TestWriteHandlersFailedRequest(t *testing.T) {
	for _, tc := range writeHandlers() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			options := writeOptions(t, tc)
			assert.Error(t, tc.handler.Handle(&options, s.Client))
			_, err := os.Stat(options.Output)
			assert.True(t, os.IsNotExist(err), "nothing is written when the request fails")
		})
	}
}
