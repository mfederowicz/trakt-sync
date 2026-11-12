package trakt

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestPeopleServiceGetSinglePersonSocialIDs(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/people/bryan-cranston", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		if got, want := r.URL.RawQuery, "extended=full&language=pt"; got != want {
			t.Errorf("query is %q, want %q", got, want)
		}
		test.SafeFprint(w, `{"name":"Bryan Cranston","social_ids":{"twitter":"BryanCranston","facebook":null,"instagram":"bryancranston",`+
			`"tiktok":"bryancranston","youtube":"UCabc","wikidata":"Q23547","wikipedia":"Bryan_Cranston","wikipedia_language":"pt"}}`)
	})

	got, _, err := setup.Client.People.GetSinglePerson(context.Background(), "bryan-cranston", &uri.ListOptions{Extended: "full", Language: "pt"})
	test.AssertNilError(t, err)
	want := &str.Person{
		Name: str.String("Bryan Cranston"),
		SocialIDs: &str.SocialIDs{
			Twitter: str.String("BryanCranston"), Instagram: str.String("bryancranston"), Wikipedia: str.String("Bryan_Cranston"),
			TikTok: str.String("bryancranston"), YouTube: str.String("UCabc"), Wikidata: str.String("Q23547"), WikipediaLanguage: str.String("pt"),
		},
	}
	test.AssertNoDiff(t, want, got)
}
