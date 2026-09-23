// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

var (
	usersReportOptions     = str.Options{Module: "users", Action: "report", UserName: "sean", Reason: "spam", Msg: "spam"}
	usersListReportOptions = str.Options{Module: "users", Action: "list_report", UserName: "sean", ID: "1", Reason: "spam", Msg: "spam"}
)

// TestHandlersWithoutResponse checks that handlers return an error instead of
// panicking when the request fails before any HTTP response (resp is nil).
func TestHandlersWithoutResponse(t *testing.T) {
	tests := []struct {
		name    string
		handler Handler
		options str.Options
	}{
		{name: "checkin delete", handler: CheckinDeleteHandler{}},
		{name: "comments like", handler: CommentsLikeHandler{}, options: str.Options{CommentID: 1}},
		{name: "lists like", handler: ListsLikeHandler{}, options: str.Options{InternalID: "1"}},
		{name: "lists list", handler: ListsListHandler{}, options: str.Options{InternalID: "1"}},
		{name: "movies refresh", handler: MoviesRefreshHandler{}, options: str.Options{InternalID: "tron-legacy-2010"}},
		{name: "people refresh", handler: PeopleRefreshHandler{}, options: str.Options{ID: "bryan-cranston"}},
		{name: "shows refresh", handler: ShowsRefreshHandler{}, options: str.Options{InternalID: "breaking-bad"}},
		{name: "users list", handler: UsersListHandler{}, options: str.Options{ID: "1"}},
		{name: "users report", handler: UsersReportHandler{}, options: usersReportOptions},
		{name: "users follow", handler: UsersFollowHandler{}, options: str.Options{UserName: "sean"}},
		{name: "users unfollow", handler: UsersUnfollowHandler{}, options: str.Options{UserName: "sean"}},
		{name: "users block", handler: UsersBlockHandler{}, options: str.Options{UserName: "sean"}},
		{name: "users unblock", handler: UsersUnblockHandler{}, options: str.Options{UserName: "sean"}},
		{name: "users watching", handler: UsersWatchingHandler{}, options: str.Options{UserName: "sean"}},
		{name: "users delete list", handler: UsersDeleteListHandler{}, options: str.Options{UserName: "sean", ID: "1"}},
		{name: "users list report", handler: UsersListReportHandler{}, options: usersListReportOptions},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			s.Teardown() // closed server: every request fails without a response

			options := tt.options
			if err := tt.handler.Handle(&options, s.Client); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

// TestHandlersServerError checks that handlers reading the result do not panic on a 500.
func TestHandlersServerError(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		handler Handler
		options str.Options
	}{
		{name: "lists list", path: "/lists/1", handler: ListsListHandler{}, options: str.Options{InternalID: "1"}},
		{name: "users list", path: "/users/sean/lists/1", handler: UsersListHandler{}, options: str.Options{UserName: "sean", ID: "1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			s.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			options := tt.options
			err := tt.handler.Handle(&options, s.Client)
			if err == nil || !strings.Contains(err.Error(), "500") {
				t.Fatalf("error is %v, want a 500 error", err)
			}
		})
	}
}

// TestUsersReportHandlersWithoutMessage checks 400/409 responses that have no message field.
func TestUsersReportHandlersWithoutMessage(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		status  int
		handler Handler
		options str.Options
	}{
		{name: "users report 400", path: "/users/sean/report", status: http.StatusBadRequest, handler: UsersReportHandler{}, options: usersReportOptions},
		{name: "users report 409", path: "/users/sean/report", status: http.StatusConflict, handler: UsersReportHandler{}, options: usersReportOptions},
		{name: "users list report 400", path: "/users/sean/lists/1/report", status: http.StatusBadRequest, handler: UsersListReportHandler{}, options: usersListReportOptions},
		{name: "users list report 409", path: "/users/sean/lists/1/report", status: http.StatusConflict, handler: UsersListReportHandler{}, options: usersListReportOptions},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			s.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodPost)
				w.WriteHeader(tt.status)
				test.SafeFprint(w, `{}`)
			})

			options := tt.options
			err := tt.handler.Handle(&options, s.Client)
			if err == nil || !strings.Contains(err.Error(), "reason error") {
				t.Fatalf("error is %v, want a reason error", err)
			}
		})
	}
}
