// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	Client    *internal.Client
	Mux       *http.ServeMux
	ServerURL string
	Teardown  func()
}

// setup sets up a test HTTP server along with a trakt.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup(t *testing.T) *TestSetup {
	t.Helper()
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(consts.BaseURLPath+"/", http.StripPrefix(consts.BaseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the Trakt client being tested and is
	// configured to use test server.
	client := internal.NewClient(nil)
	uri, _ := url.Parse(server.URL + consts.BaseURLPath + "/")
	client.BaseURL = uri

	return &TestSetup{
		Client:    client,
		Mux:       mux,
		ServerURL: server.URL,
		Teardown:  server.Close,
	}
}

func MuxUserSettings(t *testing.T, mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/users/settings", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, "GET")
		s := str.UserSettings{}
		val := true
		connections := str.Connections{}
		connections.Facebook = &val
		s.Connections = &connections
		user, _ := json.Marshal(s)
		test.SafeFprint(w, string(user))
	})
	return mux
}

func MuxShow(t *testing.T, mux *http.ServeMux, o *str.Options) *http.ServeMux {
	mux.HandleFunc("/shows/"+o.InternalID, func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, "GET")
		test.SafeFprint(w,
			`{
			  "title": "Test show",
			  "year": 2011,
			  "ids": {
				"trakt": 353,
				"slug": "game-of-thrones",
				"tvdb": 121361,
				"imdb": "tt0944947",
				"tmdb": 1399
			  }
			}`,
		)
	})
	return mux
}

func TestEmptyServeMux(t *testing.T) {
	// Verify that a ServeMux with nothing registered
	// doesn't panic.
	testSetup := setup(t)
	mux := testSetup.Mux
	var r http.Request
	r.Method = "GET"
	r.Host = "example.com"
	r.URL = &url.URL{Path: "/"}
	_, p := mux.Handler(&r)
	if p != "" {
		t.Errorf(`got %q, want ""`, p)
	}
}

func TestCreateCheckinUserSettingsError(t *testing.T) {
	testSetup := setup(t)
	c := &CommonLogic{}
	o := &str.Options{}
	s := &str.UserSettings{}
	a := &str.UserAccount{}
	tz := "Europe/Warsaw"
	a.Timezone = &tz
	s.Account = a
	o.UserSettings = *s
	_, err := c.CreateCheckin(testSetup.Client, o)
	assert.Contains(t, err.Error(), "user settings error")
}

func TestCreateCheckinUnknownAction(t *testing.T) {
	testSetup := setup(t)
	mux := testSetup.Mux
	mux = MuxUserSettings(t, mux)
	c := &CommonLogic{}
	o := &str.Options{}
	s := &str.UserSettings{}
	a := &str.UserAccount{}
	tz := "Europe/Warsaw"
	a.Timezone = &tz
	s.Account = a
	o.UserSettings = *s
	_, err := c.CreateCheckin(testSetup.Client, o)
	assert.Equal(t, err.Error(), "uknown checkin action")
}

func TestCreateCheckinForMovie(t *testing.T) {
	testSetup := setup(t)
	mux := testSetup.Mux
	mux = MuxUserSettings(t, mux)
	c := &CommonLogic{}
	o := &str.Options{}
	s := &str.UserSettings{}
	a := &str.UserAccount{}
	tz := "Europe/Warsaw"
	a.Timezone = &tz
	s.Account = a
	o.UserSettings = *s
	o.Action = consts.Movie
	o.InternalID = "despicable-me-4-2024"
	mux.HandleFunc("/movies/despicable-me-4-2024", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, "GET")
		test.SafeFprint(w,
			`{
			  "title": "Despicable Me 4x",
			  "year": 2024,
			  "ids": {
				"trakt": 367444,
				"slug": "despicable-me-4-2024",
				"imdb": "tt7510222",
				"tmdb": 519182
			  }
			}`,
		)
	})
	checkin, _ := c.CreateCheckin(testSetup.Client, o)
	test.AssertType(t, checkin, consts.Fupper(consts.Checkin))
	assert.Equal(t, checkin.Movie.IDs.Trakt, test.Ptr(int64(consts.TestMovieTraktID)))
}

func TestCreateCheckinForEpisode(t *testing.T) {
	testSetup := setup(t)
	mux := testSetup.Mux
	mux = MuxUserSettings(t, mux)
	c := &CommonLogic{}
	o := &str.Options{}
	s := &str.UserSettings{}
	a := &str.UserAccount{}
	tz := "Europe/Warsaw"
	a.Timezone = &tz
	s.Account = a
	o.UserSettings = *s
	o.Action = consts.Episode
	o.InternalID = "12345"
	mux.HandleFunc("/episodes/12345", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, "GET")
		test.SafeFprint(w,
			`{
				  "season": 6,
				  "number": 21,
				  "title": "Made in America",
				  "ids": {
					"trakt": 73629,
					"tvdb": 329768,
					"imdb": "tt0995839",
					"tmdb": 63055,
					"tvrage": null
				  }
				}`,
		)
	})

	checkin, _ := c.CreateCheckin(testSetup.Client, o)
	test.AssertType(t, checkin, consts.Fupper(consts.Checkin))
	assert.Equal(t, checkin.Episode.IDs.Trakt, test.Ptr(int64(consts.TestEpisodeTraktID)))
}

func TestCreateCheckinForShowEpisodeInvalidLength(t *testing.T) {
	testSetup := setup(t)
	mux := testSetup.Mux
	mux = MuxUserSettings(t, mux)
	c := &CommonLogic{}
	o := &str.Options{}
	s := &str.UserSettings{}
	a := &str.UserAccount{}
	tz := "Europe/Warsaw"
	a.Timezone = &tz
	s.Account = a
	o.UserSettings = *s
	o.Action = consts.ShowEpisode
	o.InternalID = "12345"
	o.EpisodeCode = "12"
	mux = MuxShow(t, mux, o)
	_, err := c.CreateCheckin(testSetup.Client, o)
	assert.Contains(t, err.Error(), "invalid length")
}

func TestCreateCheckinForShowEpisodeInvalidFormat(t *testing.T) {
	testSetup := setup(t)
	mux := testSetup.Mux
	mux = MuxUserSettings(t, mux)
	c := &CommonLogic{}
	o := &str.Options{}
	s := &str.UserSettings{}
	a := &str.UserAccount{}
	tz := "Europe/Warsaw"
	a.Timezone = &tz
	s.Account = a
	o.UserSettings = *s
	o.Action = consts.ShowEpisode
	o.EpisodeCode = "123456"
	mux = MuxShow(t, mux, o)
	_, err := c.CreateCheckin(testSetup.Client, o)
	assert.Contains(t, err.Error(), "invalid format")
}

func TestCreateCheckinForShowEpisodeEpisodeCode(t *testing.T) {
	testSetup := setup(t)
	mux := testSetup.Mux
	mux = MuxUserSettings(t, mux)
	c := &CommonLogic{}
	o := &str.Options{}
	s := &str.UserSettings{}
	a := &str.UserAccount{}
	tz := "Europe/Warsaw"
	a.Timezone = &tz
	s.Account = a
	o.UserSettings = *s
	o.Action = consts.ShowEpisode
	o.EpisodeCode = "6x10"
	o.InternalID = "353"
	mux = MuxShow(t, mux, o)
	checkin, _ := c.CreateCheckin(testSetup.Client, o)
	assert.Equal(t, checkin.Episode.Season, test.Ptr(consts.TestEpisodeSeason6))
	assert.Equal(t, checkin.Episode.Number, test.Ptr(consts.TestEpisodeNumber10))
	test.AssertType(t, checkin, consts.Fupper(consts.Checkin))
}

func TestCreateCheckinForShowEpisodeAbs(t *testing.T) {
	testSetup := setup(t)
	mux := testSetup.Mux
	mux = MuxUserSettings(t, mux)
	c := &CommonLogic{}
	o := &str.Options{}
	s := &str.UserSettings{}
	a := &str.UserAccount{}
	tz := "Europe/Warsaw"
	a.Timezone = &tz
	s.Account = a
	o.UserSettings = *s
	o.Action = consts.ShowEpisode
	o.EpisodeAbs = consts.TestEpisodeAbs
	o.InternalID = "353"
	mux = MuxShow(t, mux, o)

	checkin, _ := c.CreateCheckin(testSetup.Client, o)
	assert.Equal(t, checkin.Episode.NumberAbs, test.Ptr(consts.TestEpisodeAbs))
	test.AssertType(t, checkin, consts.Fupper(consts.Checkin))
}

func TestConvertDateString(t *testing.T) {
	testSetup := setup(t)
	mux := testSetup.Mux
	mux = MuxUserSettings(t, mux)
	c := &CommonLogic{}
	o := &str.Options{}
	o.ResetAt = "2025-01-24"
	out := c.ConvertDateString(o.ResetAt, consts.DefaultStartDateFormat, "Europe/Warsaw", true)
	assert.Contains(t, out, "+01:00")
	o.ResetAt = "2025-05-24"
	out = c.ConvertDateString(o.ResetAt, consts.DefaultStartDateFormat, "Europe/Warsaw", true)
	assert.Contains(t, out, "+02:00")
}

func TestCurrnetDateString(t *testing.T) {
	testSetup := setup(t)
	mux := testSetup.Mux
	mux = MuxUserSettings(t, mux)
	c := &CommonLogic{}
	out := c.CurrentDateString(time.UTC.String(), true)
	currentTime := time.Now().UTC().Truncate(time.Hour)
	assert.Contains(t, out, currentTime.Format(time.RFC3339))
}
