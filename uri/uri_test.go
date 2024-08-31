package uri

import (
	"testing"
)

func TestBuildQueryBasic(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?extended=full&limit=10&page=1"

	url := "http://example.com"
	list := ListOptions{}
	list.Page = 1
	list.Limit = 10
	list.Extended = "full"

	got, _ := AddQuery(url, list)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func testBuildQueryCommonFilters(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?years=2016&genres=action,adventure,comedy&studio_ids=1,2,3"

	url := "http://example.com"
	core := ListOptions{}

	core.Genres = []string{"action", "adventure", "comedy"}
	core.Years = "2016"
	core.StudioIDs = []int{1, 2, 3}

	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryRatingFiltersRatings(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?ratings=10-45"

	url := "http://example.com"
	core := ListOptions{}
	core.Ratings = RatingRange{Min: 10, Max: 45}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryRatingFiltersInvalidRatings(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com"

	url := "http://example.com"
	core := ListOptions{}
	core.Ratings = RatingRange{Min: 100, Max: 45}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryRatingFiltersVotes(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?votes=10-45"

	url := "http://example.com"
	core := ListOptions{}
	core.Votes = VotesRange{Min: 10, Max: 45}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryRatingFiltersTmdbRatings(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?tmdb_ratings=5.5-10.0"

	url := "http://example.com"
	core := ListOptions{}
	core.TmdbRatings = TmdbRatingRange{Min: 5.5, Max: 10.0}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryRatingFiltersTmdbVotes(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?tmdb_votes=25-40"

	url := "http://example.com"
	core := ListOptions{}
	core.TmdbVotes = VotesRange{Min: 25, Max: 40}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryRatingFiltersImdb(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?imdb_ratings=3-6&imdb_votes=10-25"

	url := "http://example.com"
	core := ListOptions{}
	core.ImdbRatings = RatingRange{Min: 3, Max: 6}
	core.ImdbVotes = ImdbVotesRange{Min: 10, Max: 25}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}
func TestBuildQueryRatingFiltersRt(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?rt_meters=55-100&rt_user_meters=65-100"

	url := "http://example.com"
	core := ListOptions{}

	core.RtMeters = RatingRange{Min: 55, Max: 100}
	core.RtUserMeters = RatingRange{Min: 65, Max: 100}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryRatingFiltersMeta(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?metascores=55.0-100.0"

	url := "http://example.com"
	core := ListOptions{}
	core.Metascores = RatingRangeFloat{Min: 55, Max: 100}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryMovieFilters(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?certifications=pg-13,pg-16"

	url := "http://example.com"
	core := ListOptions{}
	core.Certifications = []string{"pg-13", "pg-16"}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryShowFilters(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?certifications=pg-13,pg-16&network_ids=1,2,45&status=pilot,ended"

	url := "http://example.com"
	core := ListOptions{}

	core.Certifications = []string{"pg-13", "pg-16"}
	core.NetworkIDs = []int{1, 2, 45}
	core.Status = []string{"pilot", "ended"}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}

func TestBuildQueryEpisodeFilters(t *testing.T) {
	t.Helper()

	expectedURL := "http://example.com?certifications=pg-13,pg-16&episode_types=standard,series_premiere&network_ids=1,2,45"

	url := "http://example.com"
	core := ListOptions{}

	core.Certifications = []string{"pg-13", "pg-16"}
	core.NetworkIDs = []int{1, 2, 45}
	core.EpisodeTypes = []string{"standard", "series_premiere"}
	got, _ := AddQuery(url, core)
	if string(got) != expectedURL {
		t.Fatalf("Expected %q, got %q", expectedURL, string(got))
	}

}
