// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// movies, shows and media -a anticipated_recommendations send page and limit only and write the items to -o.
func TestAnticipatedRecommendationsHandlers(t *testing.T) {
	handlers := []struct {
		name    string
		handler Handler
		path    string
		body    string
	}{
		{name: "movies", handler: MoviesAnticipatedRecommendationsHandler{}, path: "/movies/recommendations/anticipated",
			body: `[{"score":0.92,"released":"2026-12-18","movie":{"title":"Avatar: Fire and Ash","ids":{"trakt":1}}}]`},
		{name: "shows", handler: ShowsAnticipatedRecommendationsHandler{}, path: "/shows/recommendations/anticipated",
			body: `[{"score":0.81,"released":null,"show":{"title":"Andor","ids":{"trakt":2}}}]`},
		{name: "media", handler: MediaAnticipatedRecommendationsHandler{}, path: "/media/recommendations/anticipated",
			body: `[{"score":0.92,"released":"2026-12-18","movie":{"title":"Avatar: Fire and Ash","ids":{"trakt":1}}},{"score":0.81,"released":null,"show":{"title":"Andor","ids":{"trakt":2}}}]`},
	}
	for _, h := range handlers {
		h := h
		t.Run(h.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			queries := []string{}
			s.Mux.HandleFunc(h.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				queries = append(queries, r.URL.RawQuery)
				test.SafeFprint(w, h.body)
			})

			out := filepath.Join(t.TempDir(), "out.json")
			options := str.Options{PerPage: 100, ExtendedInfo: "full", Output: out}
			test.AssertNilError(t, h.handler.Handle(&options, s.Client))
			assert.Equal(t, []string{"limit=100&page=1"}, queries)
			data, err := os.ReadFile(out)
			test.AssertNilError(t, err)
			// a null released date is left out of the export (omitempty), like other optional fields
			assert.JSONEq(t, strings.ReplaceAll(h.body, `"released":null,`, ""), string(data))
		})
	}
}

func TestAnticipatedRecommendationsHandlerEmpty(t *testing.T) {
	s := setup(t)
	defer s.Teardown()
	s.Mux.HandleFunc("/movies/recommendations/anticipated", func(w http.ResponseWriter, _ *http.Request) {
		test.SafeFprint(w, `[]`)
	})

	options := str.Options{Output: filepath.Join(t.TempDir(), "out.json")}
	err := MoviesAnticipatedRecommendationsHandler{}.Handle(&options, s.Client)
	assert.ErrorContains(t, err, consts.EmptyResult)
}
