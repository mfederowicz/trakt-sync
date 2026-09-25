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

func TestUsersReviewAndActivityFlags(t *testing.T) {
	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-users-reviews/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	tests := []struct {
		name       string
		args       []string
		wantYear   int
		wantMonth  int
		wantType   string
		wantOutput string
	}{
		{name: "month in review", args: []string{"-a", "month_in_review", "-year", "2026", "-month", "8"}, wantYear: 2026, wantMonth: 8,
			wantOutput: "export_users_month_in_review_2026-08.json"},
		{name: "year in review", args: []string{"-a", "year_in_review", "-year", "2025"}, wantYear: 2025,
			wantOutput: "export_users_year_in_review_2025.json"},
		{name: "activities", args: []string{"-a", "activities", "-t", "following"}, wantType: "following",
			wantOutput: "export_users_activities_following.json"},
		{name: "comment reactions", args: []string{"-a", "comment_reactions"}, wantOutput: "export_users_comment_reactions.json"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			var got *str.Options
			command := &Command{Name: consts.Users, Flag: UsersCmd.Flag, Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				return nil
			}}
			assert.NoError(t, command.Exec(fs, internal.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.wantYear, got.Year)
			assert.Equal(t, tt.wantMonth, got.Month)
			assert.Equal(t, tt.wantType, got.Type)
			assert.Equal(t, tt.wantOutput, got.Output)
		})
	}
}

func TestUsersPlexFlags(t *testing.T) {
	fs := afero.NewMemMapFs()
	tmpPath := "/tmp-users-plex/"
	assert.NoError(t, fs.MkdirAll(tmpPath, consts.X755))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, tmpPath+"user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))

	fileConfig := cfg.DefaultConfig()
	fileConfig.ClientID = "a"
	fileConfig.ClientSecret = "b"
	fileConfig.TokenPath = tmpPath + "token.json"
	fileConfig.SettingsPath = tmpPath + "user_settings.json"

	tests := []struct {
		name    string
		args    []string
		wantURL string
		wantAll bool
		wantID  string
	}{
		{name: "connect default", args: []string{"-a", "plex_connect"}, wantURL: consts.DefaultReturnURL},
		{name: "connect own url", args: []string{"-a", "plex_connect", "-return_url", "http://localhost:8080"}, wantURL: "http://localhost:8080"},
		{name: "sync full", args: []string{"-a", "plex_sync", "-i", "abc", "-all_data"}, wantURL: consts.DefaultReturnURL, wantAll: true, wantID: "abc"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			var got *str.Options
			command := &Command{Name: consts.Users, Flag: UsersCmd.Flag, Run: func(c *Command, _ ...string) error {
				got = c.UpdateOptionsWithCommandFlags(c.Options)
				return nil
			}}
			assert.NoError(t, command.Exec(fs, internal.NewClient(nil), fileConfig, tt.args))
			assert.Equal(t, tt.wantURL, got.ReturnURL)
			assert.Equal(t, tt.wantAll, got.AllData)
			assert.Equal(t, tt.wantID, got.ID)
		})
	}
}
