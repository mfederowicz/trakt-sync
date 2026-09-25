// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestTeamMembersHandler(t *testing.T) {
	s := setup(t)
	defer s.Teardown()

	s.Mux.HandleFunc("/team", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		if got := r.URL.Query().Get("extended"); got != "images" {
			t.Errorf("extended is %q, want %q", got, "images")
		}
		test.SafeFprint(w, `[{"user":{"username":"justin","ids":{"slug":"justin","trakt":1}}},{"user":{"username":"sean","ids":{"slug":"sean","trakt":2}}}]`)
	})

	output := filepath.Join(t.TempDir(), "export_team_members.json")
	options := &str.Options{Action: consts.Members, ExtendedInfo: "images", Output: output}
	test.AssertNilError(t, TeamMembersHandler{}.Handle(options, s.Client))

	data, err := os.ReadFile(output)
	test.AssertNilError(t, err)
	got := []*str.TeamMember{}
	test.AssertNilError(t, json.Unmarshal(data, &got))
	test.AssertNoDiff(t, []*str.TeamMember{
		{User: &str.UserProfile{Username: str.String("justin"), IDs: &str.IDs{Slug: str.String("justin"), Trakt: test.Ptr(int64(1))}}},
		{User: &str.UserProfile{Username: str.String("sean"), IDs: &str.IDs{Slug: str.String("sean"), Trakt: test.Ptr(int64(2))}}},
	}, got)
}

func TestTeamMembersHandlerErrors(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr string
	}{
		{name: "empty list", status: http.StatusOK, body: `[]`, wantErr: consts.EmptyResult},
		{name: "server error", status: http.StatusInternalServerError, body: `{}`, wantErr: "fetch team members error"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			s.Mux.HandleFunc("/team", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			output := filepath.Join(t.TempDir(), "export_team_members.json")
			options := &str.Options{Action: consts.Members, Output: output}
			err := TeamMembersHandler{}.Handle(options, s.Client)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
			}
			if _, statErr := os.Stat(output); !os.IsNotExist(statErr) {
				t.Errorf("output file written on error: %v", statErr)
			}
		})
	}
}
