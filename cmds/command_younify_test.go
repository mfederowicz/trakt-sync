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

func TestYounifyFlags(t *testing.T) {
	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-younify/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte("{}"), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	tests := []struct {
		name       string
		args       []string
		wantAction string
		wantID     string
		wantURL    string
		wantAll    bool
		wantOutput string
	}{
		{name: "connections", args: []string{"-a", "connections"}, wantAction: consts.Connections, wantURL: consts.DefaultReturnURL,
			wantOutput: "export_younify_connections.json"},
		{name: "connect", args: []string{"-a", "connect", "-service_id", "netflix", "-return_url", "trakt://settings"}, wantAction: consts.Connect,
			wantID: "netflix", wantURL: "trakt://settings", wantOutput: "export_younify_connect_results.json"},
		{name: "refresh all data", args: []string{"-a", "refresh", "-service_id", "netflix", "-all_data"}, wantAction: consts.Refresh,
			wantID: "netflix", wantURL: consts.DefaultReturnURL, wantAll: true, wantOutput: "export_younify_refresh.json"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			var got *str.Options
			command := &Command{Name: consts.Younify, Flag: YounifyCmd.Flag, Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				return nil
			}}
			assert.NoError(t, command.Exec(fs, internal.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.wantAction, got.Action)
			assert.Equal(t, tt.wantID, got.ServiceID)
			assert.Equal(t, tt.wantURL, got.ReturnURL)
			assert.Equal(t, tt.wantAll, got.AllData)
			assert.Equal(t, tt.wantOutput, got.Output)
		})
	}
}
