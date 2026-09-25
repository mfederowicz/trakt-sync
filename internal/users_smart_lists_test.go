// Package internal used for client and services
package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersServiceSmartLists(t *testing.T) {
	write := &str.SmartListWrite{Name: str.String("Sci-Fi"), Source: str.String("popular"), MediaType: str.String("movies"),
		Filters: &str.SmartListFilters{Genres: []string{"science-fiction"}, Years: []int{2000, 2026}}, Privacy: str.String("public")}
	tests := []struct {
		name     string
		method   string
		path     string
		status   int
		body     string
		wantBody *str.SmartListWrite
		call     func(u *UsersService) (any, *str.Response, error)
		want     any
	}{
		{
			name: "list", method: http.MethodGet, path: "/users/sean/smart-lists", status: http.StatusOK,
			body: `[{"name":"Sci-Fi","ids":{"trakt":1,"slug":"sci-fi"}}]`,
			call: func(u *UsersService) (any, *str.Response, error) {
				return u.GetSmartLists(context.Background(), str.String("sean"))
			},
			want: []*str.SmartList{{Name: str.String("Sci-Fi"), IDs: &str.IDs{Trakt: test.Ptr(int64(1)), Slug: str.String("sci-fi")}}},
		},
		{
			name: "get", method: http.MethodGet, path: "/users/sean/smart-lists/sci-fi", status: http.StatusOK,
			body: `{"name":"Sci-Fi","source":"popular","media_type":"movies","filters":{"genres":["science-fiction"]}}`,
			call: func(u *UsersService) (any, *str.Response, error) {
				return u.GetSmartList(context.Background(), str.String("sean"), str.String("sci-fi"))
			},
			want: &str.SmartList{Name: str.String("Sci-Fi"), Source: str.String("popular"), MediaType: str.String("movies"),
				Filters: &str.SmartListFilters{Genres: []string{"science-fiction"}}},
		},
		{
			name: "create", method: http.MethodPost, path: "/users/me/smart-lists", status: http.StatusCreated,
			body: `{"ids":{"trakt":7,"slug":"sci-fi"}}`, wantBody: write,
			call: func(u *UsersService) (any, *str.Response, error) {
				return u.AddSmartList(context.Background(), str.String("me"), write)
			},
			want: &str.SmartList{IDs: &str.IDs{Trakt: test.Ptr(int64(7)), Slug: str.String("sci-fi")}},
		},
		{
			name: "update", method: http.MethodPut, path: "/users/me/smart-lists/sci-fi", status: http.StatusOK,
			body: `{"name":"Sci-Fi","privacy":"private"}`, wantBody: &str.SmartListWrite{Privacy: str.String("private")},
			call: func(u *UsersService) (any, *str.Response, error) {
				return u.UpdateSmartList(context.Background(), str.String("me"), str.String("sci-fi"), &str.SmartListWrite{Privacy: str.String("private")})
			},
			want: &str.SmartList{Name: str.String("Sci-Fi"), Privacy: str.String("private")},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, tt.method)
				if tt.wantBody != nil {
					got := new(str.SmartListWrite)
					test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
					test.AssertNoDiff(t, tt.wantBody, got)
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			got, _, err := tt.call(setup.Client.Users)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, tt.want, got)
		})
	}
}

func TestUsersServiceDeleteSmartList(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/users/me/smart-lists/sci-fi", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := setup.Client.Users.DeleteSmartList(context.Background(), str.String("me"), str.String("sci-fi"))
	test.AssertNilError(t, err)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		t.Errorf("response is %v, want status %d", resp, http.StatusNoContent)
	}
}
