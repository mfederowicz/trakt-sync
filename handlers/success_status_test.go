// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/stretchr/testify/assert"
)

// TestSuccessStatusIsNotAnError checks that success responses (201 created, 204 removed) do not come back as errors.
func TestSuccessStatusIsNotAnError(t *testing.T) {
	input := filepath.Join(t.TempDir(), "list.json")
	assert.NoError(t, os.WriteFile(input, []byte(`{"name":"list"}`), 0o600))

	tests := []struct {
		name      string
		handler   Handler
		options   str.Options
		path      string
		status    int
		body      string
		wantWrite bool
	}{
		{
			name:      "users add_list created",
			handler:   UsersAddListHandler{},
			options:   str.Options{Module: "users", Action: "add_list", UserName: "sean", Items: input},
			path:      "/users/sean/lists",
			status:    http.StatusCreated,
			body:      `{"name":"list","ids":{"trakt":1}}`,
			wantWrite: true,
		},
		{
			name:    "sync remove_playback removed",
			handler: SyncRemovePlaybackHandler{},
			options: str.Options{Module: "sync", Action: "remove_playback", PlaybackID: 13},
			path:    "/sync/playback/13",
			status:  http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSetup := setup(t)
			defer testSetup.Teardown()
			testSetup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			})

			output := filepath.Join(t.TempDir(), "out.json")
			options := tt.options
			options.Output = output
			assert.NoError(t, tt.handler.Handle(&options, testSetup.Client))
			if tt.wantWrite {
				assert.FileExists(t, output)
			}
		})
	}
}
