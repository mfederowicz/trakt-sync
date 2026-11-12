package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestUsersServiceGetStatsTotals(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/users/sean/stats", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `{"progress":{"started":12,"finished":7,"dropped":2},"lists":5,"total_minutes":98765,"total_plays":4321}`)
	})

	got, _, err := setup.Client.Users.GetStats(context.Background(), "sean")
	test.AssertNilError(t, err)
	want := &str.UserStats{
		Progress:     &str.UserStatsProgress{Started: test.Ptr(12), Finished: test.Ptr(7), Dropped: test.Ptr(2)},
		Lists:        test.Ptr(5),
		TotalMinutes: test.Ptr(98765),
		TotalPlays:   test.Ptr(4321),
	}
	test.AssertNoDiff(t, want, got)
}

func TestUsersServiceGetStatsNullTotals(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/users/me/stats", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `{"progress":null,"lists":null,"total_minutes":null,"total_plays":null}`)
	})

	got, _, err := setup.Client.Users.GetStats(context.Background(), "")
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, &str.UserStats{}, got)
}
