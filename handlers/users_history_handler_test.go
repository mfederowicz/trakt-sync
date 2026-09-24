// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersHistoryHandlerRoutes(t *testing.T) {
	tests := []struct {
		name     string
		histType string
		itemID   int
		path     string
		wantErr  string
	}{
		{name: "defaults", histType: cfg.DefaultConfig().UsersType, itemID: cfg.DefaultConfig().ItemID, path: "/users/sean/history/"},
		{name: "movies", histType: "movies", path: "/users/sean/history/movies"},
		{name: "shows", histType: "shows", path: "/users/sean/history/shows"},
		{name: "episodes", histType: "episodes", path: "/users/sean/history/episodes"},
		{name: "movie item", histType: "movies", itemID: 12601, path: "/users/sean/history/movies/12601"},
		{name: "show item", histType: "shows", itemID: 1388, path: "/users/sean/history/shows/1388"},
		{name: "episode item", histType: "episodes", itemID: 73640, path: "/users/sean/history/episodes/73640"},
		{name: "item without type", itemID: 12601, wantErr: consts.EmptyHistoryItemTypeMsg},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/users/sean/history/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.URL.Path]++
				test.AssertMethod(t, r, http.MethodGet)
				test.SafeFprint(w, `[]`)
			})

			options := &str.Options{
				Module: "users", Action: "history", UserName: "sean",
				Type: tt.histType, ItemID: tt.itemID,
				Output: filepath.Join(t.TempDir(), "out.json"),
			}
			err := UsersHistoryHandler{}.Handle(options, s.Client)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("error is %v, want %q", err, tt.wantErr)
				}
				test.AssertNoDiff(t, map[string]int{}, calls)
				return
			}
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, map[string]int{tt.path: 1}, calls)
		})
	}
}
