// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// movies and shows -a related send personalized=true only when -personalized is set.
func TestRelatedPersonalized(t *testing.T) {
	handlers := []struct {
		name    string
		handler Handler
		path    string
		body    string
	}{
		{name: "movies related", handler: MoviesRelatedHandler{}, path: "/movies/tron-legacy-2010/related", body: `[{"title":"TRON","ids":{"trakt":1}}]`},
		{name: "shows related", handler: ShowsRelatedHandler{}, path: "/shows/tron-legacy-2010/related", body: `[{"title":"Andor","ids":{"trakt":2}}]`},
	}
	cases := []struct {
		name         string
		personalized bool
		query        string
	}{
		{name: "personalized", personalized: true, query: "page=1&personalized=true"},
		{name: "not personalized", query: "page=1"},
	}
	for _, h := range handlers {
		for _, tc := range cases {
			h, tc := h, tc
			t.Run(h.name+" "+tc.name, func(t *testing.T) {
				s := setup(t)
				defer s.Teardown()
				queries := []string{}
				s.Mux.HandleFunc(h.path, func(w http.ResponseWriter, r *http.Request) {
					test.AssertMethod(t, r, http.MethodGet)
					queries = append(queries, r.URL.RawQuery)
					test.SafeFprint(w, h.body)
				})

				options := str.Options{InternalID: "tron-legacy-2010", Personalized: tc.personalized, Output: filepath.Join(t.TempDir(), "out.json")}
				test.AssertNilError(t, h.handler.Handle(&options, s.Client))
				assert.Equal(t, []string{tc.query}, queries)
			})
		}
	}
}
