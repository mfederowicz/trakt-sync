// Package internal used for client and services
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/mfederowicz/trakt-sync/uri"
)

func TestListsServiceListsByType(t *testing.T) {
	tests := []struct {
		name string
		path string
		call func(s *ListsService, opts *uri.ListOptions) ([]*str.List, *str.Response, error)
	}{
		{
			name: "trending",
			path: "/lists/trending/personal",
			call: func(s *ListsService, opts *uri.ListOptions) ([]*str.List, *str.Response, error) {
				return s.GetTrendingListsByType(context.Background(), str.String("personal"), opts)
			},
		},
		{
			name: "popular",
			path: "/lists/popular/official",
			call: func(s *ListsService, opts *uri.ListOptions) ([]*str.List, *str.Response, error) {
				return s.GetPopularListsByType(context.Background(), str.String("official"), opts)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()

			setup.Mux.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodGet)
				if got, want := r.URL.Query().Get("page"), "2"; got != want {
					t.Errorf("page query is %q, want %q", got, want)
				}
				test.SafeFprint(w, `[{"like_count":5,"comment_count":1}]`)
			})

			got, _, err := tt.call(setup.Client.Lists, &uri.ListOptions{Page: 2})
			test.AssertNilError(t, err)
			test.AssertNoDiff(t, []*str.List{{LikeCount: test.Ptr(5), CommentCount: test.Ptr(1)}}, got)
		})
	}
}

func TestListsServiceGetListItemsSortQuery(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/lists/55/items/movie", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		q := r.URL.Query()
		if q.Get("sort_by") != "added" || q.Get("sort_how") != "desc" {
			t.Errorf("sort query is %q, want sort_by=added&sort_how=desc", r.URL.RawQuery)
		}
		test.SafeFprint(w, `[]`)
	})

	_, _, err := setup.Client.Lists.GetListItems(context.Background(), str.String("55"), str.String("movie"), &uri.ListOptions{SortBy: "added", SortHow: "desc"})
	test.AssertNilError(t, err)
}

func TestListsServiceReportList(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	setup.Mux.HandleFunc("/lists/55/report", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodPost)
		got := new(str.ListReport)
		test.AssertNilError(t, json.NewDecoder(r.Body).Decode(got))
		test.AssertNoDiff(t, &str.ListReport{Reason: str.String("spam"), Message: str.String("ads")}, got)
		w.WriteHeader(http.StatusCreated)
	})

	resp, err := setup.Client.Lists.ReportList(context.Background(), str.String("55"), &str.ListReport{Reason: str.String("spam"), Message: str.String("ads")})
	test.AssertNilError(t, err)
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("status code is %d, want %d", got, want)
	}
}

func TestListsServiceReportListConflict(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()
	conflictMux(t, setup.Mux, "/lists/55/report", `{}`)

	_, err := setup.Client.Lists.ReportList(context.Background(), str.String("55"), &str.ListReport{Reason: str.String("spam")})
	if got, want := fmt.Sprint(err), fmt.Sprintf(consts.ListReportPending, "55"); got != want {
		t.Errorf("error is %q, want %q", got, want)
	}
}
