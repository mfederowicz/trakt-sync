// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestSyncGetCollectionHandlerTypes(t *testing.T) {
	tests := []struct {
		collectionType string
		wantErr        bool
	}{
		{collectionType: cfg.DefaultConfig().Type}, // no -t given: default type from cfg.DefaultConfig()
		{collectionType: "shows"},
		{collectionType: "episodes"},
		{collectionType: "media"},
		{collectionType: "movie", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.collectionType, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			called := false
			s.Mux.HandleFunc("/sync/collection/"+tt.collectionType, func(w http.ResponseWriter, r *http.Request) {
				called = true
				test.AssertMethod(t, r, http.MethodGet)
				test.SafeFprint(w, `[{"type":"movie","movie":{"title":"TRON: Legacy"}}]`)
			})

			options := &str.Options{Module: "sync", Action: consts.GetCollection, Type: tt.collectionType, Output: filepath.Join(t.TempDir(), "out.json")}
			err := SyncGetCollectionHandler{}.Handle(options, s.Client)
			if tt.wantErr {
				if err == nil || !strings.Contains(err.Error(), "not found type") {
					t.Fatalf("error is %v, want a type error", err)
				}
				if called {
					t.Error("API was called for an invalid type")
				}
				return
			}
			test.AssertNilError(t, err)
			if !called {
				t.Errorf("/sync/collection/%s was not called", tt.collectionType)
			}
		})
	}
}
