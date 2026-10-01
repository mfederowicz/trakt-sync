// Package handlers used to handle module actions
package handlers

import (
	"io"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// postHandler is a handler that builds one object from flags, looking the item up first when needed, and posts it.
type postHandler struct {
	name     string
	handler  Handler
	options  str.Options
	requests []string
	sends    string
}

func postHandlers() []postHandler {
	const (
		episodeID   = `"episode":{"ids":{"trakt":55}}`
		fetched     = `{"title":"Tron","ids":{"trakt":1}}`
		getMovie    = "GET /movies/55"
		getPerson   = "GET /people/55"
		getSettings = "GET /users/settings"
		getShow     = "GET /shows/55"
		postCheckin = "POST /checkin"
		postNotes   = "POST /notes"
		postPause   = "POST /scrobble/pause"
		postStart   = "POST /scrobble/start"
		postStop    = "POST /scrobble/stop"
	)
	note := func(itemType string) str.Options {
		return str.Options{InternalID: "55", Notes: "rewatch", Type: itemType}
	}
	scrobble := func(itemType string) str.Options {
		return str.Options{InternalID: "55", Progress: 42.5, Type: itemType}
	}
	byCode := scrobble("show_episode")
	byCode.EpisodeCode = "2x05"
	byAbs := scrobble("show_episode")
	byAbs.EpisodeAbs = 15
	checkinCode := str.Options{Action: "show_episode", TraktID: 55, InternalID: "55", EpisodeCode: "2x05"}
	checkinAbs := str.Options{Action: "show_episode", TraktID: 55, InternalID: "55", EpisodeAbs: 15}
	return []postHandler{
		{name: "checkin episode", handler: CheckinEpisodeHandler{}, options: str.Options{Action: "episode", TraktID: 55}, requests: []string{getSettings, postCheckin}, sends: episodeID},
		{name: "checkin show episode abs", handler: CheckinShowEpisodeHandler{}, options: checkinAbs, requests: []string{getSettings, getShow, postCheckin}, sends: `"episode":{"number_abs":15}`},
		{name: "checkin show episode code", handler: CheckinShowEpisodeHandler{}, options: checkinCode, requests: []string{getSettings, getShow, postCheckin}, sends: `"episode":{"season":2,"number":5}`},
		{name: "notes episode", handler: NotesNotesEpisodeHandler{}, options: note("episode"), requests: []string{postNotes}, sends: episodeID},
		{name: "notes history", handler: NotesNotesHistoryHandler{}, options: note("history"), requests: []string{postNotes}, sends: `"attached_to":{"id":55,"type":"history"}`},
		{name: "notes movie", handler: NotesNotesMovieHandler{}, options: note("movie"), requests: []string{getMovie, postNotes}, sends: `"movie":` + fetched},
		{name: "notes person", handler: NotesNotesPersonHandler{}, options: note("person"), requests: []string{getPerson, postNotes}, sends: `"person":`},
		{name: "notes season", handler: NotesNotesSeasonHandler{}, options: note("season"), requests: []string{postNotes}, sends: `"season":{"ids":{"trakt":55}}`},
		{name: "notes show", handler: NotesNotesShowHandler{}, options: note("show"), requests: []string{getShow, postNotes}, sends: `"show":` + fetched},
		{name: "scrobble pause episode", handler: ScrobblePauseEpisodeHandler{}, options: scrobble("episode"), requests: []string{postPause}, sends: episodeID},
		{name: "scrobble pause movie", handler: ScrobblePauseMovieHandler{}, options: scrobble("movie"), requests: []string{getMovie, postPause}, sends: `"movie":` + fetched},
		{name: "scrobble pause show episode abs", handler: ScrobblePauseShowEpisodeHandler{}, options: byAbs, requests: []string{getShow, postPause}, sends: `"number_abs":15`},
		{name: "scrobble pause show episode code", handler: ScrobblePauseShowEpisodeHandler{}, options: byCode, requests: []string{getShow, postPause}, sends: `"season":2,"number":5`},
		{name: "scrobble start episode", handler: ScrobbleStartEpisodeHandler{}, options: scrobble("episode"), requests: []string{postStart}, sends: episodeID},
		{name: "scrobble start movie", handler: ScrobbleStartMovieHandler{}, options: scrobble("movie"), requests: []string{getMovie, postStart}, sends: `"progress":42.5`},
		{name: "scrobble start show episode abs", handler: ScrobbleStartShowEpisodeHandler{}, options: byAbs, requests: []string{getShow, postStart}, sends: `"number_abs":15`},
		{name: "scrobble start show episode code", handler: ScrobbleStartShowEpisodeHandler{}, options: byCode, requests: []string{getShow, postStart}, sends: `"season":2,"number":5`},
		{name: "scrobble stop episode", handler: ScrobbleStopEpisodeHandler{}, options: scrobble("episode"), requests: []string{postStop}, sends: episodeID},
		{name: "scrobble stop movie", handler: ScrobbleStopMovieHandler{}, options: scrobble("movie"), requests: []string{getMovie, postStop}, sends: `"movie":` + fetched},
		{name: "scrobble stop show episode abs", handler: ScrobbleStopShowEpisodeHandler{}, options: byAbs, requests: []string{getShow, postStop}, sends: `"number_abs":15`},
		{name: "scrobble stop show episode code", handler: ScrobbleStopShowEpisodeHandler{}, options: byCode, requests: []string{getShow, postStop}, sends: `"season":2,"number":5`},
	}
}

// servePostHandler answers the lookups with a small item and the POST with postStatus. It returns the requests
// in the order they arrived and the body of the POST.
func servePostHandler(t *testing.T, mux *http.ServeMux, postStatus int) (requests *[]string, posted *string) {
	t.Helper()
	requests, posted = &[]string{}, new(string)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		*requests = append(*requests, r.Method+" "+r.URL.Path)
		switch {
		case r.URL.Path == "/users/settings":
			test.SafeFprint(w, `{"connections":{"twitter":true}}`)
		case r.Method == http.MethodGet:
			test.SafeFprint(w, `{"title":"Tron","ids":{"trakt":1}}`)
		default:
			body, err := io.ReadAll(r.Body)
			test.AssertNilError(t, err)
			*posted = string(body)
			w.WriteHeader(postStatus)
			test.SafeFprint(w, `{"id":9}`)
		}
	})
	return requests, posted
}

func TestPostHandlersSendItem(t *testing.T) {
	for _, tc := range postHandlers() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			requests, posted := servePostHandler(t, s.Mux, http.StatusCreated)

			options := tc.options
			test.AssertNilError(t, tc.handler.Handle(&options, s.Client))
			assert.Equal(t, tc.requests, *requests)
			assert.Contains(t, *posted, tc.sends)
		})
	}
}

func TestPostHandlersFailedPost(t *testing.T) {
	for _, tc := range postHandlers() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			requests, _ := servePostHandler(t, s.Mux, http.StatusInternalServerError)

			options := tc.options
			assert.Error(t, tc.handler.Handle(&options, s.Client))
			assert.Equal(t, tc.requests, *requests)
		})
	}
}
