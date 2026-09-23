// Package internal used for client and services
package internal

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestUsersServiceListLike(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(s *UsersService, user *string, listID *string) (*str.Response, error)
	}{
		{
			name:   "like a list",
			method: http.MethodPost,
			call: func(s *UsersService, user *string, listID *string) (*str.Response, error) {
				return s.ListLike(context.Background(), user, listID)
			},
		},
		{
			name:   "remove like on a list",
			method: http.MethodDelete,
			call: func(s *UsersService, user *string, listID *string) (*str.Response, error) {
				return s.RemoveListLike(context.Background(), user, listID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc("/users/sean/lists/star-wars/like", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, tt.method)
				w.WriteHeader(http.StatusNoContent)
			})

			resp, err := tt.call(setup.Client.Users, str.String("sean"), str.String("star-wars"))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got, want := resp.StatusCode, http.StatusNoContent; got != want {
				t.Errorf("status code is %d, want %d", got, want)
			}
		})
	}
}

func TestUsersServiceFollowRequests(t *testing.T) {
	tests := []struct {
		name string
		path string
		call func(s *UsersService, opts *uri.ListOptions) ([]*str.FollowRequest, *str.Response, error)
	}{
		{
			name: "follow requests",
			path: "/users/requests",
			call: func(s *UsersService, opts *uri.ListOptions) ([]*str.FollowRequest, *str.Response, error) {
				return s.GetFollowRequests(context.Background(), opts)
			},
		},
		{
			name: "pending following requests",
			path: "/users/requests/following",
			call: func(s *UsersService, opts *uri.ListOptions) ([]*str.FollowRequest, *str.Response, error) {
				return s.GetPendingFollowingRequests(context.Background(), opts)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				test.AssertFormValues(t, r, map[string]string{"extended": "full"})
				fmt.Fprint(w, `[{"id":3,"requested_at":"2014-09-22T06:48:14.000Z","user":{"username":"sean"}}]`)
			})

			list, _, err := tt.call(setup.Client.Users, &uri.ListOptions{Extended: "full"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got, want := len(list), 1; got != want {
				t.Fatalf("got %d requests, want %d", got, want)
			}
		})
	}
}

func TestUsersServiceGetUserProfile(t *testing.T) {
	tests := []struct {
		name string
		id   *string
		path string
	}{
		{name: "given user", id: str.String("sean"), path: "/users/sean"},
		{name: "authenticated user", id: nil, path: "/users/me"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				fmt.Fprint(w, `{"username":"sean","name":"Sean Rudford"}`)
			})

			profile, _, err := setup.Client.Users.GetUserProfile(context.Background(), tt.id)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got, want := *profile.Name, "Sean Rudford"; got != want {
				t.Errorf("name is %q, want %q", got, want)
			}
		})
	}
}
