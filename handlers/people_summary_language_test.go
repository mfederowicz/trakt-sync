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

// people -a summary sends language only when -language is set and explains a 400.
func TestPeopleSummaryLanguage(t *testing.T) {
	cases := []struct {
		name     string
		language string
		status   int
		query    string
		wantErr  string
	}{
		{name: "language", language: "pt-br", status: http.StatusOK, query: "extended=full&language=pt-br"},
		{name: "no language", status: http.StatusOK, query: "extended=full"},
		{name: "malformed language", language: "x_y", status: http.StatusBadRequest, query: "extended=full&language=x_y",
			wantErr: `language "x_y" rejected by Trakt, use a code like pt or pt-br`},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			queries := []string{}
			s.Mux.HandleFunc("/people/bryan-cranston", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				queries = append(queries, r.URL.RawQuery)
				w.WriteHeader(tc.status)
				if tc.status == http.StatusOK {
					test.SafeFprint(w, `{"name":"Bryan Cranston","social_ids":{"wikipedia":"Bryan_Cranston","wikipedia_language":"pt"}}`)
				}
			})

			options := str.Options{ID: "bryan-cranston", ExtendedInfo: "full", Language: tc.language, Output: filepath.Join(t.TempDir(), "out.json")}
			err := PeopleSummaryHandler{}.Handle(&options, s.Client)
			assert.Equal(t, []string{tc.query}, queries)
			if tc.wantErr == "" {
				test.AssertNilError(t, err)
				return
			}
			assert.ErrorContains(t, err, tc.wantErr)
		})
	}
}
