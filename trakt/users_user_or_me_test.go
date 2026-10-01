// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestUserOrMe(t *testing.T) {
	if got := userOrMe(""); got != "me" {
		t.Errorf("userOrMe(\"\") = %q, want %q", got, "me")
	}
	if got := userOrMe("sean"); got != "sean" {
		t.Errorf("userOrMe(\"sean\") = %q, want %q", got, "sean")
	}
}

// TestUsersServiceEmptyUserIsMe checks that an empty user becomes "me" in the path,
// for read and write routes.
func TestUsersServiceEmptyUserIsMe(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{}
	tests := []struct {
		name   string
		method string
		path   string
		status int
		body   string
		call   func(u *UsersService) error
	}{
		{name: "GetHistory", method: http.MethodGet, path: "/users/me/history/movies", body: `[]`, call: func(u *UsersService) error {
			_, _, err := u.GetHistory(ctx, "", "movies", 0, opts)
			return err
		}},
		{name: "GetWatchlist", method: http.MethodGet, path: "/users/me/watchlist/movies/rank/asc", body: `[]`, call: func(u *UsersService) error {
			_, _, err := u.GetWatchlist(ctx, "", "movies", "rank", "asc", opts)
			return err
		}},
		{name: "GetListItems", method: http.MethodGet, path: "/users/me/lists/55/items/movie/rank/asc", body: `[]`, call: func(u *UsersService) error {
			_, _, err := u.GetListItems(ctx, "", "55", "movie", "rank", "asc", opts)
			return err
		}},
		{name: "GetFollowers", method: http.MethodGet, path: "/users/me/followers", body: `[]`, call: func(u *UsersService) error {
			_, _, err := u.GetFollowers(ctx, "", opts)
			return err
		}},
		{name: "GetWatching", method: http.MethodGet, path: "/users/me/watching", body: `{}`, call: func(u *UsersService) error {
			_, _, err := u.GetWatching(ctx, "", opts)
			return err
		}},
		{name: "GetSmartLists", method: http.MethodGet, path: "/users/me/smart-lists", body: `[]`, call: func(u *UsersService) error {
			_, _, err := u.GetSmartLists(ctx, "")
			return err
		}},
		{name: "GetMonthInReview", method: http.MethodGet, path: "/users/me/mir/2026/9", body: `{}`, call: func(u *UsersService) error {
			_, _, err := u.GetMonthInReview(ctx, "", 2026, 9, opts)
			return err
		}},
		{name: "AddListItems", method: http.MethodPost, path: "/users/me/lists/55/items", status: http.StatusCreated, body: `{}`, call: func(u *UsersService) error {
			_, _, err := u.AddListItems(ctx, "", "55", new(str.HistoryItems))
			return err
		}},
		{name: "UpdateList", method: http.MethodPut, path: "/users/me/lists/55", body: `{}`, call: func(u *UsersService) error {
			_, _, err := u.UpdateList(ctx, "", "55", new(str.PersonalList))
			return err
		}},
		{name: "LikeList", method: http.MethodPost, path: "/users/me/lists/55/like", status: http.StatusNoContent, call: func(u *UsersService) error {
			_, err := u.LikeList(ctx, "", "55")
			return err
		}},
		{name: "DeleteList", method: http.MethodDelete, path: "/users/me/lists/55", status: http.StatusNoContent, call: func(u *UsersService) error {
			_, err := u.DeleteList(ctx, "", "55")
			return err
		}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			calls := 0
			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, tt.method)
				if r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				if tt.status != 0 {
					w.WriteHeader(tt.status)
				}
				if tt.body != "" {
					test.SafeFprint(w, tt.body)
				}
			})

			test.AssertNilError(t, tt.call(setup.Client.Users))
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
		})
	}
}
