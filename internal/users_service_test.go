// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
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
