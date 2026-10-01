// Package cmds used for commands modules
package cmds

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/trakt/trakttest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestCommentsSpoilerFlag checks comments -spoiler reaches the posted reply (it was ignored before).
func TestCommentsSpoilerFlag(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/spoiler/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/spoiler/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/spoiler/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/spoiler/token.json", "/spoiler/user_settings.json"

	for _, spoiler := range []bool{true, false} {
		spoiler := spoiler
		name := "without -spoiler"
		if spoiler {
			name = "with -spoiler"
		}
		t.Run(name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			setup := trakttest.Setup()
			defer setup.Teardown()
			var sent struct {
				Spoiler *bool `json:"spoiler"`
			}
			calls := 0
			setup.Mux.HandleFunc("/comments/7/replies", func(w http.ResponseWriter, r *http.Request) {
				calls++
				assert.Equal(t, http.MethodPost, r.Method)
				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(body, &sent))
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":8}`))
			})

			args := []string{"-o", filepath.Join(t.TempDir(), "out.json"), "-a", "replies", "-comment_id", "7",
				"-reply", "this reply has more than five words"}
			if spoiler {
				args = append(args, "-spoiler")
			}
			captureStdout(t, func() {
				assert.NoError(t, CommentsCmd.Exec(fs, setup.Client, config, args))
			})
			assert.Equal(t, 1, calls)
			if assert.NotNil(t, sent.Spoiler) {
				assert.Equal(t, spoiler, *sent.Spoiler)
			}
		})
	}
}

// TestCommentsUnknownType checks comments -a comments with a missing or unknown -t ends with an error (it panicked before).
func TestCommentsUnknownType(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/unknown/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/unknown/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/unknown/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/unknown/token.json", "/unknown/user_settings.json"

	cases := map[string]struct {
		args []string
		want string
	}{
		"without -t":       {want: `comments/comments: unknown type "movies"`},
		"plural -t movies": {args: []string{"-t", "movies"}, want: `comments/comments: unknown type "movies"`},
		"unknown -t":       {args: []string{"-t", "bogus"}, want: `comments/comments: unknown type "bogus"`},
	}
	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			setup := trakttest.Setup()
			defer setup.Teardown()
			calls := 0
			setup.Mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {
				calls++
			})

			args := append([]string{"-a", "comments", "-i", "55", "-comment", "this comment has more than five words"}, tc.args...)
			captureStdout(t, func() {
				assert.EqualError(t, CommentsCmd.Exec(fs, setup.Client, config, args), tc.want)
			})
			assert.Zero(t, calls)
		})
	}
}
