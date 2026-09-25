// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

// TestUsersServiceTypedListItems checks the typed list items routes of the contract.
func TestUsersServiceTypedListItems(t *testing.T) {
	for _, typ := range []string{"movie", "show", "movie,show", "movie,show,season,episode"} {
		typ := typ
		t.Run(typ, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			wantPath := "/users/sean/lists/55/items/" + typ
			calls := 0
			setup.Mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				if r.URL.Path != wantPath {
					t.Errorf("path is %q, want %q", r.URL.Path, wantPath)
				}
				test.SafeFprint(w, `[{"rank":1,"type":"movie","movie":{"title":"Arrival"}}]`)
			})

			got, _, err := setup.Client.Users.GetItemstOnAPersonalList(context.Background(), str.String("sean"), str.String("55"), &typ)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, []*str.UserListItem{{Rank: test.Ptr(1), Type: str.String("movie"), Movie: &str.Movie{Title: str.String("Arrival")}}}, got)
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
		})
	}
}
