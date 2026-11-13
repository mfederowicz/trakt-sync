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

// a list action that gets an empty list from Trakt writes [] and returns no error.
func TestListHandlersEmptyResult(t *testing.T) {
	handlers := []struct {
		module  string
		name    string
		handler Handler
	}{
		{module: "movies", name: "trending", handler: MoviesTrendingHandler{}},
		{module: "movies", name: "popular", handler: MoviesPopularHandler{}},
		{module: "movies", name: "anticipated", handler: MoviesAnticipatedHandler{}},
		{module: "movies", name: "boxoffice", handler: MoviesBoxofficeHandler{}},
		{module: "movies", name: "watched", handler: MoviesWatchedHandler{}},
		{module: "movies", name: "favorited", handler: MoviesFavoritedHandler{}},
		{module: "movies", name: "played", handler: MoviesPlayedHandler{}},
		{module: "movies", name: "collected", handler: MoviesCollectedHandler{}},
		{module: "movies", name: "hot", handler: MoviesHotHandler{}},
		{module: "movies", name: "streaming", handler: MoviesStreamingHandler{}},
		{module: "movies", name: "updates", handler: MoviesUpdatesHandler{}},
		{module: "movies", name: "updated_ids", handler: MoviesUpdatedIDsHandler{}},
		{module: "shows", name: "trending", handler: ShowsTrendingHandler{}},
		{module: "shows", name: "popular", handler: ShowsPopularHandler{}},
		{module: "shows", name: "anticipated", handler: ShowsAnticipatedHandler{}},
		{module: "shows", name: "watched", handler: ShowsWatchedHandler{}},
		{module: "shows", name: "favorited", handler: ShowsFavoritedHandler{}},
		{module: "shows", name: "played", handler: ShowsPlayedHandler{}},
		{module: "shows", name: "collected", handler: ShowsCollectedHandler{}},
		{module: "shows", name: "updates", handler: ShowsUpdatesHandler{}},
		{module: "shows", name: "updated_ids", handler: ShowsUpdatedIDsHandler{}},
		{module: "people", name: "lists", handler: PeopleListsHandler{}},
		{module: "people", name: "updates", handler: PeopleUpdatesHandler{}},
		{module: "people", name: "updated_ids", handler: PeopleUpdatedIDsHandler{}},
		{module: "networks", name: "lists", handler: NetworksListsHandler{}},
		{module: "lists", name: "popular", handler: ListsPopularHandler{}},
		{module: "lists", name: "trending", handler: ListsTrendingHandler{}},
		{module: "lists", name: "likes", handler: ListsLikesHandler{}},
		{module: "lists", name: "items", handler: ListsItemsHandler{}},
		{module: "lists", name: "comments", handler: ListsCommentsHandler{}},
		{module: "comments", name: "recent", handler: CommentsRecentHandler{}},
		{module: "comments", name: "trending", handler: CommentsTrendingHandler{}},
		{module: "comments", name: "updates", handler: CommentsUpdatesHandler{}},
		{module: "comments", name: "replies", handler: CommentsRepliesHandler{}},
		{module: "comments", name: "likes", handler: CommentsLikesHandler{}},
		{module: "comments", name: "reactions", handler: CommentsReactionsHandler{}},
	}

	for _, h := range handlers {
		h := h
		t.Run(h.module+" "+h.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			calls := 0
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				calls++
				test.SafeFprint(w, `[]`)
			})

			output := filepath.Join(t.TempDir(), "out.json")
			options := str.Options{Module: h.module, Action: h.name, Period: "weekly", ID: "keanu-reeves", InternalID: "123", CommentID: 1, Output: output}
			err := h.handler.Handle(&options, s.Client)
			test.AssertNilError(t, err)
			assert.Equal(t, 1, calls, "one request")
			data, err := os.ReadFile(output)
			test.AssertNilError(t, err)
			assert.Equal(t, "[]", string(data))
		})
	}
}
