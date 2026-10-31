package trakt

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/mfederowicz/trakt-sync/deprecation"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// lockRegisteredRoutes makes every registered endpoint answer as locked until the test ends
func lockRegisteredRoutes(t *testing.T) {
	t.Helper()
	orig := lookupRoute
	lookupRoute = func(method, route string) (deprecation.Notice, bool) {
		notice, ok := orig(method, route)
		notice.Locked = "2026-11-01"
		return notice, ok
	}
	t.Cleanup(func() { lookupRoute = orig })
}

// TestLockedEndpointsSendNoRequest checks that every method of a registered endpoint, once the endpoint is
// locked, returns an *EndpointRemovedError with its notice and never reaches the API.
func TestLockedEndpointsSendNoRequest(t *testing.T) {
	lockRegisteredRoutes(t)
	ctx := context.Background()
	tests := []struct {
		name   string
		method string
		route  string
		reason string // the reason of the 2026-10-13 routes when not set
		call   func(c *Client) (*str.Response, error)
	}{
		{name: "RefreshMovieMetadata", method: http.MethodPost, route: "movies/{id}/refresh", call: func(c *Client) (*str.Response, error) {
			return c.Movies.RefreshMovieMetadata(ctx, "tron-legacy-2010")
		}},
		{name: "RefreshShowMetadata", method: http.MethodPost, route: "shows/{id}/refresh", call: func(c *Client) (*str.Response, error) {
			return c.Shows.RefreshShowMetadata(ctx, "breaking-bad")
		}},
		{name: "RefreshPersonMetadata", method: http.MethodPost, route: "people/{id}/refresh", call: func(c *Client) (*str.Response, error) {
			return c.People.RefreshPersonMetadata(ctx, "bryan-cranston")
		}},
		{name: "GetListsContainingThisPerson", method: http.MethodGet, route: "people/{id}/lists/{type}/{sort}",
			reason: "Trakt deprecated this endpoint: lists no longer hold people, so it always returns no lists", call: func(c *Client) (*str.Response, error) {
				_, resp, err := c.People.GetListsContainingThisPerson(ctx, "bryan-cranston", "personal", "popular", &uri.ListOptions{})
				return resp, err
			}},
		{name: "RemoveRecentSearch", method: http.MethodPost, route: "search/recent/remove", call: func(c *Client) (*str.Response, error) {
			return c.Search.RemoveRecentSearch(ctx, &str.RecentSearch{Query: "reacher", ID: 139606, Type: "shows"})
		}},
		{name: "AddSavedFilters", method: http.MethodPost, route: "users/saved_filters", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Users.AddSavedFilters(ctx, []*str.SavedFilterAdd{{Name: str.String("Sci-fi")}})
			return resp, err
		}},
		{name: "DeleteSavedFilter", method: http.MethodDelete, route: "users/saved_filters/{id}", call: func(c *Client) (*str.Response, error) {
			return c.Users.DeleteSavedFilter(ctx, 101)
		}},
		{name: "UpdateFavorites", method: http.MethodPut, route: "sync/favorites", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.Sync.UpdateFavorites(ctx, new(str.PersonalList))
			return resp, err
		}},
		{name: "GetSocialMovieRecommendations", method: http.MethodGet, route: "social_recommendations/movies", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.SocialRecommendations.GetSocialMovieRecommendations(ctx, &uri.ListOptions{})
			return resp, err
		}},
		{name: "GetSocialShowRecommendations", method: http.MethodGet, route: "social_recommendations/shows", call: func(c *Client) (*str.Response, error) {
			_, resp, err := c.SocialRecommendations.GetSocialShowRecommendations(ctx, &uri.ListOptions{})
			return resp, err
		}},
	}
	if len(tests) != len(deprecation.Notices()) {
		t.Errorf("%d methods tested, %d notices registered", len(tests), len(deprecation.Notices()))
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			setup := Setup()
			defer setup.Teardown()
			var requests atomic.Int32
			setup.Mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
				requests.Add(1)
				w.WriteHeader(http.StatusNotFound)
			})

			resp, err := tt.call(setup.Client)
			if resp != nil {
				t.Errorf("response is %v, want nil", resp)
			}
			removed := new(EndpointRemovedError)
			if !errors.As(err, &removed) {
				t.Fatalf("error is %v, want *EndpointRemovedError", err)
			}
			if removed.Notice.Method != tt.method || removed.Notice.Route != tt.route {
				t.Errorf("notice is for %s %s, want %s %s", removed.Notice.Method, removed.Notice.Route, tt.method, tt.route)
			}
			reason := "Trakt dropped this endpoint from its API contract, it may stop working"
			if tt.reason != "" {
				reason = tt.reason
			}
			want := tt.method + " " + tt.route + " is locked since 2026-11-01: " + reason + "; no request was sent"
			if err.Error() != want {
				t.Errorf("error is %q, want %q", err.Error(), want)
			}
			if got := requests.Load(); got != 0 {
				t.Errorf("%d requests sent, want 0", got)
			}
		})
	}
}

// TestLockedEndpointLeavesOthersAlone checks that routes sharing a path or prefix with a registered one still pass
// when the registered ones are locked, and that a registered route passes while it is only deprecated.
func TestLockedEndpointLeavesOthersAlone(t *testing.T) {
	if err := lockedEndpoint(http.MethodPost, "movies/{id}/refresh"); err != nil {
		t.Errorf("deprecated route is locked: %v", err)
	}
	lockRegisteredRoutes(t)
	if err := lockedEndpoint(http.MethodPost, "movies/{id}/refresh"); err == nil {
		t.Error("locked route passes")
	}
	for _, tt := range []struct{ method, route string }{
		{http.MethodPost, "sync/favorites"},
		{http.MethodGet, "users/saved_filters/{section}"},
		{http.MethodPost, "movies/{id}/refresh/justwatch"},
		{http.MethodPost, "search/recent"},
	} {
		if err := lockedEndpoint(tt.method, tt.route); err != nil {
			t.Errorf("%s %s: %v, want nil", tt.method, tt.route, err)
		}
	}
}
