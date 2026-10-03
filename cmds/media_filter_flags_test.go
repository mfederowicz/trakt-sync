// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// modules that register their filter flags with newMediaFilterFlags copy them, and the global filter flags, to the options.
func TestModuleMediaFilterFlags(t *testing.T) {
	modules := []struct {
		command *Command
		action  string
		dates   bool // the module has the -start_date and -end_date filters
	}{
		{command: MediaCmd, action: "trending", dates: true},
		{command: RecommendationsCmd, action: "movies", dates: true},
		{command: ListsCmd, action: "trending", dates: true},
		{command: UsersCmd, action: "watchlist", dates: true},
		{command: CalendarsCmd, action: "all_shows"},
	}
	for _, m := range modules {
		m := m
		t.Run(m.command.Name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			fs := afero.NewMemMapFs()
			tmpPath := "/tmp-" + m.command.Name + "-filters/"
			assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
			assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
			assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

			fileConfig := cfg.DefaultConfig()
			fileConfig.ClientID = "a"
			fileConfig.ClientSecret = "b"
			fileConfig.TokenPath = tmpPath + "token.json"
			fileConfig.SettingsPath = tmpPath + "user_settings.json"

			var got *str.Options
			command := &Command{Name: m.command.Name, Flag: m.command.Flag, Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				return nil
			}}
			args := []string{"-a", m.action, "-watchnow", "free", "-genres", "action", "-subgenres", "space", "-years", "2020-2026", "-ratings", "75-100",
				"-runtimes", "90-150", "-countries", "us", "-certifications", "pg-13", "-languages", "en,pl"}
			want := str.Options{WatchNow: "free", Genres: "action", Subgenres: "space", Years: "2020-2026", Ratings: "75-100", Runtimes: "90-150",
				Countries: "us", Certifications: "pg-13", Languages: "en,pl"}
			if m.dates {
				args = append(args, "-start_date", "2026-01-01", "-end_date", "2026-12-31")
				want.MediaStartDate = "2026-01-01"
				want.MediaEndDate = "2026-12-31"
			}
			assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, args))
			assert.Equal(t, m.action, got.Action)
			assert.Equal(t, want, str.Options{WatchNow: got.WatchNow, Genres: got.Genres, Subgenres: got.Subgenres, Years: got.Years, Ratings: got.Ratings,
				Runtimes: got.Runtimes, Countries: got.Countries, Certifications: got.Certifications, MediaStartDate: got.MediaStartDate, MediaEndDate: got.MediaEndDate, Languages: got.Languages})
		})
	}
}
