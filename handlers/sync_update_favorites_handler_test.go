// Package handlers used to handle module actions
package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// sync -a update_favorites sends only the settings that are given, so the others keep their value.
func TestSyncUpdateFavoritesSendsGivenSettings(t *testing.T) {
	tests := []struct {
		name        string
		description string
		sortBy      string
		sortHow     string
		sends       string
	}{
		{name: "description only", description: "best of", sends: `{"description":"best of"}`},
		{name: "sort_by only", sortBy: "added", sends: `{"sort_by":"added"}`},
		{name: "sort_how only", sortHow: "desc", sends: `{"sort_how":"desc"}`},
		{name: "all settings", description: "best of", sortBy: "added", sortHow: "desc", sends: `{"description":"best of","sort_by":"added","sort_how":"desc"}`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			bodies := []string{}
			s.Mux.HandleFunc("/sync/favorites", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodPut)
				body, err := io.ReadAll(r.Body)
				test.AssertNilError(t, err)
				bodies = append(bodies, string(body))
				test.SafeFprint(w, `{}`)
			})

			options := str.Options{Module: "sync", Action: consts.UpdateFavorites, Type: "movies", Description: tt.description,
				SortBy: tt.sortBy, SortHow: tt.sortHow, Output: filepath.Join(t.TempDir(), "out.json")}
			test.AssertNilError(t, SyncUpdateFavoritesHandler{}.Handle(&options, s.Client))
			assert.Len(t, bodies, 1)
			assert.JSONEq(t, tt.sends, bodies[0])
		})
	}
}

func TestSyncUpdateFavoritesWithoutSettings(t *testing.T) {
	s := setup(t)
	defer s.Teardown()
	requests := 0
	s.Mux.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request) {
		requests++
	})

	options := str.Options{Module: "sync", Action: consts.UpdateFavorites, Type: "movies", Output: filepath.Join(t.TempDir(), "out.json")}
	err := SyncUpdateFavoritesHandler{}.Handle(&options, s.Client)
	assert.ErrorContains(t, err, consts.EmptyFavoritesSettingsMsg)
	assert.Zero(t, requests)
}
