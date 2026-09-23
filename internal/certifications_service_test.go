// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestCertificationsServiceGetCertifications(t *testing.T) {
	tests := []struct {
		name    string
		strType string
		body    string
		want    *str.Certifications
	}{
		{
			name:    "movie certifications",
			strType: "movies",
			body:    `{"us":[{"name":"PG","slug":"pg","description":"Parental Guidance Suggested"}]}`,
			want: &str.Certifications{Us: []*str.Certification{
				{Name: str.String("PG"), Slug: str.String("pg"), Description: str.String("Parental Guidance Suggested")},
			}},
		},
		{
			name:    "show certifications",
			strType: "shows",
			body:    `{"us":[{"name":"TV-MA","slug":"tv-ma","description":"Mature Audience Only"}]}`,
			want: &str.Certifications{Us: []*str.Certification{
				{Name: str.String("TV-MA"), Slug: str.String("tv-ma"), Description: str.String("Mature Audience Only")},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc("/certifications/"+tt.strType, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				test.SafeFprint(w, tt.body)
			})

			got, _, err := setup.Client.Certifications.GetCertifications(context.Background(), str.String(tt.strType))
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, tt.want, got)
		})
	}
}
