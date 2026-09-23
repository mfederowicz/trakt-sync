// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestHandlersErrorStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		path    string
		body    string
		handler Handler
		options str.Options
		wantErr string
	}{
		{
			name:    "users follow 404",
			status:  http.StatusNotFound,
			path:    "/users/sean/follow",
			body:    `{}`,
			handler: UsersFollowHandler{},
			options: str.Options{UserName: "sean"},
			wantErr: "user not found:sean",
		},
		{
			name:    "users block 404",
			status:  http.StatusNotFound,
			path:    "/users/sean/block",
			body:    `{}`,
			handler: UsersBlockHandler{},
			options: str.Options{UserName: "sean"},
			wantErr: "user not found:sean",
		},
		{
			name:    "users follow",
			status:  http.StatusConflict,
			path:    "/users/sean/follow",
			body:    `{}`,
			handler: UsersFollowHandler{},
			options: str.Options{UserName: "sean"},
			wantErr: consts.UserPendingFollowRequest,
		},
		{
			name:    "users block",
			status:  http.StatusConflict,
			path:    "/users/sean/block",
			body:    `{}`,
			handler: UsersBlockHandler{},
			options: str.Options{UserName: "sean"},
			wantErr: consts.UserBlockedAlready,
		},
		{
			name:    "movie checkin",
			status:  http.StatusConflict,
			path:    "/checkin",
			body:    `{"expires_at":"2026-09-24T20:00:00.000Z"}`,
			handler: CheckinMovieHandler{},
			options: str.Options{TraktID: 1, Action: consts.Movie, InternalID: "tron-legacy-2010"},
			wantErr: "409",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			MuxUserSettings(t, s.Mux)
			s.Mux.HandleFunc("/movies/tron-legacy-2010", func(w http.ResponseWriter, _ *http.Request) {
				test.SafeFprint(w, `{"title":"TRON: Legacy","ids":{"trakt":1}}`)
			})

			s.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodPost)
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			options := tt.options
			err := tt.handler.Handle(&options, s.Client)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
