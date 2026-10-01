package cfg

import (
	"flag"
	"os"
	"os/user"
	"path/filepath"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const (
	testTokenPath    = "/home/tester/token.json"
	testSettingsPath = "/home/tester/user_settings.json"
	testTokenJSON    = `{"access_token":"access","token_type":"bearer","refresh_token":"refresh","scope":"public","expires_in":7200,"created_at":1790863634}`
	testSettingsJSON = `{"user":{"username":"sean"},"account":{"timezone":"Europe/Warsaw"}}`
)

// credentialsFs returns an in-memory filesystem holding a token file and a user settings file.
func credentialsFs(t *testing.T) afero.Fs {
	t.Helper()
	fs := afero.NewMemMapFs()
	assert.NoError(t, afero.WriteFile(fs, testTokenPath, []byte(testTokenJSON), consts.X600))
	assert.NoError(t, afero.WriteFile(fs, testSettingsPath, []byte(testSettingsJSON), consts.X600))
	return fs
}

// validConfig is the default config with the fields a user has to set.
func validConfig() *Config {
	config := DefaultConfig()
	config.ClientID = "client-id"
	config.ClientSecret = "client-secret"
	config.TokenPath = testTokenPath
	config.SettingsPath = testSettingsPath
	return config
}

// useFlags replaces the process flags with the given ones for one test.
func useFlags(t *testing.T, args ...string) map[string]string {
	t.Helper()
	oldCommandLine, oldArgs := flag.CommandLine, os.Args
	t.Cleanup(func() {
		flag.CommandLine, os.Args = oldCommandLine, oldArgs
	})

	flag.CommandLine = emptyFlagset()
	for _, name := range []string{"a", "c", "f", "i", "l", "m", "o", "s", "t", "u"} {
		flag.String(name, consts.EmptyString, name)
	}
	flag.Bool("v", false, "v")
	os.Args = append([]string{consts.CMD}, args...)
	parseFlags()

	flagMap := map[string]string{}
	flag.VisitAll(func(f *flag.Flag) {
		flagMap[f.Name] = f.Value.String()
	})
	return flagMap
}

func TestValidateConfig(t *testing.T) {
	assert.True(t, ValidateConfig("watchlist", OptionsConfig{}))
	assert.True(t, ValidateConfig("watchlist", OptionsConfig{Type: []string{"movies", "shows"}, Sort: []string{"rank"}, Format: []string{"imdb"}}))
	assert.False(t, ValidateConfig("watchlist", OptionsConfig{Type: []string{"movies", "books"}}))
	assert.False(t, ValidateConfig("watchlist", OptionsConfig{Sort: []string{"random"}}))
	assert.False(t, ValidateConfig("watchlist", OptionsConfig{Format: []string{"isbn"}}))
	assert.False(t, ValidateConfig("no-such-module", OptionsConfig{Type: []string{"movies"}}))
}

func TestIsValidConfigType(t *testing.T) {
	allowed := []string{"movies", "shows"}
	assert.True(t, IsValidConfigType(allowed, "shows"))
	assert.True(t, IsValidConfigType(allowed, ""), "an empty value is always accepted")
	assert.False(t, IsValidConfigType(allowed, "books"))
	assert.False(t, IsValidConfigType(nil, "movies"))

	assert.True(t, IsValidConfigTypeSlice(allowed, nil))
	assert.True(t, IsValidConfigTypeSlice(allowed, str.Slice{"shows", "movies"}))
	assert.False(t, IsValidConfigTypeSlice(allowed, str.Slice{"movies", "movies"}), "a value may not repeat more often than it is allowed")
}

func TestGetOptionTime(t *testing.T) {
	cases := []struct {
		name    string
		options str.Options
		want    string
	}{
		{name: "history", options: str.Options{Module: "history"}, want: "watched_at"},
		{name: "watchlist", options: str.Options{Module: "watchlist"}, want: "listed_at"},
		{name: "collection", options: str.Options{Module: "collection"}, want: "listed_at"},
		{name: "other module with a user", options: str.Options{Module: "lists", UserName: "sean"}, want: "listed_at"},
		{name: "other module keeps its value", options: str.Options{Module: "lists", Time: "rated_at"}, want: "rated_at"},
		{name: "other module without a user", options: str.Options{Module: "lists"}, want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			options := tc.options
			assert.Equal(t, tc.want, GetOptionTime(&options))
			assert.Equal(t, tc.want, options.Time)
		})
	}
}

func TestOptionsFromConfig(t *testing.T) {
	config := validConfig()
	config.Module = "watchlist"
	config.Type = "shows"
	config.Sort = "added"
	config.Format = "tmdb"
	config.UserName = "sean"
	config.PerPage = 25
	config.PagesLimit = 3
	config.Verbose = true

	options, err := OptionsFromConfig(credentialsFs(t), config)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "watchlist", options.Module)
	assert.Equal(t, "shows", options.Type)
	assert.Equal(t, "added", options.Sort)
	assert.Equal(t, "tmdb", options.Format)
	assert.Equal(t, "sean", options.UserName)
	assert.Equal(t, 25, options.PerPage)
	assert.Equal(t, 3, options.PagesLimit)
	assert.True(t, options.Verbose)
	assert.Equal(t, "access", options.Token.AccessToken)
	assert.Equal(t, "refresh", options.Token.RefreshToken)
	assert.Equal(t, "sean", *options.UserSettings.User.Username)
	assert.Equal(t, "Europe/Warsaw", options.Timezone, "timezone comes from the user settings")
	assert.Equal(t, "export_shows_watchlist.json", options.Output)
}

func TestOptionsFromConfigAdjustments(t *testing.T) {
	cases := []struct {
		name   string
		change func(c *Config)
		check  func(t *testing.T, o str.Options)
	}{
		{name: "unknown module becomes history",
			change: func(c *Config) { c.Module = "podcasts"; c.Type = "movies" },
			check: func(t *testing.T, o str.Options) {
				assert.Equal(t, "history", o.Module)
				assert.Equal(t, "export_movies_history.json", o.Output)
			}},
		{name: "unknown format becomes imdb",
			change: func(c *Config) { c.Module = "watchlist"; c.Format = "isbn" },
			check:  func(t *testing.T, o str.Options) { assert.Equal(t, "imdb", o.Format) }},
		{name: "unknown sort becomes rank",
			change: func(c *Config) { c.Module = "watchlist"; c.Sort = "random" },
			check:  func(t *testing.T, o str.Options) { assert.Equal(t, "rank", o.Sort) }},
		{name: "episodes use tmdb instead of imdb",
			change: func(c *Config) { c.Module = "watchlist"; c.Type = "episodes"; c.Format = "imdb" },
			check:  func(t *testing.T, o str.Options) { assert.Equal(t, "tmdb", o.Format) }},
		{name: "output from the config wins",
			change: func(c *Config) { c.Output = "mine.json" },
			check:  func(t *testing.T, o str.Options) { assert.Equal(t, "mine.json", o.Output) }},
		{name: "lists output",
			change: func(c *Config) { c.Module = "lists"; c.Type = "" },
			check:  func(t *testing.T, o str.Options) { assert.Equal(t, "export_lists_.json", o.Output) }},
		{name: "people output uses the action",
			change: func(c *Config) { c.Module = "people"; c.Type = ""; c.Action = "movies" },
			check:  func(t *testing.T, o str.Options) { assert.Equal(t, "export_people_movies.json", o.Output) }},
		{name: "default output",
			change: func(c *Config) { c.Module = "history"; c.Type = "movies" },
			check:  func(t *testing.T, o str.Options) { assert.Equal(t, "export_movies_history.json", o.Output) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := validConfig()
			tc.change(config)
			options, err := OptionsFromConfig(credentialsFs(t), config)
			if assert.NoError(t, err) {
				tc.check(t, options)
			}
		})
	}
}

func TestOptionsFromConfigErrors(t *testing.T) {
	t.Run("type not valid for the module", func(t *testing.T) {
		config := validConfig()
		config.Module = "watchlist"
		config.Type = "books"
		_, err := OptionsFromConfig(credentialsFs(t), config)
		assert.EqualError(t, err, "type 'books' is not valid for module 'watchlist'")
	})

	t.Run("no access token and no client id", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		assert.NoError(t, afero.WriteFile(fs, testTokenPath, []byte(`{}`), consts.X600))
		assert.NoError(t, afero.WriteFile(fs, testSettingsPath, []byte(testSettingsJSON), consts.X600))
		config := validConfig()
		config.ClientID = consts.EmptyString
		_, err := OptionsFromConfig(fs, config)
		assert.EqualError(t, err, "no valid Authorization header")
	})

	t.Run("client id is enough without a token", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		assert.NoError(t, afero.WriteFile(fs, testTokenPath, []byte(`{}`), consts.X600))
		assert.NoError(t, afero.WriteFile(fs, testSettingsPath, []byte(`{}`), consts.X600))
		config := validConfig()
		config.Timezone = "Asia/Tokyo"
		options, err := OptionsFromConfig(fs, config)
		if assert.NoError(t, err) {
			assert.Empty(t, options.Token.AccessToken)
			assert.Equal(t, "Asia/Tokyo", options.Timezone, "timezone falls back to the config without account settings")
		}
	})
}

// Missing or broken credential files are replaced with defaults written to the real filesystem,
// so these cases use a temp directory instead of the in-memory one.
func TestOptionsFromConfigDefaultCredentialFiles(t *testing.T) {
	t.Run("missing files", func(t *testing.T) {
		dir := t.TempDir()
		config := validConfig()
		config.TokenPath = filepath.Join(dir, "token.json")
		config.SettingsPath = filepath.Join(dir, "user_settings.json")

		options, err := OptionsFromConfig(afero.NewOsFs(), config)
		if !assert.NoError(t, err) {
			return
		}
		assert.Empty(t, options.Token.AccessToken)
		created := time.Unix(options.Token.CreatedAt, 0)
		assert.WithinDuration(t, time.Now().Add(-24*time.Hour), created, time.Minute, "the default token is a day old, so it counts as expired")
		assert.True(t, options.Token.Expired())
		assert.Equal(t, "UTC", options.Timezone)

		for _, path := range []string{config.TokenPath, config.SettingsPath} {
			info, statErr := os.Stat(path)
			if assert.NoError(t, statErr) {
				assert.Equal(t, os.FileMode(consts.X600), info.Mode().Perm(), path)
			}
		}
	})

	t.Run("files that are not json", func(t *testing.T) {
		dir := t.TempDir()
		config := validConfig()
		config.TokenPath = filepath.Join(dir, "token.json")
		config.SettingsPath = filepath.Join(dir, "user_settings.json")
		assert.NoError(t, os.WriteFile(config.TokenPath, []byte("not json"), consts.X600))
		assert.NoError(t, os.WriteFile(config.SettingsPath, []byte("not json"), consts.X600))

		options, err := OptionsFromConfig(afero.NewOsFs(), config)
		if !assert.NoError(t, err) {
			return
		}
		assert.True(t, options.Token.Expired())
		assert.Equal(t, "UTC", options.Timezone)

		data, readErr := os.ReadFile(config.SettingsPath)
		if assert.NoError(t, readErr) {
			assert.JSONEq(t, `{"account":{"timezone":"UTC"}}`, string(data))
		}
	})
}

func TestSyncOptionsFromFlags(t *testing.T) {
	flagMap := useFlags(t, "-m=watchlist", "-t=shows", "-s=added", "-f=tvdb", "-u=sean", "-o=out.json", "-v")
	file := &Config{ClientID: "client-id", ClientSecret: "client-secret", TokenPath: testTokenPath, SettingsPath: testSettingsPath, Type: "movies", PerPage: 25}

	options, err := SyncOptionsFromFlags(credentialsFs(t), file, flagMap)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "watchlist", options.Module)
	assert.Equal(t, "shows", options.Type, "the -t flag wins over the config file")
	assert.Equal(t, "added", options.Sort)
	assert.Equal(t, "tvdb", options.Format)
	assert.Equal(t, "sean", options.UserName)
	assert.Equal(t, "out.json", options.Output)
	assert.Equal(t, 25, options.PerPage)
	assert.True(t, options.Verbose)
}

func TestSyncOptionsFromFlagsConfigError(t *testing.T) {
	flagMap := useFlags(t)
	_, err := SyncOptionsFromFlags(credentialsFs(t), &Config{TokenPath: testTokenPath, SettingsPath: testSettingsPath}, flagMap)
	assert.EqualError(t, err, "error sync options from flags: config error : client_id and client_secret are required fields, update your config file")
}

func TestMergeConfigs(t *testing.T) {
	t.Run("file values replace the defaults", func(t *testing.T) {
		flagMap := useFlags(t)
		file := &Config{
			ClientID: "client-id", ClientSecret: "client-secret", TokenPath: testTokenPath, SettingsPath: testSettingsPath,
			RedirectURI: "urn:ietf:wg:oauth:2.0:oob", ErrorCode: 3, WarningCode: 4, PerPage: 25, PagesLimit: 7,
			ConfigPath: "/etc/trakt-sync.toml", Output: "file.json", Type: "shows", Format: "tmdb", UserName: "sean",
			List: "watchlist", ID: "star-wars", Module: "lists", Action: "items", Sort: "added",
		}

		got, err := MergeConfigs(DefaultConfig(), file, flagMap)
		if !assert.NoError(t, err) {
			return
		}
		want := DefaultConfig()
		want.ClientID, want.ClientSecret, want.TokenPath, want.SettingsPath = "client-id", "client-secret", testTokenPath, testSettingsPath
		want.RedirectURI, want.ErrorCode, want.WarningCode = "urn:ietf:wg:oauth:2.0:oob", 3, 4
		want.PerPage, want.PagesLimit = 25, 7
		want.ConfigPath, want.Output, want.Type, want.Format = "/etc/trakt-sync.toml", "file.json", "shows", "tmdb"
		want.UserName, want.List, want.ID, want.Module, want.Action, want.Sort = "sean", "watchlist", "star-wars", "lists", "items", "added"
		assert.Equal(t, want, got)
	})

	t.Run("flags replace file values", func(t *testing.T) {
		flagMap := useFlags(t, "-a=popular", "-c=/tmp/other.toml", "-f=tvdb", "-i=vampires", "-l=favorites", "-m=movies", "-o=flag.json", "-s=title", "-t=episodes", "-u=justin", "-v=false")
		file := &Config{
			ClientID: "client-id", ClientSecret: "client-secret", TokenPath: testTokenPath, SettingsPath: testSettingsPath, Verbose: true,
			ConfigPath: "/etc/trakt-sync.toml", Output: "file.json", Type: "shows", Format: "tmdb", UserName: "sean",
			List: "watchlist", ID: "star-wars", Module: "lists", Action: "items", Sort: "added",
		}

		got, err := MergeConfigs(DefaultConfig(), file, flagMap)
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "popular", got.Action)
		assert.Equal(t, "/tmp/other.toml", got.ConfigPath)
		assert.Equal(t, "tvdb", got.Format)
		assert.Equal(t, "vampires", got.ID)
		assert.Equal(t, "favorites", got.List)
		assert.Equal(t, "movies", got.Module)
		assert.Equal(t, "flag.json", got.Output)
		assert.Equal(t, "title", got.Sort)
		assert.Equal(t, "episodes", got.Type)
		assert.Equal(t, "justin", got.UserName)
		assert.False(t, got.Verbose)
	})

	t.Run("tilde in the credential paths is the home directory", func(t *testing.T) {
		usr, err := user.Current()
		if err != nil {
			t.Skipf("no current user: %v", err)
		}
		flagMap := useFlags(t)
		file := &Config{ClientID: "client-id", ClientSecret: "client-secret", TokenPath: "~/token.json", SettingsPath: "~/user_settings.json"}

		got, err := MergeConfigs(DefaultConfig(), file, flagMap)
		if assert.NoError(t, err) {
			assert.Equal(t, filepath.Join(usr.HomeDir, "token.json"), got.TokenPath)
			assert.Equal(t, filepath.Join(usr.HomeDir, "user_settings.json"), got.SettingsPath)
		}
	})

	errorCases := []struct {
		name string
		file Config
		want string
	}{
		{name: "no client secret", file: Config{ClientID: "client-id", TokenPath: testTokenPath, SettingsPath: testSettingsPath},
			want: "config error : client_id and client_secret are required fields, update your config file"},
		{name: "token path is not json", file: Config{ClientID: "client-id", ClientSecret: "client-secret", TokenPath: "/home/tester/token.txt", SettingsPath: testSettingsPath},
			want: "config error : token_path should be json file, update your config file"},
		{name: "no settings path", file: Config{ClientID: "client-id", ClientSecret: "client-secret", TokenPath: testTokenPath},
			want: "config error : settings_path should be json file, update your config file"},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			flagMap := useFlags(t)
			file := tc.file
			got, err := MergeConfigs(DefaultConfig(), &file, flagMap)
			assert.Nil(t, got)
			assert.EqualError(t, err, tc.want)
		})
	}
}

func TestGetConfig(t *testing.T) {
	const path = "/home/tester/trakt-sync.toml"
	write := func(t *testing.T, content string) afero.Fs {
		t.Helper()
		fs := afero.NewMemMapFs()
		assert.NoError(t, afero.WriteFile(fs, path, []byte(content), consts.X600))
		return fs
	}

	t.Run("config file", func(t *testing.T) {
		fs := write(t, "client_id = \"client-id\"\nclient_secret = \"client-secret\"\ntoken_path = \"/home/tester/token.json\"\nsettings_path = \"/home/tester/user_settings.json\"\nper_page = 25\n")
		config, err := GetConfig(fs, path)
		if assert.NoError(t, err) {
			assert.Equal(t, "client-id", config.ClientID)
			assert.Equal(t, testTokenPath, config.TokenPath)
			assert.Equal(t, 25, config.PerPage)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		config, err := GetConfig(afero.NewMemMapFs(), path)
		assert.Nil(t, config)
		assert.EqualError(t, err, "cannot read the config file")
	})

	t.Run("file that is not toml", func(t *testing.T) {
		config, err := GetConfig(write(t, "client_id = "), path)
		assert.Nil(t, config)
		assert.ErrorContains(t, err, "cannot parse the config file")
	})

	t.Run("file without credentials", func(t *testing.T) {
		config, err := GetConfig(write(t, "per_page = 25\n"), path)
		assert.Nil(t, config)
		assert.EqualError(t, err, "config normalize error : client_id and client_secret are required fields, update your config file")
	})

	t.Run("no path uses the defaults, which have no credentials", func(t *testing.T) {
		config, err := GetConfig(afero.NewMemMapFs(), consts.EmptyString)
		assert.Nil(t, config)
		assert.ErrorContains(t, err, "client_id and client_secret are required fields")
	})
}

func TestGetOutputForModuleByAction(t *testing.T) {
	calendar := str.Options{Module: "calendars", StartDate: "2026-10-01", Days: 7}
	comment := str.Options{Module: "comments", CommentID: 417, Type: "movies"}
	sync := str.Options{Module: "sync", Type: "movies"}
	cases := []struct {
		options str.Options
		action  string
		want    string
	}{
		{options: calendar, action: "my_shows", want: "export_calendars_shows_20261001_7.json"},
		{options: calendar, action: "all_shows", want: "export_calendars_shows_20261001_7.json"},
		{options: calendar, action: "my_new_shows", want: "export_calendars_new_shows_20261001_7.json"},
		{options: calendar, action: "all_new_shows", want: "export_calendars_new_shows_20261001_7.json"},
		{options: calendar, action: "my_season_premieres", want: "export_calendars_season_premieres_20261001_7.json"},
		{options: calendar, action: "all_season_premieres", want: "export_calendars_season_premieres_20261001_7.json"},
		{options: calendar, action: "my_finales", want: "export_calendars_finales_20261001_7.json"},
		{options: calendar, action: "all_finales", want: "export_calendars_finales_20261001_7.json"},
		{options: calendar, action: "my_movies", want: "export_calendars_movies_20261001_7.json"},
		{options: calendar, action: "all_movies", want: "export_calendars_movies_20261001_7.json"},
		{options: calendar, action: "my_dvd", want: "export_calendars_dvd_20261001_7.json"},
		{options: calendar, action: "all_dvd", want: "export_calendars_dvd_20261001_7.json"},
		{options: calendar, action: "my_media", want: "export_calendars_media_20261001_7.json"},
		{options: calendar, action: "all_media", want: "export_calendars_media_20261001_7.json"},
		{options: calendar, action: "my_streaming", want: "export_calendars_streaming_20261001_7.json"},
		{options: calendar, action: "all_streaming", want: "export_calendars_streaming_20261001_7.json"},
		{options: calendar, action: "hot_releases", want: "export_calendars_hot_releases_20261001_7.json"},
		{options: calendar, action: "hot_premieres", want: "export_calendars_hot_premieres_20261001_7.json"},
		{options: calendar, action: "hot_new_shows", want: "export_calendars_hot_new_shows_20261001_7.json"},
		{options: calendar, action: "hot_finales", want: "export_calendars_hot_finales_20261001_7.json"},
		{options: calendar, action: "unknown", want: "export_calendars.json"},
		{options: comment, action: "comment", want: "export_comments_comment_417.json"},
		{options: comment, action: "replies", want: "export_comments_replies_417.json"},
		{options: comment, action: "item", want: "export_comments_item_417.json"},
		{options: comment, action: "likes", want: "export_comments_likes_417.json"},
		{options: comment, action: "reactions", want: "export_comments_reactions_417.json"},
		{options: comment, action: "reactions_summary", want: "export_comments_reactions_summary_417.json"},
		{options: comment, action: "trending", want: "export_comments_trending.json"},
		{options: comment, action: "recent", want: "export_comments_recent.json"},
		{options: comment, action: "updates", want: "export_comments_updates.json"},
		{options: comment, action: "unknown", want: "export_comments_movies.json"},
		{options: sync, action: "update_watchlist", want: "export_sync_update_watchlist.json"},
		{options: sync, action: "update_favorites", want: "export_sync_update_favorites.json"},
		{options: sync, action: "get_watchlist", want: "export_sync_watchlist_movies.json"},
		{options: sync, action: "get_favorites", want: "export_sync_favorites_movies.json"},
		{options: sync, action: "get_ratings", want: "export_sync_ratings_movies.json"},
		{options: sync, action: "get_history", want: "export_sync_history_movies.json"},
		{options: sync, action: "get_watched", want: "export_sync_watched_movies.json"},
		{options: sync, action: "get_collection", want: "export_sync_collection_movies.json"},
		{options: sync, action: "last_activities", want: "export_sync_last_activities.json"},
		{options: sync, action: "unknown", want: "export_sync_movies.json"},
		{options: str.Options{Module: "certifications", Type: "shows"}, want: "export_certifications_shows.json"},
		{options: str.Options{Module: "certifications", Type: "books"}, want: "export_certifications.json"},
		{options: str.Options{Module: "no-such-module", Type: "movies", Format: "imdb"}, want: "export_no-such-module_movies_imdb.json"},
	}
	for _, tc := range cases {
		t.Run(tc.options.Module+" "+tc.action, func(t *testing.T) {
			options := tc.options
			options.Action = tc.action
			assert.Equal(t, tc.want, GetOutputForModule(&options))
		})
	}
}
