// Package handlers used to handle module actions
package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// exportOnlyKeys are the keys an export made with -ex full,images,colors holds and a write request must not send back.
var exportOnlyKeys = []string{
	"after_credits", "colors", "during_credits", "genres", "images", "last_aired",
	"original_title", "overview", "social_ids", "subgenres", "total_runtime",
}

// fullExportItems returns an -items file with one movie or show as an export with -ex full,images,colors writes it.
func fullExportItems(kind string) string {
	const media = `"title":"Ida","year":2013,"ids":{"trakt":1},"overview":"A novice nun.","genres":["drama"],"subgenres":["road-trip"],"original_title":"Ida",` +
		`"images":{"fanart":["f.webp"],"poster":["p.webp"],"logo":["l.webp"],"clearart":["c.webp"],"banner":["b.webp"],"thumb":["t.webp"]},` +
		`"colors":{"poster":["#111111","#222222"]},"social_ids":{"twitter":"ida","facebook":"idafilm","instagram":"ida","wikipedia":"Ida_(film)"},` +
		`"after_credits":true,"during_credits":false,"last_aired":"2026-09-30T20:00:00Z","total_runtime":480`
	return `[{"id":7,"type":"` + kind + `","rating":8,"rated_at":"2026-10-01T10:00:00.000Z","watched_at":"2026-10-01T10:00:00.000Z","` + kind + `":{` + media + `}}]`
}

func assertNoExportOnlyKeys(t *testing.T, body string) {
	t.Helper()
	assert.Contains(t, body, `"trakt":1`, "the item is sent")
	for _, key := range exportOnlyKeys {
		assert.NotContains(t, body, `"`+key+`"`)
	}
}

// the write actions that read a movie or a show from -items send its ids, title and dates, not the rest of the export.
func TestWriteHandlersDropExportOnlyKeys(t *testing.T) {
	for _, tc := range writeHandlers() {
		if tc.options.Type != consts.Movies || tc.items == consts.EmptyString || strings.Contains(tc.sends, "rank") {
			continue
		}
		for kind, stype := range map[string]string{consts.Movie: consts.Movies, consts.Show: consts.Shows} {
			tc, kind, stype := tc, kind, stype
			t.Run(tc.name+" "+kind, func(t *testing.T) {
				s := setup(t)
				defer s.Teardown()
				sent := consts.EmptyString
				s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					test.AssertNilError(t, err)
					sent = string(body)
					test.SafeFprint(w, `{}`)
				})

				tc.items = fullExportItems(kind)
				options := writeOptions(t, tc)
				options.Type = stype
				test.AssertNilError(t, tc.handler.Handle(&options, s.Client))
				assertNoExportOnlyKeys(t, sent)
			})
		}
	}
}

// add_to_history sends two requests, the cleanup and the add; neither carries the rest of the export.
func TestSyncAddToHistoryDropsExportOnlyKeys(t *testing.T) {
	for kind, stype := range map[string]string{consts.Movie: consts.Movies, consts.Show: consts.Shows} {
		kind, stype := kind, stype
		t.Run(kind, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			bodies := map[string]string{}
			s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				test.AssertNilError(t, err)
				bodies[r.URL.Path] = string(body)
				test.SafeFprint(w, `{}`)
			})

			dir := inTempDir(t)
			options := str.Options{Module: "sync", Action: consts.AddToHistory, Type: stype, Output: filepath.Join(dir, "out.json"), Items: filepath.Join(dir, "items.json")}
			test.AssertNilError(t, os.WriteFile(options.Items, []byte(fullExportItems(kind)), 0o600))

			test.AssertNilError(t, SyncAddToHistoryHandler{}.Handle(&options, s.Client))
			assert.Len(t, bodies, consts.TwoValue)
			for _, body := range bodies {
				assertNoExportOnlyKeys(t, body)
			}
		})
	}
}
