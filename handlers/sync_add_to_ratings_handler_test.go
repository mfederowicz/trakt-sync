// Package handlers used to handle module actions
package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// TestSyncAddToRatingsMixedShowExport sends a ratings export that holds a show, one of its seasons and one of its
// episodes, the episode first. The show must keep its own rating, and seasons and episodes go in their own lists.
func TestSyncAddToRatingsMixedShowExport(t *testing.T) {
	const export = `[
		{"rated_at":"2026-10-01T10:00:00.000Z","rating":7,"type":"episode","episode":{"season":1,"number":2,"ids":{"trakt":502}},"show":{"title":"Dark","year":2017,"ids":{"trakt":1}}},
		{"rated_at":"2026-10-02T10:00:00.000Z","rating":8,"type":"season","season":{"number":1,"ids":{"trakt":301}},"show":{"title":"Dark","year":2017,"ids":{"trakt":1}}},
		{"rated_at":"2026-10-03T10:00:00.000Z","rating":9,"type":"show","show":{"title":"Dark","year":2017,"ids":{"trakt":1}}}
	]`
	const show = `{"ids":{"trakt":1},"rating":9,"rated_at":"2026-10-03T10:00:00Z"}`
	tests := []struct {
		name  string
		stype string
		want  string
	}{
		{name: "shows", stype: consts.Shows, want: `{"movies":[],"shows":[` + show + `],"seasons":[],"episodes":[]}`},
		{name: "all", stype: consts.ActionTypeAll, want: `{"movies":[],"shows":[` + show + `],` +
			`"seasons":[{"ids":{"trakt":301},"rating":8,"rated_at":"2026-10-02T10:00:00Z"}],` +
			`"episodes":[{"ids":{"trakt":502},"rating":7,"rated_at":"2026-10-01T10:00:00Z"}]}`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			sent := ""
			s.Mux.HandleFunc("/sync/ratings", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodPost)
				body, err := io.ReadAll(r.Body)
				test.AssertNilError(t, err)
				sent = string(body)
				w.WriteHeader(http.StatusCreated)
				test.SafeFprint(w, `{}`)
			})

			dir := t.TempDir()
			options := str.Options{Module: consts.Sync, Action: consts.AddToRatings, Type: tt.stype,
				Items: filepath.Join(dir, "items.json"), Output: filepath.Join(dir, "out.json")}
			test.AssertNilError(t, os.WriteFile(options.Items, []byte(export), 0o600))

			test.AssertNilError(t, SyncAddToRatingsHandler{}.Handle(&options, s.Client))
			assert.JSONEq(t, tt.want, sent)
		})
	}
}
