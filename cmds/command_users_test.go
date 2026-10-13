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
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt/trakttest"
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

	setup := trakttest.Setup()
	defer setup.Teardown()
	var gotStartAt, gotEndAt string
	setup.Mux.HandleFunc("/users/sean/history/movies", func(w http.ResponseWriter, r *http.Request) {
		gotStartAt, gotEndAt = r.URL.Query().Get("start_at"), r.URL.Query().Get("end_at")
		_, _ = w.Write([]byte(`[{"id":1}]`))
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

	setup := trakttest.Setup()
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

// TestUsersUpdateListSettings checks that users -a update_list sends a list setting from its flag,
// and without the flag keeps the value the list has.
func TestUsersUpdateListSettings(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/usr/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/usr/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/usr/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))

	// every setting differs from its flag default (private, no comments, numbers, rank, asc)
	const list = `{"name":"list","ids":{"trakt":123456},"privacy":"public","display_numbers":false,"allow_comments":true,"sort_by":"added","sort_how":"desc"}`
	kept := map[string]any{"privacy": "public", "display_numbers": false, "allow_comments": true, "sort_by": "added", "sort_how": "desc"}
	with := func(key string, value any) map[string]any {
		want := map[string]any{key: value}
		for k, v := range kept {
			if k != key {
				want[k] = v
			}
		}
		return want
	}

	tests := []struct {
		name string
		list string
		args []string
		want map[string]any // a key that is missing must not be sent
	}{
		{name: "no flags keep every setting", list: list, want: kept},
		{name: "description only keeps every setting", list: list, args: []string{"-description", "short description"}, want: kept},
		{name: "display_numbers on", list: list, args: []string{"-display_numbers"}, want: with("display_numbers", true)},
		{name: "display_numbers off", list: `{"name":"list","ids":{"trakt":123456},"display_numbers":true}`, args: []string{"-display_numbers=false"},
			want: map[string]any{"display_numbers": false}},
		{name: "allow_comments off", list: list, args: []string{"-allow_comments=false"}, want: with("allow_comments", false)},
		{name: "privacy", list: list, args: []string{"-privacy", "friends"}, want: with("privacy", "friends")},
		{name: "privacy same as the default", list: list, args: []string{"-privacy", "private"}, want: with("privacy", "private")},
		{name: "sort_by", list: list, args: []string{"-sort_by", "rank"}, want: with("sort_by", "rank")},
		{name: "sort_how", list: list, args: []string{"-sort_how", "asc"}, want: with("sort_how", "asc")},
		{name: "list without settings and no flags", list: `{"name":"list","ids":{"trakt":123456}}`, want: map[string]any{}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			config := cfg.DefaultConfig()
			config.ClientID, config.ClientSecret = "a", "b"
			config.TokenPath, config.SettingsPath = "/usr/token.json", "/usr/user_settings.json"

			resetAllFlags()
			t.Cleanup(resetAllFlags)

			setup := trakttest.Setup()
			defer setup.Teardown()
			sent := map[string]any{}
			setup.Mux.HandleFunc("/users/sean/lists/123456", func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					_, _ = w.Write([]byte(tt.list))
					return
				}
				assert.NoError(t, json.NewDecoder(r.Body).Decode(&sent))
				// fail the update so the handler stops before writing its result file
				w.WriteHeader(http.StatusInternalServerError)
			})

			args := append([]string{"-a", "update_list", "-u", "sean", "-i", "123456"}, tt.args...)
			assert.Error(t, UsersCmd.Exec(fs, setup.Client, config, args))
			assert.Equal(t, "list", sent["name"], "the update request was sent")
			assert.NotContains(t, sent, "display_strings")
			got := map[string]any{}
			for _, key := range []string{"privacy", "display_numbers", "allow_comments", "sort_by", "sort_how"} {
				if value, ok := sent[key]; ok {
					got[key] = value
				}
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
