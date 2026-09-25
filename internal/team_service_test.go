// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestTeamServiceGetTeamMembers(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/team", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		if got := r.URL.RawQuery; got != "extended=full%2Cimages" {
			t.Errorf("query is %q, want %q", got, "extended=full%2Cimages")
		}
		test.SafeFprint(w, `[{"user":{"username":"justin","private":false,"deleted":false,"name":"Justin Nemeth",`+
			`"vip":true,"vip_ep":true,"director":true,"ids":{"slug":"justin","trakt":1},"location":"San Diego, CA",`+
			`"images":{"avatar":{"full":"https://example.com/justin.png"}}}}]`)
	})

	got, _, err := setup.Client.Team.GetTeamMembers(context.Background(), &uri.ListOptions{Extended: "full,images"})
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []*str.TeamMember{
		{User: &str.UserProfile{
			Username: str.String("justin"),
			Private:  test.Ptr(false),
			Deleted:  test.Ptr(false),
			Name:     str.String("Justin Nemeth"),
			Vip:      test.Ptr(true),
			VipEp:    test.Ptr(true),
			Director: test.Ptr(true),
			IDs:      &str.IDs{Slug: str.String("justin"), Trakt: test.Ptr(int64(1))},
			Location: str.String("San Diego, CA"),
			Images:   &str.Images{Avatar: &str.Avatar{Full: str.String("https://example.com/justin.png")}},
		}},
	}, got)
}

func TestTeamServiceGetTeamMembersError(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/team", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	got, resp, err := setup.Client.Team.GetTeamMembers(context.Background(), &uri.ListOptions{})
	if err == nil {
		t.Fatal("expected an error")
	}
	if got != nil {
		t.Errorf("list is %v, want nil", got)
	}
	if resp == nil || resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("response is %v, want status %d", resp, http.StatusInternalServerError)
	}
}
