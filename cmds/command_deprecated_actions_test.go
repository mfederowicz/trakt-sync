// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/deprecation"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestDeprecatedActionsPrintNoticeAndRun checks that an action whose endpoint Trakt dropped from the contract
// prints its deprecation notice and still runs, and that the same action name in another module prints nothing.
func TestDeprecatedActionsPrintNoticeAndRun(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/deprecated/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/deprecated/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/deprecated/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/deprecated/token.json", "/deprecated/user_settings.json"

	const status = " is deprecated since 2026-10-13: Trakt dropped this endpoint from its API contract, it may stop working\n"
	tests := []struct {
		name    string
		command *Command
		args    []string
		want    string
	}{
		{name: "people lists", command: PeopleCmd, args: []string{"-a", "lists", "-i", "bryan-cranston"},
			want: "people -a lists is deprecated since 2026-10-20: Trakt deprecated this endpoint: lists no longer hold people, so it always returns no lists\n"},
		{name: "sync update_favorites", command: SyncCmd, args: []string{"-a", "update_favorites", "-sort_by", "added"}, want: "sync -a update_favorites" + status},
		{name: "younify refresh", command: YounifyCmd, args: []string{"-a", "refresh", "-service_id", "netflix"}},
		{name: "movies refresh_justwatch", command: MoviesCmd, args: []string{"-a", "refresh_justwatch", "-i", "tron-legacy-2010"}},
		{name: "users saved_filters", command: UsersCmd, args: []string{"-a", "saved_filters"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			ran := false
			command := &Command{Name: tt.command.Name, Flag: tt.command.Flag, Run: func(_ *Command, _ ...string) error {
				ran = true
				return nil
			}}
			var err error
			out := captureStdout(t, func() { err = command.Exec(fs, trakt.NewClient(nil), config, tt.args) })
			assert.NoError(t, err)
			assert.True(t, ran)
			assert.Equal(t, tt.want, out)
		})
	}
}

// TestDeprecationResult checks what Exec does with a notice in each stage: a deprecated action goes on after
// its notice is printed, a locked or removed one ends with the notice as its error.
func TestDeprecationResult(t *testing.T) {
	notice := deprecation.Notice{Module: consts.Movies, Action: consts.Refresh, Reason: "Trakt removed this endpoint", Deprecated: "2026-10-13"}
	var err error
	out := captureStdout(t, func() { err = deprecationResult(notice) })
	assert.NoError(t, err)
	assert.Equal(t, "movies -a refresh is deprecated since 2026-10-13: Trakt removed this endpoint\n", out)

	notice.Locked = "2026-11-01"
	out = captureStdout(t, func() { err = deprecationResult(notice) })
	assert.EqualError(t, err, "movies -a refresh is locked since 2026-11-01: Trakt removed this endpoint; no request was sent")
	assert.Empty(t, out)

	notice.Removed = "2026-12-01"
	assert.EqualError(t, deprecationResult(notice), "movies -a refresh is removed since 2026-12-01: Trakt removed this endpoint; no request was sent")
}

// TestLockedActionsSendNothing checks that an action whose endpoint is locked ends with its notice as the error,
// prints nothing and never runs.
func TestLockedActionsSendNothing(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/locked/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/locked/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/locked/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/locked/token.json", "/locked/user_settings.json"

	const status = ": Trakt removed this endpoint from its API, the live API no longer serves it; no request was sent"
	const gone = " is locked since 2026-11-20" + status
	const retired = " is locked since 2026-11-20: Trakt retired saved filters, the live API answers 410 saved_filters_retired; no request was sent"
	tests := []struct {
		name string
		cmd  *Command
		args []string
		want string
	}{
		{name: "movies hot", cmd: MoviesCmd, args: []string{"-a", "hot"}, want: "movies -a hot is locked since 2026-11-17" + status},
		{name: "movies streaming", cmd: MoviesCmd, args: []string{"-a", "streaming", "-period", "daily"}, want: "movies -a streaming is locked since 2026-11-17" + status},
		{name: "users plex_connect", cmd: UsersCmd, args: []string{"-a", "plex_connect", "-return_url", "http://localhost:8080"}, want: "users -a plex_connect is locked since 2026-11-18" + status},
		{name: "users plex_servers", cmd: UsersCmd, args: []string{"-a", "plex_servers"}, want: "users -a plex_servers is locked since 2026-11-18" + status},
		{name: "users plex_server", cmd: UsersCmd, args: []string{"-a", "plex_server", "-i", "abc"}, want: "users -a plex_server is locked since 2026-11-18" + status},
		{name: "users plex_sync", cmd: UsersCmd, args: []string{"-a", "plex_sync", "-i", "abc", "-all_data"}, want: "users -a plex_sync is locked since 2026-11-18" + status},
		{name: "movies refresh", cmd: MoviesCmd, args: []string{"-a", "refresh", "-i", "tron-legacy-2010"}, want: "movies -a refresh" + gone},
		{name: "shows refresh", cmd: ShowsCmd, args: []string{"-a", "refresh", "-i", "the-sopranos"}, want: "shows -a refresh" + gone},
		{name: "people refresh", cmd: PeopleCmd, args: []string{"-a", "refresh", "-i", "bryan-cranston"}, want: "people -a refresh" + gone},
		{name: "search remove_recent", cmd: SearchCmd, args: []string{"-a", "remove_recent", "-q", "reacher", "-i", "139606"}, want: "search -a remove_recent" + gone},
		{name: "users add_saved_filters", cmd: UsersCmd, args: []string{"-a", "add_saved_filters"}, want: "users -a add_saved_filters" + retired},
		{name: "users delete_saved_filter", cmd: UsersCmd, args: []string{"-a", "delete_saved_filter", "-i", "101"}, want: "users -a delete_saved_filter" + retired},
		{name: "social_recommendations movies", cmd: SocialRecommendationsCmd, args: []string{"-a", "movies"}, want: "social_recommendations -a movies" + gone},
		{name: "social_recommendations shows", cmd: SocialRecommendationsCmd, args: []string{"-a", "shows", "-watch_window", "30"}, want: "social_recommendations -a shows" + gone},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			ran := false
			command := &Command{Name: tt.cmd.Name, Flag: tt.cmd.Flag, Run: func(_ *Command, _ ...string) error {
				ran = true
				return nil
			}}
			var err error
			out := captureStdout(t, func() { err = command.Exec(fs, trakt.NewClient(nil), config, tt.args) })
			assert.EqualError(t, err, tt.want)
			assert.False(t, ran)
			assert.Empty(t, out)
		})
	}
}
