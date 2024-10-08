package uri

import (
	"testing"
)

var (
	DefaultRange                     = RatingRange{Min: 10, Max: 45}
	ListOptionsBasic                 = ListOptions{Page: 1, Limit: 10, Extended: "full"}
	ListOptionsCommon                = ListOptions{Genres: []string{"action", "adventure", "comedy"}, Years: "2016", StudioIDs: []int{1, 2, 3}}
	ListOptionsRatings               = ListOptions{Ratings: RatingRange{Min: 10, Max: 45}}
	ListOptionsInvalidRatings        = ListOptions{Ratings: RatingRange{Min: 100, Max: 48}}
	ListOptionsVotes                 = ListOptions{Votes: VotesRange{Min: 10, Max: 45}}
	ListOptionsEpisodesFilters       = ListOptions{Certifications: []string{"pg-13", "pg-16"}, NetworkIDs: []int{1, 2, 45}, EpisodeTypes: []string{"standard", "series_premiere"}}
	ListOptionsTmdbRatingsFilters    = ListOptions{TmdbRatings: TmdbRatingRange{Min: 5.5, Max: 10.0}}
	ListOptionsShowsFilters          = ListOptions{Certifications: []string{"pg-13", "pg-16"}, NetworkIDs: []int{1, 2, 45}, Status: []string{"pilot", "ended"}}
	ListOptionsCertificationsFilters = ListOptions{Certifications: []string{"pg-13", "pg-16"}}
	ListOptionsTmdbVotes             = ListOptions{TmdbVotes: VotesRange{Min: 25, Max: 40}}
	ListOptionsImdbVotes             = ListOptions{ImdbRatings: RatingRange{Min: 3, Max: 6}, ImdbVotes: ImdbVotesRange{Min: 10, Max: 25}}
	ListOptionsRt                    = ListOptions{RtMeters: RatingRange{Min: 55, Max: 100}, RtUserMeters: RatingRange{Min: 65, Max: 100}}
	ListOptionsMetascores            = ListOptions{Metascores: RatingRangeFloat{Min: 55, Max: 100}}
)

const (
	BaseURL  = "http://example.com"
	Expected = "Expected %q, got %q"
)

func TestBuildQueryBasic(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?extended=full&limit=10&page=1"

	list := ListOptionsBasic

	got, _ := AddQuery(BaseURL, list)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func testBuildQueryCommonFilters(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?years=2016&genres=action,adventure,comedy&studio_ids=1,2,3"
	got, _ := AddQuery(BaseURL, ListOptionsCommon)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func TestBuildQueryRatingFiltersRatings(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?ratings=10-45"
	got, _ := AddQuery(BaseURL, ListOptionsRatings)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func TestBuildQueryRatingFiltersInvalidRatings(t *testing.T) {
	t.Helper()

	got, _ := AddQuery(BaseURL, ListOptionsInvalidRatings)
	if string(got) != BaseURL {
		t.Fatalf(Expected, BaseURL, string(got))
	}
}

func TestBuildQueryRatingFiltersVotes(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?votes=10-45"

	got, _ := AddQuery(BaseURL, ListOptionsVotes)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func TestBuildQueryRatingFiltersTmdbRatings(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?tmdb_ratings=5.5-10.0"

	got, _ := AddQuery(BaseURL, ListOptionsTmdbRatingsFilters)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func TestBuildQueryRatingFiltersTmdbVotes(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?tmdb_votes=25-40"

	got, _ := AddQuery(BaseURL, ListOptionsTmdbVotes)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func TestBuildQueryRatingFiltersImdb(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?imdb_ratings=3-6&imdb_votes=10-25"

	got, _ := AddQuery(BaseURL, ListOptionsImdbVotes)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}
func TestBuildQueryRatingFiltersRt(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?rt_meters=55-100&rt_user_meters=65-100"
	got, _ := AddQuery(BaseURL, ListOptionsRt)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func TestBuildQueryRatingFiltersMeta(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?metascores=55.0-100.0"

	got, _ := AddQuery(BaseURL, ListOptionsMetascores)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func TestBuildQueryCertificationsFilters(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?certifications=pg-13,pg-16"

	got, _ := AddQuery(BaseURL, ListOptionsCertificationsFilters)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func TestBuildQueryShowFilters(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?certifications=pg-13,pg-16&network_ids=1,2,45&status=pilot,ended"

	got, _ := AddQuery(BaseURL, ListOptionsShowsFilters)
	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}

func TestBuildQueryEpisodeFilters(t *testing.T) {
	t.Helper()

	expectedURL := BaseURL + "?certifications=pg-13,pg-16&episode_types=standard,series_premiere&network_ids=1,2,45"

	got, err := AddQuery(BaseURL, ListOptionsEpisodesFilters)
	if err != nil {
		t.Logf("error:%s", err)
	}

	if string(got) != expectedURL {
		t.Fatalf(Expected, expectedURL, string(got))
	}
}
