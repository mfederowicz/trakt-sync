// Package handlers used to handle module actions
package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
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
			wantCalls: map[string]int{"GET /search/movie?fields=title&page=1&query=freddy": 1},
		},
		{
			name:      "text query original_title",
			handler:   SearchTextQueryHandler{},
			options:   str.Options{Action: consts.TextQuery, SearchType: str.Slice{"movie", "show"}, SearchField: str.Slice{"original_title"}, Query: "freddy"},
			wantCalls: map[string]int{"GET /search/movie,show?fields=original_title&page=1&query=freddy": 1},
		},
		{
			name:      "text query show_title",
			handler:   SearchTextQueryHandler{},
			options:   str.Options{Action: consts.TextQuery, SearchType: str.Slice{"episode"}, SearchField: str.Slice{"show_title"}, Query: "freddy"},
			wantCalls: map[string]int{"GET /search/episode?fields=show_title&page=1&query=freddy": 1},
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
		{
			name:      "exact query",
			handler:   SearchExactQueryHandler{},
			options:   str.Options{Action: consts.ExactQuery, SearchType: str.Slice{"show"}, Query: "dark"},
			wantCalls: map[string]int{"GET /search/show/exact?page=1&query=dark": 1},
		},
		{
			name:      "exact query without query",
			handler:   SearchExactQueryHandler{},
			options:   str.Options{Action: consts.ExactQuery, SearchType: str.Slice{"movie"}},
			wantCalls: map[string]int{},
			wantErr:   consts.EmptySearchQueryMsg,
		},
		{
			name:      "exact query type not in contract",
			handler:   SearchExactQueryHandler{},
			options:   str.Options{Action: consts.ExactQuery, SearchType: str.Slice{"person"}, Query: "keanu"},
			wantCalls: map[string]int{},
			wantErr:   "set one -t value for exact_query",
		},
		{
			name:      "exact query two types",
			handler:   SearchExactQueryHandler{},
			options:   str.Options{Action: consts.ExactQuery, SearchType: str.Slice{"movie", "show"}, Query: "dark"},
			wantCalls: map[string]int{},
			wantErr:   "set one -t value for exact_query",
		},
		{
			name:      "trending",
			handler:   SearchTrendingHandler{},
			options:   str.Options{Action: consts.Trending, SearchType: str.Slice{"movies"}},
			wantCalls: map[string]int{"GET /search/recent_by_id/global/movies?page=1": 1},
		},
		{
			name:      "trending with query",
			handler:   SearchTrendingHandler{},
			options:   str.Options{Action: consts.Trending, SearchType: str.Slice{"people"}, Query: "keanu"},
			wantCalls: map[string]int{"GET /search/recent_by_id/global/people?page=1&query=keanu": 1},
		},
		{
			name:      "trending singular type",
			handler:   SearchTrendingHandler{},
			options:   str.Options{Action: consts.Trending, SearchType: str.Slice{"movie"}},
			wantCalls: map[string]int{},
			wantErr:   "set one -t value for trending",
		},
		{
			name:      "trending without type",
			handler:   SearchTrendingHandler{},
			options:   str.Options{Action: consts.Trending},
			wantCalls: map[string]int{},
			wantErr:   "set one -t value for trending",
		},
		{
			name:      "add recent",
			handler:   SearchAddRecentHandler{},
			options:   str.Options{Action: consts.AddRecent, SearchType: str.Slice{"shows"}, Query: "reacher", ID: "139606"},
			wantCalls: map[string]int{`POST /search/recent {"query":"reacher","id":139606,"type":"shows"}`: 1},
		},
		{
			name:      "remove recent",
			handler:   SearchRemoveRecentHandler{},
			options:   str.Options{Action: consts.RemoveRecent, SearchType: str.Slice{"lists"}, Query: "best of", ID: "55"},
			wantCalls: map[string]int{`POST /search/recent/remove {"query":"best of","id":55,"type":"lists"}`: 1},
		},
		{
			name:      "add recent singular type",
			handler:   SearchAddRecentHandler{},
			options:   str.Options{Action: consts.AddRecent, SearchType: str.Slice{"show"}, Query: "reacher", ID: "139606"},
			wantCalls: map[string]int{},
			wantErr:   "set one -t value for add_recent",
		},
		{
			name:      "add recent without query",
			handler:   SearchAddRecentHandler{},
			options:   str.Options{Action: consts.AddRecent, SearchType: str.Slice{"shows"}, ID: "139606"},
			wantCalls: map[string]int{},
			wantErr:   consts.EmptySearchQueryMsg,
		},
		{
			name:      "remove recent without id",
			handler:   SearchRemoveRecentHandler{},
			options:   str.Options{Action: consts.RemoveRecent, SearchType: str.Slice{"shows"}, Query: "reacher"},
			wantCalls: map[string]int{},
			wantErr:   consts.EmptySearchRecentIDMsg,
		},
		{
			name:      "remove recent slug id",
			handler:   SearchRemoveRecentHandler{},
			options:   str.Options{Action: consts.RemoveRecent, SearchType: str.Slice{"shows"}, Query: "reacher", ID: "reacher"},
			wantCalls: map[string]int{},
			wantErr:   "trakt id must be a positive number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := map[string]int{}
			s.Mux.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					body, _ := io.ReadAll(r.Body)
					calls[r.Method+" "+r.URL.RequestURI()+" "+strings.TrimSpace(string(body))]++
					w.WriteHeader(http.StatusCreated)
					return
				}
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
			_, statErr := os.Stat(options.Output)
			if wantFile := !strings.HasSuffix(options.Action, "_recent"); wantFile != (statErr == nil) {
				t.Errorf("output file written is %v, want %v", statErr == nil, wantFile)
			}
		})
	}
}

func TestSearchTrendingHandlerPages(t *testing.T) {
	s := setup(t)
	defer s.Teardown()

	pages := []string{}
	s.Mux.HandleFunc("/search/recent_by_id/global/shows", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		pages = append(pages, page)
		w.Header().Set(internal.HeaderPaginationPage, page)
		w.Header().Set(internal.HeaderPaginationPageCount, "2")
		test.SafeFprint(w, `[{"id":1,"count":3,"type":"show","show":{"title":"Dark"}}]`)
	})

	options := str.Options{Module: consts.Search, Action: consts.Trending, SearchType: str.Slice{"shows"}, Output: filepath.Join(t.TempDir(), "out.json")}
	err := SearchTrendingHandler{}.Handle(&options, s.Client)
	test.AssertNilError(t, err)
	test.AssertNoDiff(t, []string{"1", "2"}, pages)
}
