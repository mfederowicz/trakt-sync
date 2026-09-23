// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestMediaTrendingHandlerFetchesAllPages(t *testing.T) {
	s := setup(t)
	defer s.Teardown()

	pages := map[string]string{
		"1": `[{"watchers":21,"movie":{"title":"Tron: Ares"}}]`,
		"2": `[{"watchers":15,"show":{"title":"Andor"}}]`,
	}
	s.Mux.HandleFunc("/media/trending", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		page := r.URL.Query().Get("page")
		w.Header().Set(internal.HeaderPaginationPage, page)
		w.Header().Set(internal.HeaderPaginationPageCount, "2")
		test.SafeFprint(w, pages[page])
	})

	output := filepath.Join(t.TempDir(), "export_media_trending.json")
	options := &str.Options{Action: "trending", Output: output}
	test.AssertNilError(t, MediaTrendingHandler{}.Handle(options, s.Client))

	data, err := os.ReadFile(output)
	test.AssertNilError(t, err)
	got := []*str.MediaItem{}
	test.AssertNilError(t, json.Unmarshal(data, &got))
	test.AssertNoDiff(t, []*str.MediaItem{
		{Watchers: test.Ptr(21), Movie: &str.Movie{Title: str.String("Tron: Ares")}},
		{Watchers: test.Ptr(15), Show: &str.Show{Title: str.String("Andor")}},
	}, got)
}
