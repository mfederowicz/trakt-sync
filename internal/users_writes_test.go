// Package internal used for client and services
package internal

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

func TestUsersServiceAddSavedFilters(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	filters := []*str.SavedFilterAdd{{Name: str.String("Sci-fi"), URL: str.String("https://trakt.tv/movies/trending?genres=science-fiction")}}
	setup.Mux.HandleFunc("/users/saved_filters", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := []*str.SavedFilterAdd{}
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(&got))
		test.AssertNoDiff(t, filters, got)
		w.WriteHeader(http.StatusCreated)
		test.SafeFprint(w, `{"added":[{"rank":1,"id":101,"section":"movies","name":"Sci-fi","path":"/movies/trending","query":"genres=science-fiction"}],"skipped":[]}`)
	})

	got, _, err := setup.Client.Users.AddSavedFilters(context.Background(), filters)
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, &str.SavedFiltersResult{
		Added: []*str.SavedFilter{{Rank: test.Ptr(1), ID: test.Ptr(int64(101)), Section: str.String("movies"), Name: str.String("Sci-fi"),
			Path: str.String("/movies/trending"), Query: str.String("genres=science-fiction")}},
		Skipped: []*str.SavedFilterAdd{},
	}, got)
}

func TestUsersServiceDeleteSavedFilter(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/users/saved_filters/101", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := setup.Client.Users.DeleteSavedFilter(context.Background(), 101)
	test.AssertNilError(t, err)
	if resp == nil || resp.StatusCode != http.StatusNoContent {
		t.Errorf("response is %v, want status %d", resp, http.StatusNoContent)
	}
}
