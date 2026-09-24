// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

// The comments module fills only InternalID (-i / -trakt_id); the body carries just the numeric Trakt ID.
func TestCommentsCommentsSeasonEpisodeHandlers(t *testing.T) {
	tests := []struct {
		name    string
		handler Handler
		id      string
		want    str.Comment
		wantErr string
	}{
		{name: "season", handler: CommentsCommentsSeasonHandler{}, id: "3950", want: str.Comment{Season: &str.Season{IDs: &str.IDs{Trakt: test.Ptr(int64(3950))}}}},
		{name: "episode", handler: CommentsCommentsEpisodeHandler{}, id: "73482", want: str.Comment{Episode: &str.Episode{IDs: &str.IDs{Trakt: test.Ptr(int64(73482))}}}},
		{name: "season without id", handler: CommentsCommentsSeasonHandler{}, wantErr: consts.EmptyTraktIDMsg},
		{name: "episode slug is rejected", handler: CommentsCommentsEpisodeHandler{}, id: "the-sopranos", wantErr: "trakt id must be a positive number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			MuxUserSettings(t, s.Mux)

			posted := 0
			s.Mux.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
				posted++
				test.AssertMethod(t, r, http.MethodPost)
				got := new(str.Comment)
				test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
				test.AssertNoDiff(t, tt.want.Season, got.Season)
				test.AssertNoDiff(t, tt.want.Episode, got.Episode)
				w.WriteHeader(http.StatusCreated)
				test.SafeFprint(w, `{"id":1}`)
			})
			s.Mux.HandleFunc("/seasons/", func(_ http.ResponseWriter, r *http.Request) {
				t.Errorf("unexpected season lookup: %s", r.URL.Path)
			})
			s.Mux.HandleFunc("/episodes/", func(_ http.ResponseWriter, r *http.Request) {
				t.Errorf("unexpected episode lookup: %s", r.URL.Path)
			})

			err := tt.handler.Handle(&str.Options{InternalID: tt.id, Comment: "great, more than five words here"}, s.Client)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				if posted != 0 {
					t.Error("comment was posted on error")
				}
				return
			}
			test.AssertNilError(t, err)
			if posted != 1 {
				t.Errorf("posted %d comments, want 1", posted)
			}
		})
	}
}
