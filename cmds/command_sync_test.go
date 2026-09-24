// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSyncPlaybackType(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		typ     string
		typeSet bool
		want    string
	}{
		{name: "playback without -t is all", action: consts.Playback, typ: cfg.DefaultConfig().Type, want: consts.ActionTypeAll},
		{name: "playback -t movies", action: consts.Playback, typ: "movies", typeSet: true, want: "movies"},
		{name: "playback -t episodes", action: consts.Playback, typ: "episodes", typeSet: true, want: "episodes"},
		{name: "playback -t all", action: consts.Playback, typ: consts.ActionTypeAll, typeSet: true, want: consts.ActionTypeAll},
		{name: "other action keeps default type", action: consts.GetWatched, typ: cfg.DefaultConfig().Type, want: cfg.DefaultConfig().Type},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := syncPlaybackType(&str.Options{Action: tt.action, Type: tt.typ}, tt.typeSet)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSyncPlaybackTypeFromFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "no -t", want: consts.ActionTypeAll},
		{name: "-t episodes", args: []string{"-t", "episodes"}, want: "episodes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// -t shares its value with the global flag, so restore it for other tests
			t.Cleanup(func() { *_strType = cfg.DefaultConfig().Type })

			fs := afero.NewMemMapFs()
			tmpPath := "/tmp-playback/"
			assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
			assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
			assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

			fileConfig := cfg.DefaultConfig()
			fileConfig.ClientID = "a"
			fileConfig.ClientSecret = "b"
			fileConfig.TokenPath = tmpPath + "token.json"
			fileConfig.SettingsPath = tmpPath + "user_settings.json"

			var got string
			command := &Command{Name: "watchlist", Run: func(c *Command, _ ...string) error {
				got = syncPlaybackType(&str.Options{Action: consts.Playback, Type: c.Options.Type}, c.flagIsSet("t"))
				return nil
			}}
			assert.NoError(t, command.Exec(fs, internal.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSyncProgressSort(t *testing.T) {
	tests := []struct {
		name       string
		action     string
		sortBySet  bool
		sortHowSet bool
		wantBy     string
		wantHow    string
	}{
		{name: "up next drops default sort", action: consts.GetUpNext, wantBy: "", wantHow: ""},
		{name: "watched progress keeps set sort", action: consts.GetWatchedProgress, sortBySet: true, sortHowSet: true, wantBy: "rank", wantHow: "asc"},
		{name: "up next keeps only sort_how", action: consts.GetUpNext, sortHowSet: true, wantBy: "", wantHow: "asc"},
		{name: "other action untouched", action: consts.GetWatchlist, wantBy: "rank", wantHow: "asc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &str.Options{Action: tt.action, SortBy: cfg.DefaultConfig().SortBy, SortHow: cfg.DefaultConfig().SortHow}
			syncProgressSort(options, tt.sortBySet, tt.sortHowSet)
			assert.Equal(t, tt.wantBy, options.SortBy)
			assert.Equal(t, tt.wantHow, options.SortHow)
		})
	}
}

func TestSyncProgressSortFromFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantBy  string
		wantHow string
	}{
		{name: "no sort flags", wantBy: "", wantHow: ""},
		{name: "-sort_by after the module", args: []string{"-sort_by", "added"}, wantBy: "added", wantHow: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(resetAllFlags)

			fs := afero.NewMemMapFs()
			tmpPath := "/tmp-progress/"
			assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
			assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
			assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

			fileConfig := cfg.DefaultConfig()
			fileConfig.ClientID = "a"
			fileConfig.ClientSecret = "b"
			fileConfig.TokenPath = tmpPath + "token.json"
			fileConfig.SettingsPath = tmpPath + "user_settings.json"

			var got *str.Options
			command := &Command{Name: "watchlist", Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				got.Action = consts.GetUpNext
				syncProgressSort(got, c.flagIsSet("sort_by"), c.flagIsSet("sort_how"))
				return nil
			}}
			assert.NoError(t, command.Exec(fs, internal.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.wantBy, got.SortBy)
			assert.Equal(t, tt.wantHow, got.SortHow)
		})
	}
}
