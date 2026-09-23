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

func TestCommentsReactionHandler(t *testing.T) {
	tests := []struct {
		name    string
		options str.Options
		method  string
		wantErr string
	}{
		{name: "add", options: str.Options{CommentID: 417, Reaction: "love"}, method: http.MethodPost},
		{name: "remove", options: str.Options{CommentID: 417, Reaction: "love", Remove: true}, method: http.MethodDelete},
		{name: "missing comment id", options: str.Options{Reaction: "love"}, wantErr: consts.EmptyCommentIDMsg},
		{name: "missing reaction", options: str.Options{CommentID: 417}, wantErr: consts.EmptyReactionMsg},
		{name: "invalid reaction", options: str.Options{CommentID: 417, Reaction: "angry"}, wantErr: "reaction 'angry' is not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			called := false
			s.Mux.HandleFunc("/comments/417/reactions/love", func(w http.ResponseWriter, r *http.Request) {
				called = true
				test.AssertMethod(t, r, tt.method)
				w.WriteHeader(http.StatusNoContent)
			})

			options := tt.options
			options.Module = consts.Comments
			options.Action = consts.Reaction
			err := CommentsReactionHandler{}.Handle(&options, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				if called {
					t.Error("API was called for invalid options")
				}
				return
			}
			test.AssertNilError(t, err)
			if !called {
				t.Error("API was not called")
			}
		})
	}
}

func TestCommentsReportHandler(t *testing.T) {
	tests := []struct {
		name    string
		options str.Options
		wantErr string
	}{
		{name: "report", options: str.Options{CommentID: 417, Reason: "spoilers"}},
		{name: "missing reason", options: str.Options{CommentID: 417}, wantErr: consts.EmptyReasonMsg},
		{name: "invalid reason", options: str.Options{CommentID: 417, Reason: "boring"}, wantErr: "reason 'boring' is not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			s.Mux.HandleFunc("/comments/417/report", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodPost)
				w.WriteHeader(http.StatusCreated)
			})

			options := tt.options
			options.Module = consts.Comments
			options.Action = consts.Report
			err := CommentsReportHandler{}.Handle(&options, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			test.AssertNilError(t, err)
		})
	}
}
