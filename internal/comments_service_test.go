// Package internal used for client and services
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestCommentsServiceGetCommentReactions(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/comments/417/reactions", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		if got, want := r.URL.Query().Get("page"), "2"; got != want {
			t.Errorf("page query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `[{"reaction":{"type":"love"},"user":{"username":"sean"}}]`)
	})

	got, _, err := setup.Client.Comments.GetCommentReactions(context.Background(), test.Ptr(417), &uri.ListOptions{Page: 2})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.CommentReaction{
		{Reaction: &str.Reaction{Type: str.String("love")}, User: &str.UserProfile{Userame: str.String("sean")}},
	}, got)
}

func TestCommentsServiceGetCommentReactionsSummary(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/comments/417/reactions/summary", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `{"reaction_count":3,"user_count":2,"distribution":{"like":2,"love":1}}`)
	})

	got, _, err := setup.Client.Comments.GetCommentReactionsSummary(context.Background(), test.Ptr(417))
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, &str.ReactionSummary{
		ReactionCount: test.Ptr(3),
		UserCount:     test.Ptr(2),
		Distribution:  map[string]int{"like": 2, "love": 1},
	}, got)
}

func TestCommentsServiceCommentReaction(t *testing.T) {
	tests := []struct {
		name   string
		method string
		status int
		call   func(s *CommentsService, id *int, reaction *string) (*str.Response, error)
	}{
		{
			name:   "add reaction",
			method: http.MethodPost,
			status: http.StatusCreated,
			call: func(s *CommentsService, id *int, reaction *string) (*str.Response, error) {
				return s.AddCommentReaction(context.Background(), id, reaction)
			},
		},
		{
			name:   "remove reaction",
			method: http.MethodDelete,
			status: http.StatusNoContent,
			call: func(s *CommentsService, id *int, reaction *string) (*str.Response, error) {
				return s.RemoveCommentReaction(context.Background(), id, reaction)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc("/comments/417/reactions/love", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, tt.method)
				w.WriteHeader(tt.status)
			})

			resp, err := tt.call(setup.Client.Comments, test.Ptr(417), str.String("love"))
			test.AssertNilError(t, err)
			if got := resp.StatusCode; got != tt.status {
				t.Errorf("status code is %d, want %d", got, tt.status)
			}
		})
	}
}

func TestCommentsServiceReportComment(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/comments/417/report", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.CommentReport)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.CommentReport{Reason: str.String("spam"), Message: str.String("ad link")}, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Comments.ReportComment(context.Background(), test.Ptr(417), &str.CommentReport{Reason: str.String("spam"), Message: str.String("ad link")})
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}

func TestCommentsServiceReportCommentConflict(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/comments/417/report", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
		test.SafeFprint(w, `{"message":"report already pending"}`)
	})

	_, err := setup.Client.Comments.ReportComment(context.Background(), test.Ptr(417), &str.CommentReport{Reason: str.String("spam")})
	if got, want := fmt.Sprint(err), fmt.Sprintf(consts.CommentReportPending, 417); got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}
