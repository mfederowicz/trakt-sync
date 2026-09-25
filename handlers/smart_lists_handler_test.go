// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
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

func TestSmartListsHandlers(t *testing.T) {
	tests := []struct {
		name      string
		handler   Handler
		options   str.Options
		status    int
		body      string
		wantErr   string
		wantNoAPI bool
	}{
		{name: "summary", handler: SmartListsSummaryHandler{}, options: str.Options{Action: consts.Summary, InternalID: "top-sci-fi"}, status: http.StatusOK, body: `{"name":"Top Sci-Fi"}`},
		{name: "summary not found", handler: SmartListsSummaryHandler{}, options: str.Options{Action: consts.Summary, InternalID: "top-sci-fi"}, status: http.StatusNotFound, body: `{}`,
			wantErr: "not found smart list for:top-sci-fi"},
		{name: "summary without id", handler: SmartListsSummaryHandler{}, options: str.Options{Action: consts.Summary}, wantErr: consts.EmptySmartListIDMsg, wantNoAPI: true},
		{name: "items", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi"}, status: http.StatusOK,
			body: `[{"rank":1,"type":"movie","movie":{"title":"Arrival"}}]`},
		{name: "items empty", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi"}, status: http.StatusOK, body: `[]`,
			wantErr: consts.EmptyResult},
		{name: "items not found", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi"}, status: http.StatusNotFound, body: `{}`,
			wantErr: "not found smart list for:top-sci-fi"},
		{name: "items without id", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items}, wantErr: consts.EmptySmartListIDMsg, wantNoAPI: true},
		{name: "items invalid watchnow", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi", WatchNow: "cable"},
			wantErr: "watchnow 'cable' is not valid", wantNoAPI: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/smart-lists/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, http.MethodGet)
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			options := tt.options
			options.Output = filepath.Join(t.TempDir(), "out.json")
			err := tt.handler.Handle(&options, s.Client)
			if tt.wantNoAPI && calls != 0 {
				t.Error("API was called with invalid options")
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				if _, statErr := os.Stat(options.Output); !os.IsNotExist(statErr) {
					t.Errorf("output file written on error: %v", statErr)
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

func TestSmartListsItemsHandlerFetchesAllPages(t *testing.T) {
	s := setup(t)
	defer s.Teardown()

	pages := map[string]string{
		"1": `[{"rank":1,"type":"show","show":{"title":"Andor"}}]`,
		"2": `[{"rank":2,"type":"show","show":{"title":"Severance"}}]`,
	}
	s.Mux.HandleFunc("/smart-lists/best-shows/items", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("genres") != "drama" || q.Get("ignore_watched") != "true" || q.Get("limit") != "1" {
			t.Errorf("query is %q, want genres=drama, ignore_watched=true and limit=1", r.URL.RawQuery)
		}
		page := q.Get("page")
		w.Header().Set(internal.HeaderPaginationPage, page)
		w.Header().Set(internal.HeaderPaginationPageCount, "2")
		test.SafeFprint(w, pages[page])
	})

	output := filepath.Join(t.TempDir(), "export_smart_lists_items_best-shows.json")
	options := &str.Options{Action: consts.Items, InternalID: "best-shows", PerPage: 1, Genres: "drama", IgnoreWatched: "true", Output: output}
	test.AssertNilError(t, SmartListsItemsHandler{}.Handle(options, s.Client))

	data, err := os.ReadFile(output)
	test.AssertNilError(t, err)
	got := []*str.UserListItem{}
	test.AssertNilError(t, json.Unmarshal(data, &got))
	test.AssertNoDiff(t, []*str.UserListItem{
		{Rank: test.Ptr(1), Type: str.String("show"), Show: &str.Show{Title: str.String("Andor")}},
		{Rank: test.Ptr(2), Type: str.String("show"), Show: &str.Show{Title: str.String("Severance")}},
	}, got)
}
