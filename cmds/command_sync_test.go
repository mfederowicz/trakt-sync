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
