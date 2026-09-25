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

// TestUsersWriteHandlersUseOutput checks that users write handlers save their result to options.Output (-o).
func TestUsersWriteHandlersUseOutput(t *testing.T) {
	// handlers that also write a second, fixed-name file do it in the working directory
	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.Chdir(t.TempDir()))
	t.Cleanup(func() { _ = os.Chdir(wd) })

	listJSON := `{"name":"list","item_count":1,"ids":{"trakt":1},"user":{"name":"sean"}}`
	input := filepath.Join(t.TempDir(), "list.json")
	assert.NoError(t, os.WriteFile(input, []byte(`{"name":"list"}`), 0o600))

	tests := []struct {
		name      string
		handler   Handler
		options   str.Options
		routes    map[string]string
		wantExtra string
	}{
		{
			name:    "watching",
			handler: UsersWatchingHandler{},
			options: str.Options{Module: "users", Action: "watching", UserName: "sean"},
			routes:  map[string]string{"/users/sean/watching": `{"action":"watch"}`},
		},
		{
			name:    "update_list",
			handler: UsersUpdateListHandler{},
			options: str.Options{Module: "users", Action: "update_list", UserName: "sean", ID: "1", Privacy: "private"},
			routes:  map[string]string{"/users/sean/lists/1": listJSON},
		},
		{
			name:    "add_list",
			handler: UsersAddListHandler{},
			options: str.Options{Module: "users", Action: "add_list", UserName: "sean", Items: input},
			routes:  map[string]string{"/users/sean/lists": listJSON},
		},
		{
			name:      "lists with -i",
			handler:   UsersListsHandler{},
			options:   str.Options{Module: "users", Action: "lists", UserName: "sean", ID: "1", Type: "movies"},
			routes:    map[string]string{"/users/sean/lists": "[" + listJSON + "]", "/users/sean/lists/1/items/movies": `[{"rank":1}]`},
			wantExtra: "export_users_lists.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSetup := setup(t)
			defer testSetup.Teardown()
			for path, body := range tt.routes {
				testSetup.Mux.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte(body))
				})
			}

			output := filepath.Join(t.TempDir(), "out.json")
			options := tt.options
			options.Output = output
			assert.NoError(t, tt.handler.Handle(&options, testSetup.Client))
			assert.FileExists(t, output)
			if tt.wantExtra != "" {
				assert.FileExists(t, tt.wantExtra)
			}
		})
	}
}
