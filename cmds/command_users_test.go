// Package cmds used for commands modules
package cmds

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestUsersHistoryDates checks that users -a history sends -start_at / -end_at.
func TestUsersHistoryDates(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/usr/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/usr/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/usr/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/usr/token.json", "/usr/user_settings.json"

	resetAllFlags()
	t.Cleanup(resetAllFlags)

	setup := internal.Setup()
	defer setup.Teardown()
	var gotStartAt, gotEndAt string
	setup.Mux.HandleFunc("/users/sean/history/movies", func(w http.ResponseWriter, r *http.Request) {
		gotStartAt, gotEndAt = r.URL.Query().Get("start_at"), r.URL.Query().Get("end_at")
		_, _ = w.Write([]byte(`[]`))
	})

	args := []string{"-o", filepath.Join(t.TempDir(), "out.json"), "-a", "history", "-u", "sean", "-start_at", "2026-01-15", "-end_at", "2026-01-20"}
	assert.NoError(t, UsersCmd.Exec(fs, setup.Client, config, args))

	// ConvertDateString gives a date-only value the current full hour, so compute it the same way
	common := handlers.CommonLogic{}
	tz := UsersCmd.Options.Timezone
	assert.Equal(t, common.ConvertDateString("2026-01-15", consts.DefaultStartDateFormat, tz, true), gotStartAt)
	assert.Equal(t, common.ConvertDateString("2026-01-20", consts.DefaultStartDateFormat, tz, true), gotEndAt)
}

// TestUsersUpdateListDescription checks that users -a update_list sends -description.
func TestUsersUpdateListDescription(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/usr/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/usr/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/usr/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/usr/token.json", "/usr/user_settings.json"

	resetAllFlags()
	t.Cleanup(resetAllFlags)

	setup := internal.Setup()
	defer setup.Teardown()
	var gotDescription string
	setup.Mux.HandleFunc("/users/sean/lists/123456", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = w.Write([]byte(`{"name":"list","ids":{"trakt":123456}}`))
			return
		}
		var list str.PersonalList
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&list))
		if list.Description != nil {
			gotDescription = *list.Description
		}
		// fail the update so the handler stops before writing its result file
		w.WriteHeader(http.StatusInternalServerError)
	})

	args := []string{"-a", "update_list", "-u", "sean", "-i", "123456", "-description", "short description"}
	assert.Error(t, UsersCmd.Exec(fs, setup.Client, config, args))
	assert.Equal(t, "short description", gotDescription)
}
