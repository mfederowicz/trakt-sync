// Package cmds used for commands modules
package cmds

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// TestUsersTypedRoutes checks that the typed ratings routes of the contract are reached through users -a ratings -t.
func TestUsersTypedRoutes(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/usr-typed/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/usr-typed/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/usr-typed/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/usr-typed/token.json", "/usr-typed/user_settings.json"

	tests := []struct {
		name     string
		args     []string
		wantPath string
	}{
		{name: "ratings movies", args: []string{"-a", "ratings", "-t", "movies"}, wantPath: "/users/sean/ratings/movies"},
		{name: "ratings shows", args: []string{"-a", "ratings", "-t", "shows"}, wantPath: "/users/sean/ratings/shows"},
		{name: "ratings episodes", args: []string{"-a", "ratings", "-t", "episodes"}, wantPath: "/users/sean/ratings/episodes"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			setup := internal.Setup()
			defer setup.Teardown()
			var gotPaths []string
			setup.Mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
				gotPaths = append(gotPaths, r.URL.Path)
				_, _ = w.Write([]byte(`[]`))
			})

			args := append([]string{"-o", filepath.Join(t.TempDir(), "out.json"), "-u", "sean"}, tt.args...)
			_ = UsersCmd.Exec(fs, setup.Client, config, args) // an empty list may end in an "empty" error; only the request matters here
			assert.Contains(t, gotPaths, tt.wantPath)
		})
	}
}

// TestUsersListItemsTypeFlag checks that users -a lists passes -t unchanged, so the typed list items routes
// (items/movie, items/show, items/movie,show, items/movie,show,season,episode) are reachable;
// the handler is not run because it writes the lists overview to a fixed file in the working directory.
func TestUsersListItemsTypeFlag(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/usr-typed/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/usr-typed/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/usr-typed/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/usr-typed/token.json", "/usr-typed/user_settings.json"

	for _, typ := range []string{"movie", "show", "movie,show", "movie,show,season,episode"} {
		typ := typ
		t.Run(typ, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			var gotType, gotID string
			command := &Command{Name: consts.Users, Flag: UsersCmd.Flag, Run: func(c *Command, _ ...string) error {
				got := c.UpdateOptionsWithCommandFlags(c.Options)
				gotType, gotID = got.Type, got.ID
				return nil
			}}
			assert.NoError(t, command.Exec(fs, internal.NewClient(nil), config, []string{"-a", "lists", "-i", "55", "-t", typ}))
			assert.Equal(t, typ, gotType)
			assert.Equal(t, "55", gotID)
		})
	}
}

// TestUsersSortPathRoutes checks that -sort switches users -a watchlist|favorites to the /{type}/{sort} routes.
func TestUsersSortPathRoutes(t *testing.T) {
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/usr-sort/", consts.X755))
	assert.NoError(t, afero.WriteFile(fs, "/usr-sort/token.json", []byte("{}"), consts.X644))
	assert.NoError(t, afero.WriteFile(fs, "/usr-sort/user_settings.json", []byte(`{"user":{"username":"sean"}}`), consts.X644))
	config := cfg.DefaultConfig()
	config.ClientID, config.ClientSecret = "a", "b"
	config.TokenPath, config.SettingsPath = "/usr-sort/token.json", "/usr-sort/user_settings.json"

	tests := []struct {
		name     string
		args     []string
		wantPath string
		wantErr  string
	}{
		{name: "watchlist all", args: []string{"-a", "watchlist", "-sort", "added"}, wantPath: "/users/sean/watchlist/movie,show/added"},
		{name: "watchlist movies", args: []string{"-a", "watchlist", "-t", "movies", "-sort", "rank"}, wantPath: "/users/sean/watchlist/movies/rank"},
		{name: "watchlist shows", args: []string{"-a", "watchlist", "-t", "shows", "-sort", "title"}, wantPath: "/users/sean/watchlist/shows/title"},
		{name: "favorites all", args: []string{"-a", "favorites", "-sort", "released"}, wantPath: "/users/sean/favorites/media/released"},
		{name: "favorites movies", args: []string{"-a", "favorites", "-t", "movies", "-sort", "votes"}, wantPath: "/users/sean/favorites/movies/votes"},
		{name: "favorites shows", args: []string{"-a", "favorites", "-t", "shows", "-sort", "popularity"}, wantPath: "/users/sean/favorites/shows/popularity"},
		{name: "watchlist without -sort", args: []string{"-a", "watchlist", "-t", "movies"}, wantPath: "/users/sean/watchlist/movies/rank/asc"},
		{name: "invalid sort", args: []string{"-a", "watchlist", "-sort", "imdb_rating"}, wantErr: "sort 'imdb_rating' is not valid"},
		{name: "unsupported type", args: []string{"-a", "favorites", "-t", "episodes", "-sort", "rank"}, wantErr: "-sort works with -t all, movies or shows"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resetAllFlags()
			t.Cleanup(resetAllFlags)

			setup := internal.Setup()
			defer setup.Teardown()
			var gotPaths []string
			setup.Mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
				gotPaths = append(gotPaths, r.URL.Path)
				_, _ = w.Write([]byte(`[{"rank":1,"type":"movie","movie":{"title":"Arrival"}}]`))
			})

			args := append([]string{"-o", filepath.Join(t.TempDir(), "out.json"), "-u", "sean"}, tt.args...)
			err := UsersCmd.Exec(fs, setup.Client, config, args)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				assert.Empty(t, gotPaths, "no request with invalid -sort / -t")
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, []string{tt.wantPath}, gotPaths)
		})
	}
}
