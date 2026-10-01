// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/stretchr/testify/assert"
)

// pagedFetch starts a paginated wrapper at the first page and returns how many items it collected.
type pagedFetch func(client *trakt.Client, options *str.Options) (int, error)

func countItems[T any](fetch func(client *trakt.Client, options *str.Options, page int) ([]T, error)) pagedFetch {
	return func(client *trakt.Client, options *str.Options) (int, error) {
		list, err := fetch(client, options, consts.DefaultPage)
		return len(list), err
	}
}

func pagedFetches() map[string]pagedFetch {
	c := &CommonLogic{}
	return map[string]pagedFetch{
		"FetchBlockedUsers":           countItems(c.FetchBlockedUsers),
		"FetchCommentUserLikes":       countItems(c.FetchCommentUserLikes),
		"FetchFavorites":              countItems(c.FetchFavorites),
		"FetchFollowers":              countItems(c.FetchFollowers),
		"FetchFollowing":              countItems(c.FetchFollowing),
		"FetchFriends":                countItems(c.FetchFriends),
		"FetchHistoryList":            countItems(c.FetchHistoryList),
		"FetchMovieRecommendations":   countItems(c.FetchMovieRecommendations),
		"FetchRatings":                countItems(c.FetchRatings),
		"FetchRecentComments":         countItems(c.FetchRecentComments),
		"FetchShowRecommendations":    countItems(c.FetchShowRecommendations),
		"FetchTrendingComments":       countItems(c.FetchTrendingComments),
		"FetchUpdatedComments":        countItems(c.FetchUpdatedComments),
		"FetchUsersCollaborations":    countItems(c.FetchUsersCollaborations),
		"FetchUsersCollection":        countItems(c.FetchUsersCollection),
		"FetchUsersComments":          countItems(c.FetchUsersComments),
		"FetchUsersFavorites":         countItems(c.FetchUsersFavorites),
		"FetchUsersFavoritesComments": countItems(c.FetchUsersFavoritesComments),
		"FetchUsersHiddenItems":       countItems(c.FetchUsersHiddenItems),
		"FetchUsersHistory":           countItems(c.FetchUsersHistory),
		"FetchUsersLikes":             countItems(c.FetchUsersLikes),
		"FetchUsersListComments":      countItems(c.FetchUsersListComments),
		"FetchUsersListItems":         countItems(c.FetchUsersListItems),
		"FetchUsersListLikes":         countItems(c.FetchUsersListLikes),
		"FetchUsersNotes":             countItems(c.FetchUsersNotes),
		"FetchUsersRatings":           countItems(c.FetchUsersRatings),
		"FetchUsersWatchlist":         countItems(c.FetchUsersWatchlist),
		"FetchUsersWatchlistComments": countItems(c.FetchUsersWatchlistComments),
		"FetchWatchlist":              countItems(c.FetchWatchlist),
	}
}

func pagedOptions() *str.Options {
	return &str.Options{
		UserName: "sean", ID: "star-wars", Type: "movies", CommentType: "reviews", CommentID: 417,
		Section: "calendar", Sort: "newest", SortBy: "rank", SortHow: "asc", PerPage: 10,
	}
}

// servePages answers every request with one item and pagination headers for pageCount pages;
// failFrom > 0 makes that page and the ones after it fail. It returns the requested page numbers.
func servePages(t *testing.T, mux *http.ServeMux, pageCount string, failFrom string) *[]string {
	t.Helper()
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
		test.SafeFprint(w, `[{}]`)
	})
	return &pages
}

func TestCommonFetchFollowsPages(t *testing.T) {
	for name, fetch := range pagedFetches() {
		name, fetch := name, fetch
		t.Run(name, func(t *testing.T) {
			testSetup := setup(t)
			defer testSetup.Teardown()
			pages := servePages(t, testSetup.Mux, "3", consts.EmptyString)

			count, err := fetch(testSetup.Client, pagedOptions())
			test.AssertNilError(t, err)
			assert.Equal(t, 3, count, "items from all pages")
			assert.Equal(t, []string{"1", "2", "3"}, *pages)
		})
	}
}

func TestCommonFetchStopsAtPagesLimit(t *testing.T) {
	for name, fetch := range pagedFetches() {
		name, fetch := name, fetch
		t.Run(name, func(t *testing.T) {
			testSetup := setup(t)
			defer testSetup.Teardown()
			pages := servePages(t, testSetup.Mux, "5", consts.EmptyString)

			options := pagedOptions()
			options.PagesLimit = 2
			count, err := fetch(testSetup.Client, options)
			test.AssertNilError(t, err)
			assert.Equal(t, 2, count, "items up to the pages limit")
			assert.Equal(t, []string{"1", "2"}, *pages)
		})
	}
}

func TestCommonFetchNextPageError(t *testing.T) {
	for name, fetch := range pagedFetches() {
		name, fetch := name, fetch
		t.Run(name, func(t *testing.T) {
			testSetup := setup(t)
			defer testSetup.Teardown()
			pages := servePages(t, testSetup.Mux, "3", "2")

			count, err := fetch(testSetup.Client, pagedOptions())
			assert.Error(t, err)
			assert.Zero(t, count, "no partial list on a failed page")
			assert.Equal(t, []string{"1", "2"}, *pages)
		})
	}
}
