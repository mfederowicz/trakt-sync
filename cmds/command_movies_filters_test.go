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

func TestMoviesMediaFilterFlags(t *testing.T) {
	resetAllFlags()
	t.Cleanup(resetAllFlags)

	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-movies-filters/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	var got *str.Options
	command := &Command{Name: consts.Movies, Flag: MoviesCmd.Flag, Run: func(c *Command, _ ...string) error {
		got = c.UpdateOptionsWithCommandFlags(c.Options)
		return nil
	}}
	args := []string{"-a", "trending", "-watchnow", "free", "-genres", "action", "-subgenres", "space", "-years", "2020-2026",
		"-ratings", "75-100", "-runtimes", "90-150", "-countries", "us", "-certifications", "pg-13", "-start_date", "2026-01-01", "-end_date", "2026-12-31"}
	assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, args))
	assert.Equal(t, "trending", got.Action)
	assert.Equal(t, "free", got.WatchNow)
	assert.Equal(t, "action", got.Genres)
	assert.Equal(t, "space", got.Subgenres)
	assert.Equal(t, "2020-2026", got.Years)
	assert.Equal(t, "75-100", got.Ratings)
	assert.Equal(t, "90-150", got.Runtimes)
	assert.Equal(t, "us", got.Countries)
	assert.Equal(t, "pg-13", got.Certifications)
	assert.Equal(t, "2026-01-01", got.MediaStartDate)
	assert.Equal(t, "2026-12-31", got.MediaEndDate)
}
