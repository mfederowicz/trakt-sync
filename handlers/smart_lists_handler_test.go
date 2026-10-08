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
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/trakt"
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
		{name: "summary empty object", handler: SmartListsSummaryHandler{}, options: str.Options{Action: consts.Summary, InternalID: "no-such-list"}, status: http.StatusOK, body: `{}`,
			wantErr: "not found smart list for:no-such-list"},
		{name: "summary without id", handler: SmartListsSummaryHandler{}, options: str.Options{Action: consts.Summary}, wantErr: consts.EmptySmartListIDMsg, wantNoAPI: true},
		{name: "items", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi"}, status: http.StatusOK,
			body: `[{"rank":1,"type":"movie","movie":{"title":"Arrival"}}]`},
		{name: "items empty", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi"}, status: http.StatusOK, body: `[]`,
			wantErr: consts.EmptyResult},
		{name: "items not found", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi"}, status: http.StatusNotFound, body: `{}`,
			wantErr: "not found smart list for:top-sci-fi"},
		{name: "items without id", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items}, wantErr: consts.EmptySmartListIDMsg, wantNoAPI: true},
		{name: "items invalid type", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi", Type: "episodes"},
			wantErr: "type 'episodes' is not valid", wantNoAPI: true},
		{name: "items invalid sort_by", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi", SortBy: "popularity"},
			wantErr: "sort_by 'popularity' is not valid", wantNoAPI: true},
		{name: "items invalid sort_how", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi", SortHow: "up"},
			wantErr: "sort_how 'up' is not valid", wantNoAPI: true},
		{name: "items invalid watchnow_country", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi", WatchNowCountry: "USA"},
			wantErr: "watchnow_country 'USA' is not valid", wantNoAPI: true},
		{name: "items invalid parental range", handler: SmartListsItemsHandler{}, options: str.Options{Action: consts.Items, InternalID: "top-sci-fi", Parental: str.ParentalGuide{Violence: "0-4"}},
			wantErr: "parental_violence '0-4' is not valid", wantNoAPI: true},
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
		w.Header().Set(trakt.HeaderPaginationPage, page)
		w.Header().Set(trakt.HeaderPaginationPageCount, "2")
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

func TestSmartListsItemsHandlerRoutes(t *testing.T) {
	tests := []struct {
		name      string
		options   str.Options
		wantPath  string
		wantQuery string
	}{
		{name: "no type or sort uses the plain route", options: str.Options{}, wantPath: "/smart-lists/mixed/items"},
		{name: "type only", options: str.Options{Type: "shows"}, wantPath: "/smart-lists/mixed/items/shows/rank/asc"},
		{name: "sort_by only", options: str.Options{SortBy: "imdb_rating"}, wantPath: "/smart-lists/mixed/items/all/imdb_rating/desc"},
		{name: "sort_by title defaults to asc", options: str.Options{SortBy: "title"}, wantPath: "/smart-lists/mixed/items/all/title/asc"},
		{name: "sort_how only", options: str.Options{SortHow: "desc"}, wantPath: "/smart-lists/mixed/items/all/rank/desc"},
		{name: "all three", options: str.Options{Type: "movies", SortBy: "added", SortHow: "asc"}, wantPath: "/smart-lists/mixed/items/movies/added/asc"},
		{name: "watchnow country and parental filters", wantPath: "/smart-lists/mixed/items",
			options: str.Options{WatchNow: "free", WatchNowCountry: "pl",
				Parental: str.ParentalGuide{Nudity: "0-1", Violence: "0-2", Profanity: "0-1", Alcohol: "0-3", Frightening: "1-2", IncludeUnrated: true}},
			wantQuery: "parental_alcohol=0-3&parental_frightening=1-2&parental_include_unrated=true&parental_nudity=0-1&parental_profanity=0-1&parental_violence=0-2&watchnow=free&watchnow_country=pl"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			s.Mux.HandleFunc("/smart-lists/", func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				if r.URL.Path != tt.wantPath {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.wantPath)
				}
				q := r.URL.Query()
				q.Del("page")
				if got := q.Encode(); got != tt.wantQuery {
					t.Errorf("query is %q, want %q", got, tt.wantQuery)
				}
				test.SafeFprint(w, `[{"rank":1,"type":"movie","movie":{"title":"Arrival"}}]`)
			})

			options := tt.options
			options.Action = consts.Items
			options.InternalID = "mixed"
			options.Output = filepath.Join(t.TempDir(), "out.json")
			test.AssertNilError(t, SmartListsItemsHandler{}.Handle(&options, s.Client))
		})
	}
}
