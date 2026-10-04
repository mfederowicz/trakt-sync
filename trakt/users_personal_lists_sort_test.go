package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

// GetUsersPersonalLists sends sort_by and sort_how when set and no query for nil or empty options.
func TestUsersGetUsersPersonalListsSort(t *testing.T) {
	tests := []struct {
		name      string
		opts      *uri.ListOptions
		wantQuery string
	}{
		{name: "nil options", opts: nil, wantQuery: ""},
		{name: "empty options", opts: &uri.ListOptions{}, wantQuery: ""},
		{name: "sort_how only", opts: &uri.ListOptions{SortHow: "desc"}, wantQuery: "sort_how=desc"},
		{name: "both", opts: &uri.ListOptions{SortBy: "rank", SortHow: "desc"}, wantQuery: "sort_by=rank&sort_how=desc"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()
			gotQuery := "not called"
			setup.Mux.HandleFunc("/users/sean/lists", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				gotQuery = r.URL.RawQuery
				test.SafeFprint(w, `[{"name":"Favorites","ids":{"trakt":1}}]`)
			})

			lists, _, err := setup.Client.Users.GetUsersPersonalLists(context.Background(), "sean", tt.opts)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, tt.wantQuery, gotQuery)
			test.AssertNoDiff(t, []*str.PersonalList{{Name: str.String("Favorites"), IDs: &str.IDs{Trakt: test.Ptr(int64(1))}}}, lists)
		})
	}
}
