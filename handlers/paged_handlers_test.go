// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/stretchr/testify/assert"
)

// pagedHandler is a handler that follows the pages of one list route and writes the collected items to -o.
type pagedHandler struct {
	name    string
	handler Handler
	options str.Options
	item    string
}

func pagedHandlers() []pagedHandler {
	byID := str.Options{InternalID: "55"}
	period := str.Options{Period: "weekly"}
	since := str.Options{StartDate: "2026-10-01"}
	return []pagedHandler{
		{name: "comments replies", handler: CommentsRepliesHandler{}, options: str.Options{CommentID: 417}},
		{name: "episodes comments", handler: EpisodesCommentsHandler{}, options: str.Options{InternalID: "55", Season: 1, Episode: 2}},
		{name: "episodes lists", handler: EpisodesListsHandler{}, options: str.Options{Module: "episodes", Action: "lists", InternalID: "55", Season: 1, Episode: 2}},
		{name: "lists comments", handler: ListsCommentsHandler{}, options: byID},
		{name: "lists items", handler: ListsItemsHandler{}, options: byID},
		{name: "lists likes", handler: ListsLikesHandler{}, options: byID},
		{name: "lists popular", handler: ListsPopularHandler{}},
		{name: "lists trending", handler: ListsTrendingHandler{}},
		{name: "movies anticipated", handler: MoviesAnticipatedHandler{}},
		{name: "movies collected", handler: MoviesCollectedHandler{}, options: period},
		{name: "movies comments", handler: MoviesCommentsHandler{}, options: byID},
		{name: "movies favorited", handler: MoviesFavoritedHandler{}, options: period},
		{name: "movies lists", handler: MoviesListsHandler{}, options: str.Options{Module: "movies", Action: "lists", InternalID: "55"}},
		{name: "movies played", handler: MoviesPlayedHandler{}, options: period},
		{name: "movies popular", handler: MoviesPopularHandler{}},
		{name: "movies related", handler: MoviesRelatedHandler{}, options: byID},
		{name: "movies trending", handler: MoviesTrendingHandler{}},
		{name: "movies updated ids", handler: MoviesUpdatedIDsHandler{}, options: since, item: "1"},
		{name: "movies updates", handler: MoviesUpdatesHandler{}, options: since},
		{name: "movies watched", handler: MoviesWatchedHandler{}, options: period},
		{name: "networks lists", handler: NetworksListsHandler{}},
		{name: "people lists", handler: PeopleListsHandler{}, options: str.Options{ID: "55"}},
		{name: "people updated ids", handler: PeopleUpdatedIDsHandler{}, options: since, item: "1"},
		{name: "people updates", handler: PeopleUpdatesHandler{}, options: since},
		{name: "search exact query", handler: SearchExactQueryHandler{}, options: str.Options{SearchType: str.Slice{"movie"}, Query: "tron"}},
		{name: "search text query", handler: SearchTextQueryHandler{}, options: str.Options{Module: "search", Action: "text_query", SearchType: str.Slice{"movie"}, Query: "tron"}},
		{name: "search trending", handler: SearchTrendingHandler{}, options: str.Options{SearchType: str.Slice{"movies"}}},
		{name: "seasons comments", handler: SeasonsCommentsHandler{}, options: str.Options{InternalID: "55", Season: 1}},
		{name: "seasons lists", handler: SeasonsListsHandler{}, options: str.Options{Module: "seasons", Action: "lists", InternalID: "55", Season: 1}},
		{name: "shows anticipated", handler: ShowsAnticipatedHandler{}},
		{name: "shows collected", handler: ShowsCollectedHandler{}, options: period},
		{name: "shows comments", handler: ShowsCommentsHandler{}, options: byID},
		{name: "shows favorited", handler: ShowsFavoritedHandler{}, options: period},
		{name: "shows lists", handler: ShowsListsHandler{}, options: str.Options{Module: "shows", Action: "lists", InternalID: "55"}},
		{name: "shows played", handler: ShowsPlayedHandler{}, options: period},
		{name: "shows popular", handler: ShowsPopularHandler{}},
		{name: "shows related", handler: ShowsRelatedHandler{}, options: byID},
		{name: "shows trending", handler: ShowsTrendingHandler{}},
		{name: "shows updated ids", handler: ShowsUpdatedIDsHandler{}, options: since, item: "1"},
		{name: "shows updates", handler: ShowsUpdatesHandler{}, options: since},
		{name: "shows watched", handler: ShowsWatchedHandler{}, options: period},
		{name: "smart lists items", handler: SmartListsItemsHandler{}, options: byID},
	}
}

// servePagedItems answers every GET with one item and pagination headers for pageCount pages; a non-empty
// failFrom makes that page and the ones after it fail. It returns the requested page numbers.
func servePagedItems(t *testing.T, mux *http.ServeMux, item string, pageCount string, failFrom string) *[]string {
	t.Helper()
	if item == consts.EmptyString {
		item = `{}`
	}
	pages := []string{}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		page := r.URL.Query().Get("page")
		pages = append(pages, page)
		if failFrom != consts.EmptyString && page >= failFrom {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set(trakt.HeaderPaginationPage, page)
		w.Header().Set(trakt.HeaderPaginationPageCount, pageCount)
		test.SafeFprint(w, "["+item+"]")
	})
	return &pages
}

func TestPagedHandlersWriteAllPages(t *testing.T) {
	for _, tc := range pagedHandlers() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			pages := servePagedItems(t, s.Mux, tc.item, "3", consts.EmptyString)

			options := tc.options
			options.Output = filepath.Join(t.TempDir(), "out.json")
			test.AssertNilError(t, tc.handler.Handle(&options, s.Client))
			assert.Equal(t, []string{"1", "2", "3"}, *pages)

			data, err := os.ReadFile(options.Output)
			test.AssertNilError(t, err)
			var written []json.RawMessage
			test.AssertNilError(t, json.Unmarshal(data, &written))
			assert.Len(t, written, 3, "items written to -o")
		})
	}
}

func TestPagedHandlersStopAtPagesLimit(t *testing.T) {
	for _, tc := range pagedHandlers() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			pages := servePagedItems(t, s.Mux, tc.item, "5", consts.EmptyString)

			options := tc.options
			options.PagesLimit = 2
			options.Output = filepath.Join(t.TempDir(), "out.json")
			test.AssertNilError(t, tc.handler.Handle(&options, s.Client))
			assert.Equal(t, []string{"1", "2"}, *pages)
		})
	}
}

func TestPagedHandlersFailedPage(t *testing.T) {
	for _, failFrom := range []string{"1", "2"} {
		for _, tc := range pagedHandlers() {
			tc, failFrom := tc, failFrom
			t.Run(tc.name+" page "+failFrom, func(t *testing.T) {
				s := setup(t)
				defer s.Teardown()
				servePagedItems(t, s.Mux, tc.item, "3", failFrom)

				options := tc.options
				options.Output = filepath.Join(t.TempDir(), "out.json")
				assert.Error(t, tc.handler.Handle(&options, s.Client))
				_, err := os.Stat(options.Output)
				assert.True(t, os.IsNotExist(err), "nothing is written when a page fails")
			})
		}
	}
}
