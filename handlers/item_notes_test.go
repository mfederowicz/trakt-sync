// Package handlers used to handle module actions
package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// notedItems has one list item of each type with its own notes, as users -a lists exports them
const notedItems = `[
{"id":1,"type":"movie","notes":"movie note","movie":{"title":"Tron","ids":{"trakt":1}}},
{"id":2,"type":"show","notes":"show note","show":{"title":"Dark","ids":{"trakt":2}}},
{"id":3,"type":"season","notes":"season note","season":{"number":1,"ids":{"trakt":3}}},
{"id":4,"type":"episode","notes":"episode note","episode":{"season":1,"number":1,"ids":{"trakt":4}}}
]`

func TestAddItemsSendNotes(t *testing.T) {
	users := func(action string) str.Options {
		return str.Options{Module: "users", Action: action, UserName: "me", ID: "55"}
	}
	sync := func(action string) str.Options {
		return str.Options{Module: "sync", Action: action}
	}
	cases := []struct {
		name    string
		handler Handler
		options str.Options
		path    string
		notes   bool
	}{
		{name: "users add list items", handler: UsersAddListItemsHandler{}, options: users("add_list_items"), path: "/users/me/lists/55/items", notes: true},
		{name: "sync add to watchlist", handler: SyncAddToWatchlistHandler{}, options: sync("add_to_watchlist"), path: "/sync/watchlist", notes: true},
		{name: "users remove list items", handler: UsersRemoveListItemsHandler{}, options: users("remove_list_items"), path: "/users/me/lists/55/items/remove"},
		{name: "sync remove from watchlist", handler: SyncRemoveFromWatchlistHandler{}, options: sync("remove_from_watchlist"), path: "/sync/watchlist/remove"},
		{name: "sync add to collection", handler: SyncAddToCollectionHandler{}, options: sync("add_to_collection"), path: "/sync/collection"},
		{name: "sync add to favorites", handler: SyncAddToFavoritesHandler{}, options: sync("add_to_favorites"), path: "/sync/favorites"},
	}
	types := map[string]string{"movies": "movie note", "shows": "show note", "seasons": "season note", "episodes": "episode note"}
	for _, tc := range cases {
		tc := tc
		for stype, note := range types {
			stype, note := stype, note
			t.Run(tc.name+" "+stype, func(t *testing.T) {
				s := setup(t)
				defer s.Teardown()
				bodies := []string{}
				s.Mux.HandleFunc(tc.path, func(w http.ResponseWriter, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					test.AssertNilError(t, err)
					bodies = append(bodies, string(body))
					test.SafeFprint(w, `{}`)
				})

				dir := t.TempDir()
				options := tc.options
				options.Type = stype
				options.Output = filepath.Join(dir, "out.json")
				options.Items = filepath.Join(dir, "items.json")
				test.AssertNilError(t, os.WriteFile(options.Items, []byte(notedItems), 0o600))
				test.AssertNilError(t, tc.handler.Handle(&options, s.Client))

				if !assert.Len(t, bodies, 1) {
					return
				}
				if tc.notes {
					assert.Contains(t, bodies[0], `"notes":"`+note+`"`)
					return
				}
				assert.NotContains(t, bodies[0], `"notes"`)
			})
		}
	}
}
