// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// watchlist_comments and favorites_comments export the whole comment, not only its id and user.
func TestUsersCommentsHandlersKeepCommentFields(t *testing.T) {
	const comment = `[{"id":8,"parent_id":0,"created_at":"2026-10-01T10:30:00Z","updated_at":"2026-10-02T11:30:00Z","comment":"Great list, thanks!",` +
		`"spoiler":false,"review":false,"replies":2,"likes":5,"user_stats":{"rating":9,"play_count":1,"completed_count":1},"user":{"username":"justin","ids":{"slug":"justin"}}}]`
	handlers := []struct {
		action  string
		path    string
		handler Handler
	}{
		{action: "watchlist_comments", path: "/users/sean/watchlist/comments/likes", handler: UsersWatchlistCommentsHandler{}},
		{action: "favorites_comments", path: "/users/sean/favorites/comments/likes", handler: UsersFavoritesCommentsHandler{}},
	}

	for _, h := range handlers {
		h := h
		t.Run(h.action, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			s.Mux.HandleFunc(h.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				test.SafeFprint(w, comment)
			})

			options := str.Options{Module: "users", Action: h.action, UserName: "sean", Sort: "likes", Output: filepath.Join(t.TempDir(), "out.json")}
			test.AssertNilError(t, h.handler.Handle(&options, s.Client))
			data, err := os.ReadFile(options.Output)
			test.AssertNilError(t, err)
			assert.JSONEq(t, comment, string(data))
		})
	}
}
