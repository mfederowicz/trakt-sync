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

// syncListSettingsHandlers are the sync actions that update the settings of the favorites and the watchlist.
var syncListSettingsHandlers = []struct {
	action  string
	handler Handler
	path    string
}{
	{action: consts.UpdateFavorites, handler: SyncUpdateFavoritesHandler{}, path: "/sync/favorites"},
	{action: consts.UpdateWatchlist, handler: SyncUpdateWatchlistHandler{}, path: "/sync/watchlist"},
}

// sync -a update_favorites and update_watchlist send only the settings that are given, so the others keep their value.
func TestSyncUpdateListSettingsSendsGivenSettings(t *testing.T) {
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
	for _, h := range syncListSettingsHandlers {
		h := h
		for _, tt := range tests {
			tt := tt
			t.Run(h.action+" "+tt.name, func(t *testing.T) {
				s := setup(t)
				defer s.Teardown()
				bodies := []string{}
				s.Mux.HandleFunc(h.path, func(w http.ResponseWriter, r *http.Request) {
					test.AssertMethod(t, r, http.MethodPut)
					body, err := io.ReadAll(r.Body)
					test.AssertNilError(t, err)
					bodies = append(bodies, string(body))
					test.SafeFprint(w, `{}`)
				})

				options := str.Options{Module: "sync", Action: h.action, Type: "movies", Description: tt.description,
					SortBy: tt.sortBy, SortHow: tt.sortHow, Output: filepath.Join(t.TempDir(), "out.json")}
				test.AssertNilError(t, h.handler.Handle(&options, s.Client))
				assert.Len(t, bodies, 1)
				assert.JSONEq(t, tt.sends, bodies[0])
			})
		}
	}
}

func TestSyncUpdateListSettingsWithoutSettings(t *testing.T) {
	for _, h := range syncListSettingsHandlers {
		h := h
		t.Run(h.action, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()
			requests := 0
			s.Mux.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request) {
				requests++
			})

			options := str.Options{Module: "sync", Action: h.action, Type: "movies", Output: filepath.Join(t.TempDir(), "out.json")}
			err := h.handler.Handle(&options, s.Client)
			assert.ErrorContains(t, err, consts.EmptyListSettingsMsg)
			assert.Zero(t, requests)
		})
	}
}
