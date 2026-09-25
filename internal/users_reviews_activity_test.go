// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestUsersServiceGetCommentReactions(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/users/reactions/comments", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		if got, want := r.URL.RawQuery, "extended=min&limit=10&page=1"; got != want {
			t.Errorf("query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `[{"reacted_at":"2026-09-01T10:00:00.000Z","reaction":{"type":"love"},"type":"comment","comment":{"id":42}}]`)
	})

	got, _, err := setup.Client.Users.GetCommentReactions(context.Background(), &uri.ListOptions{Page: 1, Limit: 10, Extended: "min"})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.CommentReaction{{
		ReactedAt: &str.Timestamp{Time: time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)},
		Reaction:  &str.Reaction{Type: str.String("love")},
		Type:      str.String("comment"),
		Comment:   &str.Comment{ID: test.Ptr(42)},
	}}, got)
}

func TestUsersServiceGetSocialActivity(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/users/sean/following/activities", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		if got, want := r.URL.RawQuery, "extended=full&genres=drama&limit=10&page=2&years=2020-2026"; got != want {
			t.Errorf("query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `[{"id":9001,"activity_at":"2026-09-02T20:00:00.000Z","action":"watch","user":{"username":"justin"},"user_rating":8,`+
			`"type":"episode","episode":{"season":1,"number":2},"show":{"title":"Andor"}},`+
			`{"id":9002,"activity_at":"2026-09-02T21:00:00.000Z","action":"checkin","user":{"username":"sean"},"type":"movie","movie":{"title":"Arrival"}}]`)
	})

	opts := &uri.SocialActivityOptions{Page: 2, Limit: 10, Extended: "full", Genres: "drama", Years: "2020-2026"}
	got, _, err := setup.Client.Users.GetSocialActivity(context.Background(), str.String("sean"), str.String("following"), opts)
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.SocialActivity{
		{ID: test.Ptr(int64(9001)), ActivityAt: &str.Timestamp{Time: time.Date(2026, 9, 2, 20, 0, 0, 0, time.UTC)}, Action: str.String("watch"),
			User: &str.UserProfile{Username: str.String("justin")}, UserRating: test.Ptr(8), Type: str.String("episode"),
			Episode: &str.Episode{Season: test.Ptr(1), Number: test.Ptr(2)}, Show: &str.Show{Title: str.String("Andor")}},
		{ID: test.Ptr(int64(9002)), ActivityAt: &str.Timestamp{Time: time.Date(2026, 9, 2, 21, 0, 0, 0, time.UTC)}, Action: str.String("checkin"),
			User: &str.UserProfile{Username: str.String("sean")}, Type: str.String("movie"), Movie: &str.Movie{Title: str.String("Arrival")}},
	}, got)
}

func TestUsersServiceGetReviews(t *testing.T) {
	body := `{"stats":{"all":{"minutes":{"total":1200,"yearly":0,"monthly":1200,"weekly":280.5,"daily":40},"lists_counts":{"total":2,"yearly":0,"monthly":2,"weekly":0.5,"daily":0}},` +
		`"movies":{"play_counts":{"total":3,"yearly":0,"monthly":3,"weekly":0.7,"daily":0.1}}},` +
		`"images":{"cover":"cover.jpg","story":"story.jpg"},` +
		`"first_watched":{"watched_at":"2026-08-01T12:00:00.000Z","type":"movie","movie":{"title":"Arrival"}},"last_watched":null,` +
		`"countries":{"shows":{"country_count":1,"countries":[{"country":"us","count":4}]},"movies":{"country_count":0,"countries":[]}},` +
		`"trends":{"shows":[{"month":8,"month_name":"August","watchers":120,"watched":true,"show":{"title":"Andor"}}],"movies":[]},` +
		`"thanks":{"shows":[{"show":{"title":"Severance"}}],"movies":[]},` +
		`"streaming_services":{"country":"us","services":[{"source":"netflix","name":"Netflix","shows":2,"movies":1,"all":3}]}}`
	want := &str.Review{
		Stats: &str.ReviewStats{
			All: &str.ReviewStatsCategories{
				Minutes:     &str.ReviewStat{Total: test.Ptr(1200.0), Yearly: test.Ptr(0.0), Monthly: test.Ptr(1200.0), Weekly: test.Ptr(280.5), Daily: test.Ptr(40.0)},
				ListsCounts: &str.ReviewStat{Total: test.Ptr(2.0), Yearly: test.Ptr(0.0), Monthly: test.Ptr(2.0), Weekly: test.Ptr(0.5), Daily: test.Ptr(0.0)},
			},
			Movies: &str.ReviewStatsCategories{PlayCounts: &str.ReviewStat{Total: test.Ptr(3.0), Yearly: test.Ptr(0.0), Monthly: test.Ptr(3.0), Weekly: test.Ptr(0.7), Daily: test.Ptr(0.1)}},
		},
		Images: &str.ReviewImages{Cover: str.String("cover.jpg"), Story: str.String("story.jpg")},
		FirstWatched: &str.ReviewWatchedItem{WatchedAt: &str.Timestamp{Time: time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)}, Type: str.String("movie"),
			Movie: &str.Movie{Title: str.String("Arrival")}},
		Countries: &str.ReviewCountries{
			Shows:  &str.ReviewCountryCount{CountryCount: test.Ptr(1), Countries: []*str.ReviewCountry{{Country: str.String("us"), Count: test.Ptr(4)}}},
			Movies: &str.ReviewCountryCount{CountryCount: test.Ptr(0), Countries: []*str.ReviewCountry{}},
		},
		Trends: &str.ReviewTrends{
			Shows:  []*str.ReviewTrend{{Month: test.Ptr(8), MonthName: str.String("August"), Watchers: test.Ptr(120), Watched: test.Ptr(true), Show: &str.Show{Title: str.String("Andor")}}},
			Movies: []*str.ReviewTrend{},
		},
		Thanks: &str.ReviewThanks{Shows: []*str.MediaItem{{Show: &str.Show{Title: str.String("Severance")}}}, Movies: []*str.MediaItem{}},
		StreamingServices: &str.ReviewStreamingServices{Country: str.String("us"), Services: []*str.ReviewStreamingService{
			{Source: str.String("netflix"), Name: str.String("Netflix"), Shows: test.Ptr(2), Movies: test.Ptr(1), All: test.Ptr(3)},
		}},
	}

	tests := []struct {
		name string
		path string
		call func(u *UsersService) (*str.Review, *str.Response, error)
	}{
		{name: "month", path: "/users/sean/mir/2026/8", call: func(u *UsersService) (*str.Review, *str.Response, error) {
			return u.GetMonthInReview(context.Background(), str.String("sean"), 2026, 8, &uri.ListOptions{Extended: "images"})
		}},
		{name: "year", path: "/users/sean/yir/2026", call: func(u *UsersService) (*str.Review, *str.Response, error) {
			return u.GetYearInReview(context.Background(), str.String("sean"), 2026, &uri.ListOptions{Extended: "images"})
		}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				if got := r.URL.RawQuery; got != "extended=images" {
					t.Errorf("query is %q, want %q", got, "extended=images")
				}
				test.SafeFprint(w, body)
			})

			got, _, err := tt.call(setup.Client.Users)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, want, got)
		})
	}
}
