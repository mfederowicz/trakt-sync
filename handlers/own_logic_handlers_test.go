// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// A show without a scheduled episode answers 204: the handler reports it and writes no file.
func TestShowsNextEpisodeNoContent(t *testing.T) {
	s := setup(t)
	defer s.Teardown()
	s.Mux.HandleFunc("/shows/55/next_episode", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusNoContent)
	})

	options := str.Options{InternalID: "55", Output: filepath.Join(t.TempDir(), "out.json")}
	test.AssertNilError(t, ShowsNextEpisodeHandler{}.Handle(&options, s.Client))
	_, err := os.Stat(options.Output)
	assert.True(t, os.IsNotExist(err), "nothing is written without a next episode")
}

// inTempDir runs the test in a temporary working directory, for handlers that write a file with a fixed name.
func inTempDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	previous, err := os.Getwd()
	test.AssertNilError(t, err)
	test.AssertNilError(t, os.Chdir(dir))
	t.Cleanup(func() {
		test.AssertNilError(t, os.Chdir(previous))
	})
	return dir
}

// add_to_history first removes the items from history, then adds them; each step writes its own result file.
func TestSyncAddToHistoryCleansThenAdds(t *testing.T) {
	const items = `[{"type":"movie","watched_at":"2026-10-01T10:00:00.000Z","movie":{"title":"Tron","ids":{"trakt":1}}}]`
	s := setup(t)
	defer s.Teardown()
	requests := []string{}
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.Path)
		body, err := io.ReadAll(r.Body)
		test.AssertNilError(t, err)
		assert.Contains(t, string(body), `"ids":{"trakt":1}`)
		test.SafeFprint(w, `{}`)
	})

	dir := inTempDir(t)
	options := str.Options{Module: "sync", Action: "add_to_history", Type: "movies", Output: filepath.Join(dir, "out.json"), Items: filepath.Join(dir, "items.json")}
	test.AssertNilError(t, os.WriteFile(options.Items, []byte(items), 0o600))

	test.AssertNilError(t, SyncAddToHistoryHandler{}.Handle(&options, s.Client))
	assert.Equal(t, []string{"POST /sync/history/remove", "POST /sync/history"}, requests)
	for _, name := range []string{"sync_remove_from_history_results.json", "out.json"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		test.AssertNilError(t, err)
		assert.True(t, json.Valid(data), "%s is JSON", name)
	}
}

func TestSyncAddToHistoryFailedCleanup(t *testing.T) {
	s := setup(t)
	defer s.Teardown()
	requests := []string{}
	s.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	})

	dir := inTempDir(t)
	options := str.Options{Module: "sync", Action: "add_to_history", Type: "movies", Output: filepath.Join(dir, "out.json"), Items: filepath.Join(dir, "items.json")}
	test.AssertNilError(t, os.WriteFile(options.Items, []byte(`[]`), 0o600))

	assert.Error(t, SyncAddToHistoryHandler{}.Handle(&options, s.Client))
	assert.Equal(t, []string{"POST /sync/history/remove"}, requests, "nothing is added when the cleanup fails")
	entries, err := os.ReadDir(dir)
	test.AssertNilError(t, err)
	assert.Len(t, entries, 1, "only the items file is in the directory")
}
