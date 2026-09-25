// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersReviewsAndActivityHandlers(t *testing.T) {
	review := `{"stats":{"all":{"minutes":{"total":10}}},"images":{"cover":"c.jpg","story":"s.jpg"}}`
	tests := []struct {
		name      string
		handler   Handler
		options   str.Options
		path      string
		query     string
		status    int
		body      string
		wantErr   string
		wantNoAPI bool
	}{
		{name: "comment reactions", handler: UsersCommentReactionsHandler{}, options: str.Options{Action: consts.CommentReactions, PerPage: 10},
			path: "/users/reactions/comments", query: "limit=10&page=1", status: http.StatusOK, body: `[{"type":"comment","comment":{"id":1}}]`},
		{name: "comment reactions empty", handler: UsersCommentReactionsHandler{}, options: str.Options{Action: consts.CommentReactions},
			path: "/users/reactions/comments", query: "page=1", status: http.StatusOK, body: `[]`, wantErr: consts.EmptyResult},
		{name: "activities", handler: UsersActivitiesHandler{}, options: str.Options{Action: consts.Activities, Type: "friends", PerPage: 10, Genres: "drama"},
			path: "/users/sean/friends/activities", query: "genres=drama&limit=10&page=1", status: http.StatusOK, body: `[{"id":1,"type":"movie","movie":{"title":"Arrival"}}]`},
		{name: "activities without type", handler: UsersActivitiesHandler{}, options: str.Options{Action: consts.Activities},
			wantErr: "set -t to one of [friends followers following]", wantNoAPI: true},
		{name: "activities invalid type", handler: UsersActivitiesHandler{}, options: str.Options{Action: consts.Activities, Type: "movies"},
			wantErr: "set -t to one of", wantNoAPI: true},
		{name: "month in review", handler: UsersMonthInReviewHandler{}, options: str.Options{Action: consts.MonthInReview, Year: 2026, Month: 8},
			path: "/users/sean/mir/2026/8", status: http.StatusOK, body: review},
		{name: "month in review empty object", handler: UsersMonthInReviewHandler{}, options: str.Options{Action: consts.MonthInReview, Year: 2026, Month: 8},
			path: "/users/sean/mir/2026/8", status: http.StatusOK, body: `{}`, wantErr: "no month_in_review for 2026-08 (user sean)"},
		{name: "month in review not found", handler: UsersMonthInReviewHandler{}, options: str.Options{Action: consts.MonthInReview, Year: 1999, Month: 1},
			path: "/users/sean/mir/1999/1", status: http.StatusNotFound, body: `{}`, wantErr: "no month_in_review for 1999-01"},
		{name: "month in review without year", handler: UsersMonthInReviewHandler{}, options: str.Options{Action: consts.MonthInReview, Month: 8},
			wantErr: "set -year for month_in_review", wantNoAPI: true},
		{name: "month in review month 13", handler: UsersMonthInReviewHandler{}, options: str.Options{Action: consts.MonthInReview, Year: 2026, Month: 13},
			wantErr: "set -month 1-12", wantNoAPI: true},
		{name: "month in review without month", handler: UsersMonthInReviewHandler{}, options: str.Options{Action: consts.MonthInReview, Year: 2026},
			wantErr: "set -month 1-12", wantNoAPI: true},
		{name: "year in review", handler: UsersYearInReviewHandler{}, options: str.Options{Action: consts.YearInReview, Year: 2025, ExtendedInfo: "full"},
			path: "/users/sean/yir/2025", query: "extended=full", status: http.StatusOK, body: review},
		{name: "year in review without year", handler: UsersYearInReviewHandler{}, options: str.Options{Action: consts.YearInReview},
			wantErr: "set -year for year_in_review", wantNoAPI: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				if r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				if r.URL.RawQuery != tt.query {
					t.Errorf("query is %q, want %q", r.URL.RawQuery, tt.query)
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			options := tt.options
			options.UserName = "sean"
			options.Output = filepath.Join(t.TempDir(), "out.json")
			err := tt.handler.Handle(&options, s.Client)
			if tt.wantNoAPI && calls != 0 {
				t.Error("API was called with invalid options")
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				if _, statErr := os.Stat(options.Output); !os.IsNotExist(statErr) {
					t.Errorf("output file written on error: %v", statErr)
				}
				return
			}
			test.AssertNilError(t, err)
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
			if _, statErr := os.Stat(options.Output); statErr != nil {
				t.Errorf("output file was not written: %v", statErr)
			}
		})
	}
}
