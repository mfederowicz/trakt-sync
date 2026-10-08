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

func TestSmartListsFlags(t *testing.T) {
	resetAllFlags()
	t.Cleanup(resetAllFlags)

	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-smart-lists/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	var got *str.Options
	command := &Command{Name: consts.SmartLists, Flag: SmartListsCmd.Flag, Run: func(c *Command, _ ...string) error {
		got = c.UpdateOptionsWithCommandFlags(c.Options)
		return nil
	}}
	args := []string{"-a", "items", "-i", "top-sci-fi", "-watchnow", "free", "-genres", "action", "-subgenres", "space", "-years", "2020-2026",
		"-ratings", "75-100", "-runtimes", "90-150", "-countries", "us", "-certifications", "pg-13", "-ignore_watched", "true", "-ignore_watchlisted", "false",
		"-watchnow_country", "pl", "-parental_nudity", "0-1", "-parental_violence", "0-2", "-parental_profanity", "0-1", "-parental_alcohol", "0-3",
		"-parental_frightening", "1-2", "-parental_include_unrated", "-t", "shows", "-sort_by", "imdb_rating", "-sort_how", "desc"}
	assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, args))
	assert.Equal(t, "items", got.Action)
	assert.Equal(t, "top-sci-fi", got.InternalID)
	assert.Equal(t, "free", got.WatchNow)
	assert.Equal(t, "action", got.Genres)
	assert.Equal(t, "space", got.Subgenres)
	assert.Equal(t, "2020-2026", got.Years)
	assert.Equal(t, "75-100", got.Ratings)
	assert.Equal(t, "90-150", got.Runtimes)
	assert.Equal(t, "us", got.Countries)
	assert.Equal(t, "pg-13", got.Certifications)
	assert.Equal(t, "true", got.IgnoreWatched)
	assert.Equal(t, "false", got.IgnoreWatchlisted)
	assert.Equal(t, "pl", got.WatchNowCountry)
	assert.Equal(t, str.ParentalGuide{Nudity: "0-1", Violence: "0-2", Profanity: "0-1", Alcohol: "0-3", Frightening: "1-2", IncludeUnrated: true}, got.Parental)
	assert.Equal(t, "shows", got.Type)
	assert.Equal(t, "imdb_rating", got.SortBy)
	assert.Equal(t, "desc", got.SortHow)
	assert.Equal(t, "export_smart_lists_items_top-sci-fi.json", got.Output)
}

// Without -t, -sort_by and -sort_how the config and built-in defaults (movies, rank, asc) must not reach the handler,
// or every run would use the typed and sorted route.
func TestSmartListsItemsSortKeepsOnlyGivenFlags(t *testing.T) {
	resetAllFlags()
	t.Cleanup(resetAllFlags)

	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-smart-lists-sort/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

	tests := []struct {
		name                    string
		args                    []string
		wantType, wantBy, wantH string
	}{
		{name: "no flags", args: []string{"-a", "items", "-i", "top-sci-fi"}},
		{name: "type only", args: []string{"-a", "items", "-i", "top-sci-fi", "-t", "movies"}, wantType: "movies"},
		{name: "sort_by only", args: []string{"-a", "items", "-i", "top-sci-fi", "-sort_by", "rank"}, wantBy: "rank"},
		{name: "sort_how only", args: []string{"-a", "items", "-i", "top-sci-fi", "-sort_how", "asc"}, wantH: "asc"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			fileConfig := cfg.DefaultConfig()
			fileConfig.ClientID = "a"
			fileConfig.ClientSecret = "b"
			fileConfig.TokenPath = tmpPath + "token.json"
			fileConfig.SettingsPath = tmpPath + "user_settings.json"

			var got *str.Options
			command := &Command{Name: consts.SmartLists, Flag: SmartListsCmd.Flag, Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				smartListsItemsSort(got, c.flagIsSet("t"), c.flagIsSet("sort_by"), c.flagIsSet("sort_how"))
				return nil
			}}
			assert.NoError(t, command.Exec(fs, trakt.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.wantType, got.Type)
			assert.Equal(t, tt.wantBy, got.SortBy)
			assert.Equal(t, tt.wantH, got.SortHow)
		})
	}
}
