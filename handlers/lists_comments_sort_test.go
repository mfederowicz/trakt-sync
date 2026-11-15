// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// lists -a comments takes the sorts the live API accepts; every other value
// answered 400 on 2026-11-15 and now stops before the request.
func TestListsCommentsHandlerSort(t *testing.T) {
	tests := []struct {
		sort     string
		wantPath string
	}{
		{sort: consts.EmptyString, wantPath: "/lists/55/comments"},
		{sort: "likes", wantPath: "/lists/55/comments/likes"},
		{sort: "replies", wantPath: "/lists/55/comments/replies"},
		{sort: "newest", wantPath: "/lists/55/comments/newest"},
		{sort: "oldest", wantPath: "/lists/55/comments/oldest"},
		{sort: "likes_30"},
		{sort: "replies_30"},
		{sort: "plays"},
		{sort: "rating"},
		{sort: "added"},
		{sort: "highest"},
		{sort: "lowest"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run("sort "+tt.sort, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			paths := []string{}
			s.Mux.HandleFunc("/lists/55/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				paths = append(paths, r.URL.Path)
				test.SafeFprint(w, `[{"id":1,"comment":"nice list"}]`)
			})
			s.Mux.HandleFunc("/lists/55/comments", func(w http.ResponseWriter, r *http.Request) {
				paths = append(paths, r.URL.Path)
				test.SafeFprint(w, `[{"id":1,"comment":"nice list"}]`)
			})

			options := &str.Options{Module: "lists", Action: "comments", InternalID: "55", Sort: tt.sort, Output: filepath.Join(t.TempDir(), "out.json")}
			err := ListsCommentsHandler{}.Handle(options, s.Client)
			if tt.wantPath == consts.EmptyString {
				assert.ErrorContains(t, err, "sort '"+tt.sort+"' is not valid")
				assert.Empty(t, paths, "no request is sent")
				return
			}
			test.AssertNilError(t, err)
			assert.Equal(t, []string{tt.wantPath}, paths)
		})
	}
}
