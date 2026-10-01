// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/stretchr/testify/assert"
)

type fakeHandler struct{}

func (fakeHandler) Handle(*str.Options, *trakt.Client) error {
	return nil
}

func stamp(day int) *str.Timestamp {
	return &str.Timestamp{Time: time.Date(2026, time.October, day, 12, 0, 0, 0, time.UTC)}
}

func ids(trakt int64) *str.IDs {
	return &str.IDs{Trakt: Ptr(trakt)}
}

func TestCheckDates(t *testing.T) {
	c := &CommonLogic{}
	const (
		past   = "2020-01-01T00:00:00Z"
		later  = "2020-06-01T00:00:00Z"
		future = "2999-01-01T00:00:00Z"
	)
	cases := []struct {
		name    string
		from    string
		to      string
		tz      string
		wantErr string
	}{
		{name: "no dates", tz: "UTC"},
		{name: "from only", from: past, tz: "UTC"},
		{name: "valid range", from: past, to: later, tz: "Europe/Warsaw"},
		{name: "same from and to", from: past, to: past, tz: "UTC"},
		{name: "unknown timezone", tz: "Mars/Olympus", wantErr: "invalid timezone"},
		{name: "bad from", from: "yesterday", tz: "UTC", wantErr: "invalid from date"},
		{name: "bad to", to: "2020-01-01", tz: "UTC", wantErr: "invalid to date"},
		{name: "from in the future", from: future, tz: "UTC", wantErr: "'from' date must not be later than the current full hour"},
		{name: "to in the future", to: future, tz: "UTC", wantErr: "'to' date must not be later than:"},
		{name: "from after to", from: later, to: past, tz: "UTC", wantErr: "'from' date must be earlier than or equal to 'to' date"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := c.CheckDates(tc.from, tc.to, tc.tz)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestCheckOptionsAgainstModuleConfig(t *testing.T) {
	c := &CommonLogic{}
	checks := map[string]func(*str.Options) error{
		"CheckSortAndTypes":    c.CheckSortAndTypes,
		"CheckCommentsFilters": c.CheckCommentsFilters,
		"CheckIncludeReplies":  c.CheckIncludeReplies,
		"CheckCommentTypes":    c.CheckCommentTypes,
		"CheckTypes":           c.CheckTypes,
		"ValidReason":          c.ValidReason,
		"ValidPrivacy":         c.ValidPrivacy,
	}
	for name, check := range checks {
		t.Run(name+" unknown module", func(t *testing.T) {
			err := check(&str.Options{Module: "nope", Action: "trending"})
			assert.EqualError(t, err, "not found config for module 'nope'")
		})
		t.Run(name+" nothing set", func(t *testing.T) {
			assert.NoError(t, check(&str.Options{Module: "comments", Action: "trending"}))
		})
	}

	cases := []struct {
		name    string
		check   func(*str.Options) error
		options str.Options
		wantErr string
	}{
		{name: "sort and type from the config", check: c.CheckSortAndTypes,
			options: str.Options{Module: "comments", Action: "trending", Type: "movies", Sort: "newest"}},
		{name: "unknown type", check: c.CheckSortAndTypes,
			options: str.Options{Module: "comments", Action: "trending", Type: "books", Sort: "newest"}, wantErr: "not found type for module 'comments'"},
		{name: "unknown sort", check: c.CheckSortAndTypes,
			options: str.Options{Module: "comments", Action: "trending", Type: "movies", Sort: "random"}, wantErr: "not found sort for module 'comments'"},
		{name: "CheckTypes unknown type", check: c.CheckTypes,
			options: str.Options{Module: "comments", Action: "trending", Type: "books"}, wantErr: "not found type for module 'comments'"},
		{name: "unknown comment type", check: c.CheckCommentTypes,
			options: str.Options{Module: "comments", Action: "trending", CommentType: "poems"}, wantErr: "not found comment_type for module 'comments'"},
		{name: "unknown include replies", check: c.CheckIncludeReplies,
			options: str.Options{Module: "comments", Action: "trending", IncludeReplies: "maybe"}, wantErr: "not found include_replies for module 'comments'"},
		{name: "filters stop at comment type", check: c.CheckCommentsFilters,
			options: str.Options{Module: "comments", Action: "trending", CommentType: "poems", Type: "books"}, wantErr: "not found comment_type for module 'comments'"},
		{name: "filters stop at type", check: c.CheckCommentsFilters,
			options: str.Options{Module: "comments", Action: "trending", Type: "books", IncludeReplies: "maybe"}, wantErr: "not found type for module 'comments'"},
		{name: "filters stop at include replies", check: c.CheckCommentsFilters,
			options: str.Options{Module: "comments", Action: "trending", Type: "movies", IncludeReplies: "maybe"}, wantErr: "not found include_replies for module 'comments'"},
		{name: "reason from the config", check: c.ValidReason,
			options: str.Options{Module: "comments", Action: "report", Reason: "spam"}},
		{name: "unknown reason", check: c.ValidReason,
			options: str.Options{Module: "comments", Action: "report", Reason: "boring"}, wantErr: "reason 'boring' is not valid for module 'comments' and action 'report'"},
		{name: "reason is free where the action has none", check: c.ValidReason,
			options: str.Options{Module: "comments", Action: "trending", Reason: "boring"}},
		{name: "privacy from the config", check: c.ValidPrivacy,
			options: str.Options{Module: "notes", Action: "notes", Privacy: "friends"}},
		{name: "unknown privacy", check: c.ValidPrivacy,
			options: str.Options{Module: "notes", Action: "notes", Privacy: "family"}, wantErr: "invalid privacy 'family' for module 'notes'"},
		{name: "reaction from the config", check: c.ValidReaction,
			options: str.Options{Module: "comments", Action: "reaction", Reaction: "love"}},
		{name: "unknown reaction", check: c.ValidReaction,
			options: str.Options{Module: "comments", Action: "reaction", Reaction: "meh"}, wantErr: "reaction 'meh' is not valid for module 'comments' and action 'reaction'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			options := tc.options
			err := tc.check(&options)
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestGetHandlerForMap(t *testing.T) {
	c := &CommonLogic{}
	handlers := map[string]Handler{"trending": fakeHandler{}}

	got, err := c.GetHandlerForMap("trending", handlers)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, fakeHandler{}, got)

	got, err = c.GetHandlerForMap("popular", handlers)
	assert.EqualError(t, err, "unknown handler")
	assert.Nil(t, got)
}

func TestToTimestamp(t *testing.T) {
	c := &CommonLogic{}

	got := c.ToTimestamp("2026-10-01T16:07:14+02:00")
	assert.True(t, got.Equal(time.Date(2026, time.October, 1, 14, 7, 14, 0, time.UTC)))
	assert.Equal(t, time.UTC, got.Location())

	assert.True(t, c.ToTimestamp("2026-10-01").IsZero(), "a date without time is not accepted")
	assert.True(t, c.ToTimestamp("").IsZero())
}

func TestDateLastDays(t *testing.T) {
	c := &CommonLogic{}
	const days = 7
	for _, tz := range []string{"UTC", "Europe/Warsaw", "America/Los_Angeles"} {
		t.Run(tz, func(t *testing.T) {
			loc, err := time.LoadLocation(tz)
			if !assert.NoError(t, err) {
				return
			}
			before := time.Now().In(loc).AddDate(0, 0, -days).Truncate(time.Second)

			got, err := time.Parse(time.RFC3339, c.DateLastDays(days, tz, false))
			if !assert.NoError(t, err) {
				return
			}
			after := time.Now().In(loc).AddDate(0, 0, -days)
			assert.False(t, got.Before(before), "%v is before %v", got, before)
			assert.False(t, got.After(after), "%v is after %v", got, after)

			full, err := time.Parse(time.RFC3339, c.DateLastDays(days, tz, true))
			if !assert.NoError(t, err) {
				return
			}
			assert.Zero(t, full.Minute())
			assert.Zero(t, full.Second())
			assert.False(t, full.After(after), "%v is after %v", full, after)
			assert.Less(t, after.Sub(full), time.Hour+time.Minute)
		})
	}
}

func TestTypeHelpers(t *testing.T) {
	cases := []struct {
		name  string
		check func(string) bool
		yes   []string
	}{
		{name: "movie", check: isMovieType, yes: []string{"movie", "movies"}},
		{name: "show", check: isShowType, yes: []string{"show", "shows"}},
		{name: "season", check: isSeasonType, yes: []string{"season", "seasons"}},
		{name: "episode", check: isEpisodeType, yes: []string{"episode", "episodes"}},
		{name: "people", check: isPeopleType, yes: []string{"people"}},
	}
	all := []string{"movie", "movies", "show", "shows", "season", "seasons", "episode", "episodes", "people", "person", ""}
	for _, tc := range cases {
		for _, value := range all {
			want := false
			for _, yes := range tc.yes {
				if yes == value {
					want = true
				}
			}
			assert.Equal(t, want, tc.check(value), "%s type check for %q", tc.name, value)
		}
	}
}

func TestSortRouteType(t *testing.T) {
	cases := []struct {
		name    string
		section string
		options str.Options
		want    string
		wantErr string
	}{
		{name: "movies", section: consts.Watchlist, options: str.Options{Type: "movies", SortPath: "rank"}, want: "movies"},
		{name: "shows", section: consts.Favorites, options: str.Options{Type: "shows", SortPath: "added"}, want: "shows"},
		{name: "all on the watchlist", section: consts.Watchlist, options: str.Options{Type: "all", SortPath: "title"}, want: "movie,show"},
		{name: "all on favorites", section: consts.Favorites, options: str.Options{Type: "all", SortPath: "title"}, want: "media"},
		{name: "unknown sort", section: consts.Watchlist, options: str.Options{Type: "movies", SortPath: "random"}, wantErr: "sort 'random' is not valid"},
		{name: "type without a sorted route", section: consts.Watchlist, options: str.Options{Type: "seasons", SortPath: "rank"}, wantErr: "-sort works with -t all, movies or shows, not 'seasons'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			options := tc.options
			got, err := sortRouteType(tc.section, &options)
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
				assert.Empty(t, got)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestReportMessage(t *testing.T) {
	assert.Equal(t, "already reported", reportMessage(Ptr("already reported"), errors.New("409")))
	assert.Equal(t, "409", reportMessage(Ptr(""), errors.New("409")))
	assert.Equal(t, "409", reportMessage(nil, errors.New("409")))
	assert.Empty(t, reportMessage(nil, nil))
}

func TestSeasonsWithEpisodeNumbersOnly(t *testing.T) {
	assert.Nil(t, SeasonsWithEpisodeNumbersOnly(nil))

	src := &[]str.Season{
		{Number: Ptr(1), Episodes: &[]str.Episode{{Number: Ptr(1), Title: Ptr("Pilot"), WatchedAt: stamp(1)}, {Number: Ptr(2), Title: Ptr("Second")}}},
		{Number: Ptr(2)},
	}
	got := SeasonsWithEpisodeNumbersOnly(src)

	want := &[]str.Season{
		{Number: Ptr(1), Episodes: &[]str.Episode{{Number: Ptr(1)}, {Number: Ptr(2)}}},
		{Number: Ptr(2)},
	}
	assert.Equal(t, want, got)
	assert.Equal(t, "Pilot", *(*(*src)[0].Episodes)[0].Title, "the source list must stay untouched")
}

func TestOnlyIDs(t *testing.T) {
	c := &CommonLogic{}
	items := &[]str.ExportlistItem{
		{Title: Ptr("first"), IDs: ids(1), WatchedAt: stamp(1)},
		{Title: Ptr("second"), IDs: ids(2), Seasons: &[]str.Season{{Number: Ptr(1), Episodes: &[]str.Episode{{Number: Ptr(3), Title: Ptr("Third")}}}}},
		{Title: Ptr("third"), IDs: ids(3), Seasons: &[]str.Season{}},
	}

	assert.Equal(t, &[]str.Movie{{IDs: ids(1)}, {IDs: ids(2)}, {IDs: ids(3)}}, c.OnlyMoviesIDs(items))
	assert.Equal(t, &[]str.Season{{IDs: ids(1)}, {IDs: ids(2)}, {IDs: ids(3)}}, c.OnlySeasonsIDs(items))
	assert.Equal(t, &[]str.Episode{{IDs: ids(1)}, {IDs: ids(2)}, {IDs: ids(3)}}, c.OnlyEpisodesIDs(items))
	assert.Equal(t, &[]str.Show{
		{IDs: ids(1)},
		{IDs: ids(2), Seasons: &[]str.Season{{Number: Ptr(1), Episodes: &[]str.Episode{{Number: Ptr(3)}}}}},
		{IDs: ids(3)},
	}, c.OnlyShowsIDs(items))

	empty := c.OnlyMoviesIDs(&[]str.ExportlistItem{})
	assert.Empty(t, *empty)
}

func TestCreateItemsToRemove(t *testing.T) {
	c := &CommonLogic{}
	items := c.InitItemsList()
	*items.Movies = append(*items.Movies, str.ExportlistItem{Title: Ptr("TRON"), IDs: ids(1)})
	*items.Shows = append(*items.Shows, str.ExportlistItem{IDs: ids(2)})
	*items.Seasons = append(*items.Seasons, str.ExportlistItem{IDs: ids(3)})
	*items.Episodes = append(*items.Episodes, str.ExportlistItem{IDs: ids(4)})

	want := str.ItemsToRemove{
		Movies:   &[]str.Movie{{IDs: ids(1)}},
		Shows:    &[]str.Show{{IDs: ids(2)}},
		Seasons:  &[]str.Season{{IDs: ids(3)}},
		Episodes: &[]str.Episode{{IDs: ids(4)}},
	}
	assert.Equal(t, want, c.CreateItemsToRemove(items))
	assert.Equal(t, want, c.CreateItemsToRemoveRatings(items))
}

func TestCreateItemsToReorder(t *testing.T) {
	c := &CommonLogic{}
	items := c.InitItemsList()
	*items.Movies = append(*items.Movies, str.ExportlistItem{ID: Ptr(int64(11))})
	*items.Shows = append(*items.Shows, str.ExportlistItem{ID: Ptr(int64(12))})
	*items.Seasons = append(*items.Seasons, str.ExportlistItem{ID: Ptr(int64(13))})
	*items.Episodes = append(*items.Episodes, str.ExportlistItem{ID: Ptr(int64(14))})
	*items.Lists = append(*items.Lists, str.PersonalList{IDs: ids(15)})

	got := c.CreateItemsToReorder(items)
	assert.Equal(t, []int64{11, 12, 13, 14, 15}, *got.Rank)
}

func TestCreateItemsToAdd(t *testing.T) {
	c := &CommonLogic{}
	seasons := &[]str.Season{{Number: Ptr(1)}}
	items := c.InitItemsList()
	*items.Movies = append(*items.Movies, str.ExportlistItem{Title: Ptr("TRON"), Year: Ptr(2010), IDs: ids(1), WatchedAt: stamp(1), HiddenAt: stamp(2), Notes: Ptr("imax"), Rating: Ptr(9)})
	*items.Shows = append(*items.Shows, str.ExportlistItem{Title: Ptr("Dark"), Year: Ptr(2017), IDs: ids(2), Seasons: seasons, Notes: Ptr("again"), HiddenAt: stamp(3)})
	*items.Seasons = append(*items.Seasons, str.ExportlistItem{IDs: ids(3), WatchedAt: stamp(4)})
	*items.Episodes = append(*items.Episodes, str.ExportlistItem{IDs: ids(4), WatchedAt: stamp(5)})
	*items.Users = append(*items.Users, str.ExportlistItem{IDs: ids(5), HiddenAt: stamp(6)})
	*items.People = append(*items.People, str.ExportlistItem{IDs: ids(6), HiddenAt: stamp(7)})

	want := str.HistoryItems{
		Movies:   &[]str.Movie{{Title: Ptr("TRON"), Year: Ptr(2010), IDs: ids(1), WatchedAt: stamp(1), HiddenAt: stamp(2), Notes: Ptr("imax")}},
		Shows:    &[]str.Show{{Title: Ptr("Dark"), Year: Ptr(2017), IDs: ids(2), Seasons: seasons, Notes: Ptr("again"), HiddenAt: stamp(3)}},
		Seasons:  &[]str.Season{{IDs: ids(3), WatchedAt: stamp(4)}},
		Episodes: &[]str.Episode{{IDs: ids(4), WatchedAt: stamp(5)}},
		Users:    &[]str.UserProfile{{IDs: ids(5), HiddenAt: stamp(6)}},
		People:   &[]str.Person{{IDs: ids(6), HiddenAt: stamp(7)}},
	}
	assert.Equal(t, want, c.CreateItemsToAdd(items))

	empty := c.CreateItemsToAdd(c.InitItemsList())
	assert.Empty(t, *empty.Movies)
	assert.Empty(t, *empty.People)
}

func TestCreateItemsToHidden(t *testing.T) {
	c := &CommonLogic{}
	items := c.InitItemsList()
	*items.Movies = append(*items.Movies, str.ExportlistItem{IDs: ids(1)})
	*items.Shows = append(*items.Shows, str.ExportlistItem{IDs: ids(2)})
	*items.Seasons = append(*items.Seasons, str.ExportlistItem{IDs: ids(3)})
	*items.Users = append(*items.Users, str.ExportlistItem{IDs: ids(5)})

	movies := &[]str.Movie{{IDs: ids(1)}}
	shows := &[]str.Show{{IDs: ids(2)}}
	seasons := &[]str.Season{{IDs: ids(3)}}
	users := &[]str.UserProfile{{IDs: ids(5)}}
	cases := []struct {
		section string
		want    str.HistoryItems
	}{
		{section: consts.Calendar, want: str.HistoryItems{Movies: movies, Shows: shows}},
		{section: consts.ProgressWatched, want: str.HistoryItems{Shows: shows, Seasons: seasons}},
		{section: consts.ProgressCollected, want: str.HistoryItems{Shows: shows, Seasons: seasons}},
		{section: consts.Recommendations, want: str.HistoryItems{Movies: movies, Shows: shows}},
		{section: consts.Comments, want: str.HistoryItems{Users: users}},
		{section: consts.Dropped, want: str.HistoryItems{Shows: shows}},
		{section: "unknown", want: str.HistoryItems{}},
	}
	for _, tc := range cases {
		t.Run(tc.section, func(t *testing.T) {
			assert.Equal(t, tc.want, c.CreateItemsToHidden(tc.section, items))
		})
	}
}

func TestCreateItemsToAddRatings(t *testing.T) {
	c := &CommonLogic{}
	seasons := &[]str.Season{{Number: Ptr(1)}}
	items := c.InitItemsList()
	*items.Movies = append(*items.Movies,
		str.ExportlistItem{Title: Ptr("TRON"), Year: Ptr(2010), IDs: ids(1), RatedAt: stamp(1), Rating: Ptr(9)},
		str.ExportlistItem{IDs: ids(7)},
	)
	*items.Shows = append(*items.Shows,
		str.ExportlistItem{Title: Ptr("Dark"), Year: Ptr(2017), IDs: ids(2), Seasons: seasons, RatedAt: stamp(2), Rating: Ptr(8)},
		str.ExportlistItem{IDs: ids(8), Seasons: &[]str.Season{}, RatedAt: stamp(9)},
	)
	*items.Seasons = append(*items.Seasons, str.ExportlistItem{IDs: ids(3), RatedAt: stamp(3), Rating: Ptr(7)})
	*items.Episodes = append(*items.Episodes, str.ExportlistItem{IDs: ids(4), RatedAt: stamp(4), Rating: Ptr(6)})

	want := str.RatingItems{
		Movies: &[]str.Movie{
			{Title: Ptr("TRON"), Year: Ptr(2010), IDs: ids(1), RatedAt: stamp(1), Rating: Ptr(float32(9))},
			{IDs: ids(7)},
		},
		Shows: &[]str.Show{
			{Title: Ptr("Dark"), Year: Ptr(2017), IDs: ids(2), Seasons: seasons, RatedAt: stamp(2), Rating: Ptr(float32(8))},
			{IDs: ids(8)},
		},
		Seasons:  &[]str.Season{{IDs: ids(3), RatedAt: stamp(3), Rating: Ptr(float32(7))}},
		Episodes: &[]str.Episode{{IDs: ids(4), RatedAt: stamp(4), Rating: Ptr(float32(6))}},
	}
	assert.Equal(t, want, c.CreateItemsToAddRatings(items))
}

func TestUpdateHistoryListWithType(t *testing.T) {
	c := &CommonLogic{}
	cases := map[string]string{
		"shows":    "show",
		"episodes": "episode",
		"movies":   "movie",
		"seasons":  "season",
		"other":    "movie",
	}
	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			strType := in
			data := []*str.ExportlistItem{{IDs: ids(1)}, {IDs: ids(2)}}
			got := c.UpdateHistoryListWithType(data, &strType)
			if !assert.Len(t, got, 2) {
				return
			}
			for _, item := range got {
				assert.Equal(t, want, *item.Type)
			}
			assert.Same(t, data[0], got[0])
		})
	}
}

// mixedList has one movie, one show episode, one season, one episode and one person, each with a list item id.
const mixedList = `[
	{"id":101,"watched_at":"2026-10-01T12:00:00Z","movie":{"title":"TRON: Legacy","year":2010,"ids":{"trakt":1}}},
	{"id":102,"watched_at":"2026-10-02T12:00:00Z","show":{"title":"Dark","year":2017,"ids":{"trakt":2}},"episode":{"season":1,"number":3,"ids":{"trakt":4}}},
	{"id":103,"watched_at":"2026-10-03T12:00:00Z","season":{"number":1,"ids":{"trakt":3}}},
	{"id":104,"collected_at":"2026-10-04T12:00:00Z","person":{"name":"Bryan Cranston","ids":{"trakt":5}}}
]`

func traktIDs(items *[]str.ExportlistItem) []int64 {
	out := []int64{}
	for _, item := range *items {
		out = append(out, *item.IDs.Trakt)
	}
	return out
}

func TestConvertBytesToItemsList(t *testing.T) {
	c := &CommonLogic{}

	t.Run("history of one type", func(t *testing.T) {
		got, err := c.ConvertBytesToItemsList([]byte(mixedList), consts.AddToHistory, consts.Movies)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []int64{1}, traktIDs(got.Movies))
		assert.Empty(t, *got.Shows)
		assert.Equal(t, []int64{101}, *got.IDs)
		assert.Equal(t, "TRON: Legacy", *(*got.Movies)[0].Movie.Title)
	})

	t.Run("history of all types", func(t *testing.T) {
		got, err := c.ConvertBytesToItemsList([]byte(mixedList), consts.RemoveFromHistory, consts.ActionTypeAll)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []int64{1}, traktIDs(got.Movies))
		assert.Equal(t, []int64{2}, traktIDs(got.Shows))
		assert.Equal(t, []int64{3}, traktIDs(got.Seasons))
		assert.Equal(t, []int64{4}, traktIDs(got.Episodes))

		show := (*got.Shows)[0]
		if !assert.Len(t, *show.Seasons, 1) {
			return
		}
		season := (*show.Seasons)[0]
		assert.Equal(t, 1, *season.Number)
		if !assert.Len(t, *season.Episodes, 1) {
			return
		}
		assert.Equal(t, 3, *(*season.Episodes)[0].Number)
	})

	t.Run("collection of one type", func(t *testing.T) {
		got, err := c.ConvertBytesToItemsList([]byte(mixedList), consts.AddToCollection, consts.Shows)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []int64{2}, traktIDs(got.Shows))
		assert.Empty(t, *got.Movies)
		assert.Equal(t, []int64{101, 102, 103, 104}, *got.IDs)
		show := (*got.Shows)[0]
		assert.Equal(t, "Dark", *show.Title)
		assert.Equal(t, int64(102), *show.ID)
		assert.True(t, show.WatchedAt.Equal(stamp(2).Time))
	})

	t.Run("collection of people", func(t *testing.T) {
		got, err := c.ConvertBytesToItemsList([]byte(mixedList), consts.AddListItems, consts.People)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []int64{5}, traktIDs(got.People))
		assert.True(t, (*got.People)[0].CollectedAt.Equal(stamp(4).Time))
	})

	t.Run("hidden items of one type", func(t *testing.T) {
		got, err := c.ConvertBytesToItemsList([]byte(mixedList), consts.AddHiddenItems, consts.Season)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []int64{3}, traktIDs(got.Seasons))
		assert.Empty(t, *got.Movies)
	})

	t.Run("hidden items without a type", func(t *testing.T) {
		got, err := c.ConvertBytesToItemsList([]byte(mixedList), consts.RemoveHiddenItems, consts.EmptyString)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []int64{1}, traktIDs(got.Movies))
		assert.Equal(t, []int64{2}, traktIDs(got.Shows))
		assert.Equal(t, []int64{3}, traktIDs(got.Seasons))
		assert.Equal(t, []int64{4}, traktIDs(got.Episodes))
	})

	t.Run("history and ratings with an unknown type", func(t *testing.T) {
		for _, action := range []string{consts.AddToHistory, consts.RemoveFromHistory, consts.AddToRatings, consts.RemoveFromRatings} {
			for _, stype := range []string{"dance", "movie", consts.EmptyString} {
				got, err := c.ConvertBytesToItemsList([]byte(mixedList), action, stype)
				assert.Nil(t, got, "%s -t %q", action, stype)
				assert.EqualError(t, err, "type '"+stype+"' is not valid for action '"+action+"', available types:[all movies shows seasons episodes]")
			}
		}
	})

	t.Run("unknown action", func(t *testing.T) {
		got, err := c.ConvertBytesToItemsList([]byte(mixedList), "dance", consts.Movies)
		assert.Nil(t, got)
		assert.EqualError(t, err, consts.UnknownItemsListType)
	})

	t.Run("bad json", func(t *testing.T) {
		got, err := c.ConvertBytesToItemsList([]byte(`{"not":"a list"}`), consts.AddToHistory, consts.Movies)
		assert.Nil(t, got)
		assert.Error(t, err)
	})
}

func TestConvertBytes(t *testing.T) {
	c := &CommonLogic{}

	t.Run("reorder lists", func(t *testing.T) {
		got, err := c.ConvertBytes([]byte(`[{"name":"Star Wars","ids":{"trakt":55}},{"name":"Vampires","ids":{"trakt":52}}]`), str.Options{Action: consts.ReorderLists})
		if !assert.NoError(t, err) {
			return
		}
		if !assert.Len(t, *got.Lists, 2) {
			return
		}
		assert.Equal(t, "Vampires", *(*got.Lists)[1].Name)
		assert.Empty(t, *got.Movies)
	})

	t.Run("reorder lists bad json", func(t *testing.T) {
		got, err := c.ConvertBytes([]byte(`{"name":"Star Wars"}`), str.Options{Action: consts.ReorderLists})
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("add list", func(t *testing.T) {
		got, err := c.ConvertBytes([]byte(`{"name":"Star Wars","privacy":"public"}`), str.Options{Action: consts.AddList})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "Star Wars", *got.List.Name)
		assert.Nil(t, got.Movies)
	})

	t.Run("add list bad json", func(t *testing.T) {
		got, err := c.ConvertBytes([]byte(`[]`), str.Options{Action: consts.AddList})
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("items by action and type", func(t *testing.T) {
		got, err := c.ConvertBytes([]byte(mixedList), str.Options{Action: consts.AddToWatchlist, Type: consts.Movies})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, []int64{1}, traktIDs(got.Movies))
	})
}

func TestReadInputFromFile(t *testing.T) {
	c := &CommonLogic{}
	path := filepath.Join(t.TempDir(), "items.json")
	if !assert.NoError(t, os.WriteFile(path, []byte(mixedList), 0o600)) {
		return
	}

	data, err := c.ReadInputBytes(str.Options{Items: path})
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, mixedList, string(data))

	items, err := c.ReadInput(str.Options{Items: path, Action: consts.AddToCollection, Type: consts.Movies})
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []int64{1}, traktIDs(items.Movies))

	missing := filepath.Join(t.TempDir(), "missing.json")
	_, err = c.ReadInputBytes(str.Options{Items: missing})
	assert.ErrorContains(t, err, "failed to read file "+missing)
	_, err = c.ReadInput(str.Options{Items: missing, Action: consts.AddToCollection})
	assert.ErrorContains(t, err, "failed to read file "+missing)
}
