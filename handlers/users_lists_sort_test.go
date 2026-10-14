// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// the lists overview request carries the page, and limit, sort_by and sort_how only when the options hold them.
func TestFetchUsersPersonalListsSort(t *testing.T) {
	tests := []struct {
		name      string
		options   str.Options
		wantQuery string
	}{
		{name: "no sort", options: str.Options{UserName: "sean"}, wantQuery: "page=1"},
		{name: "sort_how only", options: str.Options{UserName: "sean", SortHow: "desc"}, wantQuery: "page=1&sort_how=desc"},
		{name: "both", options: str.Options{UserName: "sean", SortBy: "rank", SortHow: "desc"}, wantQuery: "page=1&sort_by=rank&sort_how=desc"},
		{name: "per page", options: str.Options{UserName: "sean", PerPage: 100}, wantQuery: "limit=100&page=1"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			gotQuery := "not called"
			s.Mux.HandleFunc("/users/sean/lists", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				gotQuery = r.URL.RawQuery
				test.SafeFprint(w, `[{"name":"Favorites","ids":{"trakt":1}}]`)
			})

			options := tt.options
			lists, err := fetchUsersPersonalLists(s.Client, &options, consts.DefaultPage)
			test.AssertNilError(t, err)
			assert.Equal(t, tt.wantQuery, gotQuery)
			assert.Len(t, lists, 1)
		})
	}
}
