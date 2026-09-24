// Package handlers used to handle module actions
package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestSearchHandlers(t *testing.T) {
	tests := []struct {
		name      string
		handler   Handler
		options   str.Options
		wantCalls map[string]int
		wantErr   string
	}{
		{
			name:      "text query",
			handler:   SearchTextQueryHandler{},
			options:   str.Options{Action: consts.TextQuery, SearchType: str.Slice{"movie"}, SearchField: str.Slice{"title"}, Query: "freddy"},
			wantCalls: map[string]int{"GET /search/movie?field=title&page=1&query=freddy": 1},
		},
		{
			name:      "text query without type",
			handler:   SearchTextQueryHandler{},
			options:   str.Options{Action: consts.TextQuery, Query: "freddy"},
			wantCalls: map[string]int{},
			wantErr:   "invalid -t flag values",
		},
		{
			name:      "text query field not valid for type",
			handler:   SearchTextQueryHandler{},
			options:   str.Options{Action: consts.TextQuery, SearchType: str.Slice{"person"}, SearchField: str.Slice{"tagline"}, Query: "freddy"},
			wantCalls: map[string]int{},
			wantErr:   "invalid --field flag values",
		},
		{
			name:      "id lookup",
			handler:   SearchIDLookupHandler{},
			options:   str.Options{Action: consts.IDLookup, SearchIDType: "imdb", ID: "tt0266697", SearchType: str.Slice{"movie"}},
			wantCalls: map[string]int{"GET /search/imdb/tt0266697?type=movie": 1},
		},
		{
			name:      "id lookup invalid id_type",
			handler:   SearchIDLookupHandler{},
			options:   str.Options{Action: consts.IDLookup, SearchIDType: "tvrage", ID: "75725"},
			wantCalls: map[string]int{},
			wantErr:   "invalid --id_type flag value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
				calls[r.Method+" "+r.URL.RequestURI()]++
				test.SafeFprint(w, `[{"type":"movie","score":1,"movie":{"title":"Freddy"}}]`)
			})

			options := tt.options
			options.Module = consts.Search
			options.Output = filepath.Join(t.TempDir(), "out.json")
			err := tt.handler.Handle(&options, s.Client)

			test.AssertNoDiff(t, tt.wantCalls, calls)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			test.AssertNilError(t, err)
			if _, statErr := os.Stat(options.Output); statErr != nil {
				t.Errorf("output file was not written: %v", statErr)
			}
		})
	}
}
