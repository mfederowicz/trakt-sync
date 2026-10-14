// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// users -a lists without -i ends with the overview; with -i the list must be one of the user's lists.
func TestUsersListsHandlerOverview(t *testing.T) {
	// the lists overview is written to the working directory
	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.Chdir(t.TempDir()))
	t.Cleanup(func() { _ = os.Chdir(wd) })

	const lists = `[{"name":"list","item_count":1,"ids":{"trakt":1},"user":{"name":"sean"},"images":{"posters":["p.webp"]}}]`
	tests := []struct {
		name      string
		options   str.Options
		wantQuery string
		wantErr   string
	}{
		{name: "no id", options: str.Options{}, wantQuery: "page=1"},
		{name: "no id with extended info", options: str.Options{ExtendedInfo: "images"}, wantQuery: "extended=images&page=1"},
		{name: "no id with sort", options: str.Options{SortHow: "desc"}, wantQuery: "page=1&sort_how=desc"},
		{name: "id that is not a number", options: str.Options{ID: "my-list"}, wantQuery: "page=1", wantErr: "please set personal listid"},
		{name: "id of another list", options: str.Options{ID: "2"}, wantQuery: "page=1", wantErr: "unknown listid:2"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			paths := []string{}
			s.Mux.HandleFunc("/users/sean/lists", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				paths = append(paths, r.URL.Path)
				assert.Equal(t, tt.wantQuery, r.URL.RawQuery)
				test.SafeFprint(w, lists)
			})
			s.Mux.HandleFunc("/users/sean/lists/", func(_ http.ResponseWriter, r *http.Request) {
				paths = append(paths, r.URL.Path)
			})

			overview := "export_users_lists.json"
			assert.NoError(t, os.RemoveAll(overview))
			options := tt.options
			options.Module = consts.Users
			options.Action = consts.Lists
			options.UserName = "sean"
			options.Output = filepath.Join(t.TempDir(), "items.json")
			err := UsersListsHandler{}.Handle(&options, s.Client)

			assert.Equal(t, []string{"/users/sean/lists"}, paths, "only the overview is requested")
			data, readErr := os.ReadFile(overview)
			assert.NoError(t, readErr)
			assert.Contains(t, string(data), `"p.webp"`)
			assert.NoFileExists(t, options.Output)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			test.AssertNilError(t, err)
		})
	}
}

// lists -a list sends -ex and keeps the images of the list.
func TestListsListHandlerExtendedInfo(t *testing.T) {
	tests := []struct {
		name      string
		extended  string
		wantQuery string
	}{
		{name: "images", extended: "images", wantQuery: "extended=images"},
		{name: "no extended info"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			s.Mux.HandleFunc("/lists/55", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				assert.Equal(t, tt.wantQuery, r.URL.RawQuery)
				test.SafeFprint(w, `{"name":"Star Wars","images":{"posters":["p.webp"]}}`)
			})

			options := str.Options{Module: consts.Lists, Action: consts.List, InternalID: "55", ExtendedInfo: tt.extended}
			options.Output = filepath.Join(t.TempDir(), "out.json")
			test.AssertNilError(t, ListsListHandler{}.Handle(&options, s.Client))
			data, err := os.ReadFile(options.Output)
			assert.NoError(t, err)
			assert.Contains(t, string(data), `"p.webp"`)
		})
	}
}
