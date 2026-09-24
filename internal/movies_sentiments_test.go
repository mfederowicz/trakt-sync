// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestMoviesServiceGetMovieSentiments(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/movies/tron-legacy-2010/sentiments", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		test.SafeFprint(w, `{"good":[{"sentiment":"great visuals","comment_ids":[1,2]}],"bad":[{"sentiment":"weak story"}],"comment_count":2}`)
	})

	got, _, err := setup.Client.Movies.GetMovieSentiments(context.Background(), str.String("tron-legacy-2010"))
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, &str.Sentiments{
		Good:         []*str.Sentiment{{Sentiment: str.String("great visuals"), CommentIDs: &[]int{1, 2}}},
		Bad:          []*str.Sentiment{{Sentiment: str.String("weak story")}},
		CommentCount: test.Ptr(2),
	}, got)
}
