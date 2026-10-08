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

func TestShowsMediaFilterFlags(t *testing.T) {
	resetAllFlags()
	t.Cleanup(resetAllFlags)

	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-shows-filters/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	var got *str.Options
	command := &Command{Name: consts.Shows, Flag: ShowsCmd.Flag, Run: func(c *Command, _ ...string) error {
		got = c.UpdateOptionsWithCommandFlags(c.Options)
		return nil
	}}
	args := []string{"-a", "trending", "-watchnow", "free", "-genres", "action", "-subgenres", "space", "-years", "2020-2026",
		"-ratings", "75-100", "-runtimes", "90-150", "-countries", "us", "-certifications", "tv-14", "-status", "returning series,ended", "-start_date", "2026-01-01", "-end_date", "2026-12-31",
		"-languages", "en,pl", "-imdb_ratings", "8.5-10.0", "-rt_meters", "90-100", "-rt_user_meters", "80-100"}
	assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, args))
	assert.Equal(t, "trending", got.Action)
	assert.Equal(t, "free", got.WatchNow)
	assert.Equal(t, "action", got.Genres)
	assert.Equal(t, "space", got.Subgenres)
	assert.Equal(t, "2020-2026", got.Years)
	assert.Equal(t, "75-100", got.Ratings)
	assert.Equal(t, "90-150", got.Runtimes)
	assert.Equal(t, "us", got.Countries)
	assert.Equal(t, "tv-14", got.Certifications)
	assert.Equal(t, "returning series,ended", got.ShowStatus)
	assert.Equal(t, "2026-01-01", got.MediaStartDate)
	assert.Equal(t, "2026-12-31", got.MediaEndDate)
	assert.Equal(t, "en,pl", got.Languages)
	assert.Equal(t, "8.5-10.0", got.ImdbRatings)
	assert.Equal(t, "90-100", got.RtMeters)
	assert.Equal(t, "80-100", got.RtUserMeters)
}

// TestMediaFilterFlagsAreAvailable checks that movies and shows register every media filter flag that is not global.
func TestMediaFilterFlagsAreAvailable(t *testing.T) {
	filters := []string{"watchnow", "subgenres", "ratings", "certifications", "start_date", "end_date", "genres", "years", "countries", "runtimes", "languages"}
	filters = append(filters, "imdb_ratings", "rt_meters", "rt_user_meters")
	filters = append(filters, "status")
	modules := map[*Command][]string{MoviesCmd: filters, ShowsCmd: filters}
	global := []string{"genres", "years", "countries", "runtimes", "languages"}
	for command, names := range modules {
		for _, name := range names {
			if !str.ContainString(name, global) {
				assert.NotNil(t, command.Flag.Lookup(name), "%s has no -%s flag", command.Name, name)
			}
		}
	}
}
