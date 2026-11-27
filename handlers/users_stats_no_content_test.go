package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// A 204 from users stats prints a note and writes no file; kept for when Trakt stops sending placeholder zeros to non-VIP users.
func TestUsersStatsHandlerNoContent(t *testing.T) {
	s := setup(t)
	defer s.Teardown()
	s.Mux.HandleFunc("/users/me/stats", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusNoContent)
	})

	options := str.Options{UserName: "me", Output: filepath.Join(t.TempDir(), "out.json")}
	var err error
	out := captureStdout(t, func() { err = UsersStatsHandler{}.Handle(&options, s.Client) })
	test.AssertNilError(t, err)
	assert.Contains(t, out, fmt.Sprintf(consts.UsersStatsNoContent, "me"))
	assert.NotContains(t, out, "Found me user stats")
	_, statErr := os.Stat(options.Output)
	assert.True(t, os.IsNotExist(statErr), "no output file expected")
}
