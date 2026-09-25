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

func TestWatchNowFlags(t *testing.T) {
	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-watchnow/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	tests := []struct {
		name        string
		args        []string
		wantCountry string
		wantOutput  string
	}{
		{name: "all countries", args: []string{"-a", "sources"}, wantOutput: "export_watchnow_sources.json"},
		{name: "one country", args: []string{"-a", "sources", "-country", "us"}, wantCountry: "us", wantOutput: "export_watchnow_sources_us.json"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			var got *str.Options
			command := &Command{Name: consts.WatchNow, Flag: WatchNowCmd.Flag, Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				return nil
			}}
			assert.NoError(t, command.Exec(fs, internal.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, consts.Sources, got.Action)
			assert.Equal(t, tt.wantCountry, got.Country)
			assert.Equal(t, tt.wantOutput, got.Output)
		})
	}
}
