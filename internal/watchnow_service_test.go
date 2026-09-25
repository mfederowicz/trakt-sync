// Package internal used for client and services
package internal

import (
	"context"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestWatchNowServiceGetSources(t *testing.T) {
	body := `[{"us":[{"source":"netflix","name":"Netflix","free":false,"cinema":false,"amazon":false,"color":"#e50914",` +
		`"link_count":1200,"images":{"logo":"netflix.png","channel":null}}]}]`
	want := []map[string][]*str.WatchNowSource{{"us": {{
		Source:    str.String("netflix"),
		Name:      str.String("Netflix"),
		Free:      test.Ptr(false),
		Cinema:    test.Ptr(false),
		Amazon:    test.Ptr(false),
		Color:     str.String("#e50914"),
		LinkCount: test.Ptr(1200),
		Images:    &str.WatchNowSourceImages{Logo: str.String("netflix.png")},
	}}}}

	tests := []struct {
		name string
		path string
		call func(s *WatchNowService) ([]map[string][]*str.WatchNowSource, *str.Response, error)
	}{
		{
			name: "all",
			path: "/watchnow/sources",
			call: func(s *WatchNowService) ([]map[string][]*str.WatchNowSource, *str.Response, error) {
				return s.GetWatchNowSources(context.Background())
			},
		},
		{
			name: "country",
			path: "/watchnow/sources/us",
			call: func(s *WatchNowService) ([]map[string][]*str.WatchNowSource, *str.Response, error) {
				return s.GetWatchNowSourcesByCountry(context.Background(), str.String("us"))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			calls := 0
			setup.Mux.HandleFunc("/watchnow/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				if r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				if r.URL.RawQuery != "" {
					t.Errorf("query is %q, want none", r.URL.RawQuery)
				}
				test.SafeFprint(w, body)
			})

			got, _, err := tt.call(setup.Client.WatchNow)
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, want, got)
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
		})
	}
}

func TestWatchNowServiceGetSourcesError(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/watchnow/sources", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	got, resp, err := setup.Client.WatchNow.GetWatchNowSources(context.Background())
	if err == nil {
		t.Fatal("expected an error")
	}
	if got != nil {
		t.Errorf("list is %v, want nil", got)
	}
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Errorf("response is %v, want status %d", resp, http.StatusForbidden)
	}
}
