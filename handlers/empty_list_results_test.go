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
		module     string
		name       string
		handler    Handler
		searchType str.Slice
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
		{module: "movies", name: "anticipated_recommendations", handler: MoviesAnticipatedRecommendationsHandler{}},
		{module: "shows", name: "anticipated_recommendations", handler: ShowsAnticipatedRecommendationsHandler{}},
		{module: "media", name: "anticipated_recommendations", handler: MediaAnticipatedRecommendationsHandler{}},
		{module: "media", name: "trending", handler: MediaTrendingHandler{}},
		{module: "media", name: "popular", handler: MediaPopularHandler{}},
		{module: "media", name: "anticipated", handler: MediaAnticipatedHandler{}},
		{module: "calendars", name: "my-shows", handler: CalendarsShowsHandler{}},
		{module: "calendars", name: "my-new-shows", handler: CalendarsNewShowsHandler{}},
		{module: "calendars", name: "my-season-premieres", handler: CalendarsSeasonPremieresHandler{}},
		{module: "calendars", name: "my-finales", handler: CalendarsFinalesHandler{}},
		{module: "calendars", name: "my-movies", handler: CalendarsMoviesHandler{}},
		{module: "calendars", name: "my-dvd", handler: CalendarsDvdHandler{}},
		{module: "calendars", name: "my-media", handler: CalendarsMediaHandler{}},
		{module: "calendars", name: "hot-finales", handler: CalendarsHotFinalesHandler{}},
		{module: "calendars", name: "hot-releases", handler: CalendarsHotReleasesHandler{}},
		{module: "calendars", name: "hot-premieres", handler: CalendarsHotPremieresHandler{}},
		{module: "calendars", name: "hot-new-shows", handler: CalendarsHotNewShowsHandler{}},
		{module: "calendars", name: "streaming", handler: CalendarsStreamingHandler{}},
		{module: "recommendations", name: "movies", handler: RecommendationsMoviesHandler{}},
		{module: "recommendations", name: "shows", handler: RecommendationsShowsHandler{}},
		{module: "social_recommendations", name: "movies", handler: SocialRecommendationsMoviesHandler{}},
		{module: "social_recommendations", name: "shows", handler: SocialRecommendationsShowsHandler{}},
		{module: "smart_lists", name: "templates", handler: SmartListsTemplatesHandler{}},
		{module: "smart_lists", name: "items", handler: SmartListsItemsHandler{}},
		{module: "team", name: "members", handler: TeamMembersHandler{}},
		{module: "watchnow", name: "sources", handler: WatchNowSourcesHandler{}},
		{module: "younify", name: "connections", handler: YounifyConnectionsHandler{}},
		{module: "search", name: "text_query", handler: SearchTextQueryHandler{}, searchType: str.Slice{"movie"}},
		{module: "search", name: "id_lookup", handler: SearchIDLookupHandler{}},
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
			options := str.Options{Module: h.module, Action: h.name, Period: "weekly", ID: "keanu-reeves", InternalID: "123", CommentID: 1, SearchType: h.searchType, Query: "tron", Output: output}
			err := h.handler.Handle(&options, s.Client)
			test.AssertNilError(t, err)
			assert.Equal(t, 1, calls, "one request")
			data, err := os.ReadFile(output)
			test.AssertNilError(t, err)
			assert.Equal(t, "[]", string(data))
		})
	}
}
