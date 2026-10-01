// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// singleRequestHandler is a handler that sends one GET request and writes the decoded response to -o.
type singleRequestHandler struct {
	name    string
	handler Handler
	options str.Options
	path    string
	body    string
}

func singleRequestHandlers() []singleRequestHandler {
	const (
		object = `{}`
		list   = `[{}]`
	)
	byID := str.Options{InternalID: "55"}
	calendar := str.Options{StartDate: "2026-10-01", Days: 7}
	comment := str.Options{CommentID: 417}
	episode := str.Options{InternalID: "55", Season: 1, Episode: 2}
	person := str.Options{ID: "55"}
	season := str.Options{InternalID: "55", Season: 1}
	typed := str.Options{Type: "movies"}
	user := str.Options{UserName: "me"}
	return []singleRequestHandler{
		{name: "calendars dvd", handler: CalendarsDvdHandler{}, options: calendar, path: "/calendars/my/dvd/2026-10-01/7", body: list},
		{name: "calendars finales", handler: CalendarsFinalesHandler{}, options: calendar, path: "/calendars/my/shows/finales/2026-10-01/7", body: list},
		{name: "calendars hot finales", handler: CalendarsHotFinalesHandler{}, options: calendar, path: "/calendars/releases/hot/finales/2026-10-01/7", body: list},
		{name: "calendars hot new shows", handler: CalendarsHotNewShowsHandler{}, options: calendar, path: "/calendars/releases/hot/new/2026-10-01/7", body: list},
		{name: "calendars hot premieres", handler: CalendarsHotPremieresHandler{}, options: calendar, path: "/calendars/releases/hot/premieres/2026-10-01/7", body: list},
		{name: "calendars hot releases", handler: CalendarsHotReleasesHandler{}, options: calendar, path: "/calendars/releases/hot/2026-10-01/7", body: list},
		{name: "calendars media", handler: CalendarsMediaHandler{}, options: calendar, path: "/calendars/my/media/2026-10-01/7", body: list},
		{name: "calendars movies", handler: CalendarsMoviesHandler{}, options: calendar, path: "/calendars/my/movies/2026-10-01/7", body: list},
		{name: "calendars new shows", handler: CalendarsNewShowsHandler{}, options: calendar, path: "/calendars/my/shows/new/2026-10-01/7", body: list},
		{name: "calendars season premieres", handler: CalendarsSeasonPremieresHandler{}, options: calendar, path: "/calendars/my/shows/premieres/2026-10-01/7", body: list},
		{name: "calendars shows", handler: CalendarsShowsHandler{}, options: calendar, path: "/calendars/my/shows/2026-10-01/7", body: list},
		{name: "calendars streaming", handler: CalendarsStreamingHandler{}, options: calendar, path: "/calendars/my/streaming/2026-10-01/7", body: list},
		{name: "certifications types", handler: CertificationsTypesHandler{}, options: typed, path: "/certifications/movies", body: object},
		{name: "comments comment", handler: CommentsCommentHandler{}, options: comment, path: "/comments/417", body: object},
		{name: "comments item", handler: CommentsItemHandler{}, options: comment, path: "/comments/417/item", body: object},
		{name: "comments likes", handler: CommentsLikesHandler{}, options: comment, path: "/comments/417/likes", body: list},
		{name: "comments reactions", handler: CommentsReactionsHandler{}, options: comment, path: "/comments/417/reactions", body: list},
		{name: "comments reactions summary", handler: CommentsReactionsSummaryHandler{}, options: comment, path: "/comments/417/reactions/summary", body: object},
		{name: "comments recent", handler: CommentsRecentHandler{}, options: typed, path: "/comments/recent/movies", body: list},
		{name: "comments trending", handler: CommentsTrendingHandler{}, options: typed, path: "/comments/trending/movies", body: list},
		{name: "comments updates", handler: CommentsUpdatesHandler{}, options: typed, path: "/comments/updates/movies", body: list},
		{name: "countries types", handler: CountriesTypesHandler{}, options: typed, path: "/countries/movies", body: list},
		{name: "episodes people", handler: EpisodesPeopleHandler{}, options: episode, path: "/shows/55/seasons/1/episodes/2/people", body: object},
		{name: "episodes ratings", handler: EpisodesRatingsHandler{}, options: episode, path: "/shows/55/seasons/1/episodes/2/ratings", body: object},
		{name: "episodes stats", handler: EpisodesStatsHandler{}, options: episode, path: "/shows/55/seasons/1/episodes/2/stats", body: object},
		{name: "episodes summary", handler: EpisodesSummaryHandler{}, options: episode, path: "/shows/55/seasons/1/episodes/2", body: object},
		{name: "episodes translations", handler: EpisodesTranslationsHandler{}, options: str.Options{InternalID: "55", Season: 1, Episode: 2, Language: "en"}, path: "/shows/55/seasons/1/episodes/2/translations/en", body: list},
		{name: "episodes videos", handler: EpisodesVideosHandler{}, options: episode, path: "/shows/55/seasons/1/episodes/2/videos", body: list},
		{name: "episodes watching", handler: EpisodesWatchingHandler{}, options: episode, path: "/shows/55/seasons/1/episodes/2/watching", body: list},
		{name: "genres types", handler: GenresTypesHandler{}, options: typed, path: "/genres/movies", body: list},
		{name: "languages types", handler: LanguagesTypesHandler{}, options: typed, path: "/languages/movies", body: list},
		{name: "media anticipated", handler: MediaAnticipatedHandler{}, path: "/media/anticipated", body: list},
		{name: "media popular", handler: MediaPopularHandler{}, path: "/media/popular", body: list},
		{name: "movies aliases", handler: MoviesAliasesHandler{}, options: byID, path: "/movies/55/aliases", body: list},
		{name: "movies boxoffice", handler: MoviesBoxofficeHandler{}, path: "/movies/boxoffice", body: list},
		{name: "movies people", handler: MoviesPeopleHandler{}, options: byID, path: "/movies/55/people", body: object},
		{name: "movies ratings", handler: MoviesRatingsHandler{}, options: byID, path: "/movies/55/ratings", body: object},
		{name: "movies releases", handler: MoviesReleasesHandler{}, options: str.Options{InternalID: "55", Country: "us"}, path: "/movies/55/releases/us", body: list},
		{name: "movies stats", handler: MoviesStatsHandler{}, options: byID, path: "/movies/55/stats", body: object},
		{name: "movies studios", handler: MoviesStudiosHandler{}, options: byID, path: "/movies/55/studios", body: list},
		{name: "movies translations", handler: MoviesTranslationsHandler{}, options: str.Options{InternalID: "55", Language: "en"}, path: "/movies/55/translations/en", body: list},
		{name: "movies videos", handler: MoviesVideosHandler{}, options: byID, path: "/movies/55/videos", body: list},
		{name: "movies watching", handler: MoviesWatchingHandler{}, options: byID, path: "/movies/55/watching", body: list},
		{name: "notes item", handler: NotesItemHandler{}, options: byID, path: "/notes/55/item", body: object},
		{name: "notes note", handler: NotesNoteHandler{}, options: byID, path: "/notes/55", body: object},
		{name: "people movies", handler: PeopleMoviesHandler{}, options: person, path: "/people/55/movies", body: object},
		{name: "people shows", handler: PeopleShowsHandler{}, options: person, path: "/people/55/shows", body: object},
		{name: "people summary", handler: PeopleSummaryHandler{}, options: person, path: "/people/55", body: object},
		{name: "seasons episodes", handler: SeasonsEpisodesHandler{}, options: season, path: "/shows/55/seasons/1", body: list},
		{name: "seasons people", handler: SeasonsPeopleHandler{}, options: season, path: "/shows/55/seasons/1/people", body: object},
		{name: "seasons ratings", handler: SeasonsRatingsHandler{}, options: season, path: "/shows/55/seasons/1/ratings", body: object},
		{name: "seasons season", handler: SeasonsSeasonHandler{}, options: season, path: "/shows/55/seasons/1/info", body: object},
		{name: "seasons stats", handler: SeasonsStatsHandler{}, options: season, path: "/shows/55/seasons/1/stats", body: object},
		{name: "seasons summary", handler: SeasonsSummaryHandler{}, options: byID, path: "/shows/55/seasons", body: list},
		{name: "seasons translations", handler: SeasonsTranslationsHandler{}, options: str.Options{InternalID: "55", Season: 1, Language: "en"}, path: "/shows/55/seasons/1/translations/en", body: list},
		{name: "seasons videos", handler: SeasonsVideosHandler{}, options: season, path: "/shows/55/seasons/1/videos", body: list},
		{name: "seasons watching", handler: SeasonsWatchingHandler{}, options: season, path: "/shows/55/seasons/1/watching", body: list},
		{name: "shows aliases", handler: ShowsAliasesHandler{}, options: byID, path: "/shows/55/aliases", body: list},
		{name: "shows certifications", handler: ShowsCertificationsHandler{}, options: byID, path: "/shows/55/certifications", body: list},
		{name: "shows people", handler: ShowsPeopleHandler{}, options: byID, path: "/shows/55/people", body: object},
		{name: "shows ratings", handler: ShowsRatingsHandler{}, options: byID, path: "/shows/55/ratings", body: object},
		{name: "shows stats", handler: ShowsStatsHandler{}, options: byID, path: "/shows/55/stats", body: object},
		{name: "shows studios", handler: ShowsStudiosHandler{}, options: byID, path: "/shows/55/studios", body: list},
		{name: "shows translations", handler: ShowsTranslationsHandler{}, options: str.Options{InternalID: "55", Language: "en"}, path: "/shows/55/translations/en", body: list},
		{name: "shows videos", handler: ShowsVideosHandler{}, options: byID, path: "/shows/55/videos", body: list},
		{name: "shows watching", handler: ShowsWatchingHandler{}, options: byID, path: "/shows/55/watching", body: list},
		{name: "sync last activities", handler: SyncLastActivitiesHandler{}, path: "/sync/last_activities", body: object},
		{name: "users blocked users", handler: UsersBlockedUsersHandler{}, path: "/users/blocked", body: list},
		{name: "users collaborations", handler: UsersCollaborationsHandler{}, options: user, path: "/users/me/lists/collaborations", body: list},
		{name: "users follower requests", handler: UsersFollowerRequestsHandler{}, path: "/users/requests", body: list},
		{name: "users followers", handler: UsersFollowersHandler{}, options: user, path: "/users/me/followers", body: list},
		{name: "users following", handler: UsersFollowingHandler{}, options: user, path: "/users/me/following", body: list},
		{name: "users following requests", handler: UsersFollowingRequestsHandler{}, path: "/users/requests/following", body: list},
		{name: "users friends", handler: UsersFriendsHandler{}, options: user, path: "/users/me/friends", body: list},
		{name: "users hidden items", handler: UsersHiddenItemsHandler{}, path: "/users/hidden/", body: list},
		{name: "users profile", handler: UsersProfileHandler{}, options: user, path: "/users/me", body: object},
		{name: "users saved filters", handler: UsersSavedFiltersHandler{}, options: typed, path: "/users/saved_filters/movies", body: list},
		{name: "users settings", handler: UsersSettingsHandler{}, path: "/users/settings", body: object},
		{name: "users stats", handler: UsersStatsHandler{}, options: user, path: "/users/me/stats", body: object},
	}
}

func TestSingleRequestHandlersWriteResult(t *testing.T) {
	for _, tc := range singleRequestHandlers() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			paths := []string{}
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				paths = append(paths, r.URL.Path)
				test.SafeFprint(w, tc.body)
			})

			options := tc.options
			options.Output = filepath.Join(t.TempDir(), "out.json")
			test.AssertNilError(t, tc.handler.Handle(&options, s.Client))
			assert.Equal(t, []string{tc.path}, paths)

			data, err := os.ReadFile(options.Output)
			test.AssertNilError(t, err)
			var written any
			test.AssertNilError(t, json.Unmarshal(data, &written))
			if tc.body == `{}` {
				assert.IsType(t, map[string]any{}, written)
				return
			}
			assert.Len(t, written, 1, "items written to -o")
		})
	}
}

func TestSingleRequestHandlersFailedRequest(t *testing.T) {
	for _, tc := range singleRequestHandlers() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			options := tc.options
			options.Output = filepath.Join(t.TempDir(), "out.json")
			assert.Error(t, tc.handler.Handle(&options, s.Client))
			_, err := os.Stat(options.Output)
			assert.True(t, os.IsNotExist(err), "nothing is written when the request fails")
		})
	}
}
