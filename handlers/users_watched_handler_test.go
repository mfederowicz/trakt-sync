// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersWatchedHandlerRoutes(t *testing.T) {
	for _, watchType := range []string{"movies", "shows"} {
		t.Run(watchType, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/users/sean/watched/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.URL.Path]++
				test.AssertMethod(t, r, http.MethodGet)
				test.SafeFprint(w, `[]`)
			})

			options := &str.Options{UserName: "sean", Type: watchType, Output: filepath.Join(t.TempDir(), "out.json")}
			test.AssertNilError(t, UsersWatchedHandler{}.Handle(options, s.Client))
			test.AssertNoDiff(t, map[string]int{"/users/sean/watched/" + watchType: 1}, calls)
		})
	}
}
