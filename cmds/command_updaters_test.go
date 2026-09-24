// Package cmds used for commands modules
package cmds

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// resetAllFlags restores every flag default; flag values are package globals shared by all Exec calls.
func resetAllFlags() {
	reset := func(f *flag.Flag) {
		if !strings.HasPrefix(f.Name, "test.") { // leave go test's own flags alone
			_ = f.Value.Set(f.DefValue)
		}
	}
	flag.VisitAll(reset)
	for _, c := range Commands {
		c.Flag.VisitAll(reset)
	}
}

// TestModuleFlagUpdaters checks that each command's -period / -start_date reach its own request
// and that date defaults keep the last consts.DefaultStartAtDays window.
func TestModuleFlagUpdaters(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/upd/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/upd/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/upd/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/upd/token.json", "/upd/user_settings.json"

	// ConvertDateString gives a date-only -start_date the current full hour, so compute it the same way
	startDate := func(tz string) string {
		common := handlers.CommonLogic{}
		return common.ConvertDateString("2026-01-15", consts.DefaultStartDateFormat, tz, true)
	}
	window := func(tz string) string {
		common := handlers.CommonLogic{}
		return common.DateLastDays(consts.DefaultStartAtDays, tz, true)
	}
	tests := []struct {
		name string
		cmd  *Command
		args []string
		path func(tz string) string
	}{
		{name: "movies -period", cmd: MoviesCmd, args: []string{"-a", "favorited", "-period", "daily"}, path: func(string) string { return "/movies/favorited/daily" }},
		{name: "movies default period", cmd: MoviesCmd, args: []string{"-a", "favorited"}, path: func(string) string { return "/movies/favorited/" + cfg.DefaultConfig().MoviesPeriod }},
		{name: "movies -start_date", cmd: MoviesCmd, args: []string{"-a", "updates", "-start_date", "2026-01-15"}, path: func(tz string) string { return "/movies/updates/" + startDate(tz) }},
		{name: "movies default start", cmd: MoviesCmd, args: []string{"-a", "updates"}, path: func(tz string) string { return "/movies/updates/" + window(tz) }},
		{name: "movies streaming -period", cmd: MoviesCmd, args: []string{"-a", "streaming", "-period", "daily"}, path: func(string) string { return "/movies/streaming/daily" }},
		{name: "movies streaming default period", cmd: MoviesCmd, args: []string{"-a", "streaming"}, path: func(string) string { return "/movies/streaming/" + cfg.DefaultConfig().MoviesPeriod }},
		{name: "movies streaming -period all is rejected", cmd: MoviesCmd, args: []string{"-a", "streaming", "-period", "all"}, path: func(string) string { return "" }},
		{name: "movies report", cmd: MoviesCmd, args: []string{"-a", "report", "-i", "tron-legacy-2010", "-r", "spam"}, path: func(string) string { return "/movies/tron-legacy-2010/report" }},
		{name: "movies report default reason is rejected", cmd: MoviesCmd, args: []string{"-a", "report", "-i", "tron-legacy-2010"}, path: func(string) string { return "" }},
		{name: "movies refresh_justwatch", cmd: MoviesCmd, args: []string{"-a", "refresh_justwatch", "-i", "tron-legacy-2010"}, path: func(string) string { return "/movies/tron-legacy-2010/refresh/justwatch" }},
		{name: "movies sentiments", cmd: MoviesCmd, args: []string{"-a", "sentiments", "-i", "tron-legacy-2010"}, path: func(string) string { return "/movies/tron-legacy-2010/sentiments" }},
		{name: "movies watchnow", cmd: MoviesCmd, args: []string{"-a", "watchnow", "-i", "tron-legacy-2010", "-country", "us"}, path: func(string) string { return "/movies/tron-legacy-2010/watchnow/us" }},
		{name: "movies watchnow default country is rejected", cmd: MoviesCmd, args: []string{"-a", "watchnow", "-i", "tron-legacy-2010"}, path: func(string) string { return "" }},
		{name: "movies justwatch_links", cmd: MoviesCmd, args: []string{"-a", "justwatch_links", "-i", "tron-legacy-2010", "-country", "pl"}, path: func(string) string { return "/movies/tron-legacy-2010/watchnow/justwatch_links/pl" }},
		{name: "shows -period", cmd: ShowsCmd, args: []string{"-a", "favorited", "-period", "daily"}, path: func(string) string { return "/shows/favorited/daily" }},
		{name: "shows -start_date", cmd: ShowsCmd, args: []string{"-a", "updates", "-start_date", "2026-01-15"}, path: func(tz string) string { return "/shows/updates/" + startDate(tz) }},
		{name: "shows default start", cmd: ShowsCmd, args: []string{"-a", "updates"}, path: func(tz string) string { return "/shows/updates/" + window(tz) }},
		{name: "people -start_date", cmd: PeopleCmd, args: []string{"-a", "updates", "-start_date", "2026-01-15"}, path: func(tz string) string { return "/people/updates/" + startDate(tz) }},
		{name: "people default start", cmd: PeopleCmd, args: []string{"-a", "updates"}, path: func(tz string) string { return "/people/updates/" + window(tz) }},
		{name: "calendars -start_date", cmd: CalendarsCmd, args: []string{"-a", "all-shows", "-start_date", "2026-01-15", "-days", "3"}, path: func(string) string { return "/calendars/all/shows/2026-01-15/3" }},
		{name: "calendars default", cmd: CalendarsCmd, args: []string{"-a", "all-shows"}, path: func(string) string { return "/calendars/all/shows/" + time.Now().Format("2006-01-02") + "/7" }},
		{name: "sync default window", cmd: SyncCmd, args: []string{"-a", "get_history", "-t", "movies"}, path: func(string) string { return "/sync/history/movies" }},
		{name: "users history default window", cmd: UsersCmd, args: []string{"-a", "history", "-u", "sean"}, path: func(string) string { return "/users/sean/history/movies" }},
	}

	// sync and users ignore -o (their updaters recompute Output), so run in a temp dir
	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.Chdir(t.TempDir()))
	t.Cleanup(func() { _ = os.Chdir(wd) })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			setup := internal.Setup()
			defer setup.Teardown()
			var gotPath, gotStartAt string
			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotStartAt = r.URL.Path, r.URL.Query().Get("start_at")
				_, _ = w.Write([]byte(`[]`))
			})

			// -o keeps handler output files out of the working tree
			args := append([]string{"-o", filepath.Join(t.TempDir(), "out.json")}, tt.args...)
			_ = tt.cmd.Exec(fs, setup.Client, config, args)
			tz := tt.cmd.Options.Timezone
			assert.Equal(t, tt.path(tz), gotPath)
			if tt.cmd == SyncCmd || tt.cmd == UsersCmd {
				// history keeps its bounded default window instead of the whole history
				assert.Equal(t, window(tz), gotStartAt)
			}
		})
	}
}
