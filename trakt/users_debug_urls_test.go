// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

// TestUsersServiceURLsGoToDebugLogger checks that users service methods send their
// request URL to DebugLogger and print nothing to stdout.
func TestUsersServiceURLsGoToDebugLogger(t *testing.T) {
	ctx := context.Background()
	opts := &uri.ListOptions{}
	tests := []struct {
		name string
		path string
		call func(u *UsersService) error
	}{
		{name: "GetWatched", path: "users/sean/watched/movies", call: func(u *UsersService) error {
			_, _, err := u.GetWatched(ctx, "sean", "movies", opts)
			return err
		}},
		{name: "GetLikes", path: "users/sean/likes/lists", call: func(u *UsersService) error {
			_, _, err := u.GetLikes(ctx, "sean", "lists", opts)
			return err
		}},
		{name: "GetListLikes", path: "users/sean/lists/55/likes", call: func(u *UsersService) error {
			_, _, err := u.GetListLikes(ctx, "sean", "55", opts)
			return err
		}},
		{name: "GetCollection", path: "users/sean/collection/movies", call: func(u *UsersService) error {
			_, _, err := u.GetCollection(ctx, "sean", "movies", opts)
			return err
		}},
		{name: "GetComments", path: "users/sean/comments/all/all", call: func(u *UsersService) error {
			_, _, err := u.GetComments(ctx, "sean", "all", "all", opts)
			return err
		}},
		{name: "GetNotes", path: "users/sean/notes/all", call: func(u *UsersService) error {
			_, _, err := u.GetNotes(ctx, "sean", "all", opts)
			return err
		}},
		{name: "GetCollaborations", path: "users/sean/lists/collaborations", call: func(u *UsersService) error {
			_, _, err := u.GetCollaborations(ctx, "sean", opts)
			return err
		}},
		{name: "GetListItems", path: "users/sean/lists/55/items/movie/rank/asc", call: func(u *UsersService) error {
			_, _, err := u.GetListItems(ctx, "sean", "55", "movie", "rank", "asc", opts)
			return err
		}},
		{name: "GetListComments", path: "users/sean/lists/55/comments/newest", call: func(u *UsersService) error {
			_, _, err := u.GetListComments(ctx, "sean", "55", "newest", opts)
			return err
		}},
		{name: "GetBlockedUsers", path: "users/blocked", call: func(u *UsersService) error {
			_, _, err := u.GetBlockedUsers(ctx, opts)
			return err
		}},
		{name: "GetFollowers", path: "users/sean/followers", call: func(u *UsersService) error {
			_, _, err := u.GetFollowers(ctx, "sean", opts)
			return err
		}},
		{name: "GetFollowing", path: "users/sean/following", call: func(u *UsersService) error {
			_, _, err := u.GetFollowing(ctx, "sean", opts)
			return err
		}},
		{name: "GetFriends", path: "users/sean/friends", call: func(u *UsersService) error {
			_, _, err := u.GetFriends(ctx, "sean", opts)
			return err
		}},
		{name: "GetHistory", path: "users/sean/history/movies", call: func(u *UsersService) error {
			_, _, err := u.GetHistory(ctx, "sean", "movies", 0, opts)
			return err
		}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				test.SafeFprint(w, `[]`)
			})
			logged := []string{}
			setup.Client.DebugLogger = func(v ...any) {
				logged = append(logged, fmt.Sprint(v...))
			}

			var err error
			stdout := captureStdout(t, func() { err = tt.call(setup.Client.Users) })
			test.AssertNilError(t, err)
			if stdout != "" {
				t.Errorf("stdout is %q, want nothing", stdout)
			}
			if !strings.Contains(strings.Join(logged, "\n"), "url:"+tt.path) {
				t.Errorf("DebugLogger got %q, want a line with %q", logged, "url:"+tt.path)
			}
		})
	}
}

// captureStdout returns what fn writes to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return string(out)
}
