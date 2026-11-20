// Package trakt is a Go client for the Trakt API
package trakt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersServiceUpdateSettings(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	want := &str.SettingsUpdate{
		User:     &str.SettingsUpdateUser{Location: str.String("Warsaw"), Private: test.Ptr(true)},
		Browsing: &str.SettingsBrowsing{Spoilers: &str.SettingsSpoilers{Episodes: str.String("hide")}, WatchNow: &str.SettingsWatchNow{Country: str.String("pl")}},
	}
	setup.Mux.HandleFunc("/users/settings", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPut)
		got := new(str.SettingsUpdate)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, want, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Users.UpdateSettings(context.Background(), want)
	test.AssertNilError(t, err)
	if resp == nil || resp.StatusCode != http.StatusCreated {
		t.Errorf("response is %v, want status %d", resp, http.StatusCreated)
	}
}
