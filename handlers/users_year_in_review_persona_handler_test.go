// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersYearInReviewPersonaHandler(t *testing.T) {
	const persona = `{"persona":"cinephile","runner_up":null,"confidence":"strong","rarity":12,"card_number":7,"traits":["night-owl"],` +
		`"highlights":[{"kind":"movies","value":212}],"runner_up_highlights":[],"scores":{"cinephile":91},"streak":{"longest":21,"started_at":null},"monthly":[]}`
	tests := []struct {
		name      string
		year      int
		status    int
		body      string
		wantErr   string
		wantNoAPI bool
	}{
		{name: "persona", year: 2026, status: http.StatusOK, body: persona},
		{name: "without year", wantErr: "set -year 2026 or later", wantNoAPI: true},
		{name: "year before personas", year: 2025, wantErr: "set -year 2026 or later", wantNoAPI: true},
		{name: "not found", year: 2026, status: http.StatusNotFound, body: `{"error":"persona_not_found"}`, wantErr: "no year_in_review_persona for 2026 (user sean)"},
		{name: "empty object", year: 2026, status: http.StatusOK, body: `{}`, wantErr: "no year_in_review_persona for 2026 (user sean)"},
		{name: "vip required", year: 2026, status: http.StatusForbidden, body: `{"error":"vip_required"}`, wantErr: "year_in_review_persona is VIP only"},
		{name: "not open to api apps", year: 2026, status: http.StatusUnauthorized, wantErr: "not open to API apps"},
		{name: "invalid year on the server", year: 2031, status: http.StatusBadRequest, body: `{"error":"invalid_year"}`, wantErr: "fetch year_in_review_persona error"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				if want := "/users/sean/yir/" + strconv.Itoa(tt.year) + "/persona"; r.URL.Path != want {
					t.Errorf("path is %q, want %q", r.URL.Path, want)
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			options := str.Options{Action: consts.YearInReviewPersona, UserName: "sean", Year: tt.year, Output: filepath.Join(t.TempDir(), "out.json")}
			err := UsersYearInReviewPersonaHandler{}.Handle(&options, s.Client)
			if tt.wantNoAPI && calls != 0 {
				t.Error("API was called with invalid options")
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				if _, statErr := os.Stat(options.Output); !os.IsNotExist(statErr) {
					t.Errorf("output file written on error: %v", statErr)
				}
				return
			}
			test.AssertNilError(t, err)
			data, readErr := os.ReadFile(options.Output)
			test.AssertNilError(t, readErr)
			if !strings.Contains(string(data), `"persona": "cinephile"`) || !strings.Contains(string(data), `"cinephile": 91`) {
				t.Errorf("export misses the persona or its scores: %s", data)
			}
		})
	}
}
