// Package deprecation keeps the history of module actions whose Trakt endpoint is going away or gone
package deprecation

import (
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
)

const (
	endpointDropped = "Trakt dropped this endpoint from its API contract, it may stop working"
	// day20261013 is the day the routes left the Trakt API contract; the live API still served them that day
	day20261013 = "2026-10-13"
)

// notices lists every deprecated endpoint with the module action that calls it.
// A stage change is a new date on the row: Locked stops the requests, Removed marks the code as gone.
var notices = []Notice{
	{Module: consts.Movies, Action: consts.Refresh, Method: http.MethodPost, Route: "movies/{id}/refresh", Reason: endpointDropped, Deprecated: day20261013},
	{Module: consts.People, Action: consts.Refresh, Method: http.MethodPost, Route: "people/{id}/refresh", Reason: endpointDropped, Deprecated: day20261013},
	{Module: consts.Search, Action: consts.RemoveRecent, Method: http.MethodPost, Route: "search/recent/remove", Reason: endpointDropped, Deprecated: day20261013},
	{Module: consts.Shows, Action: consts.Refresh, Method: http.MethodPost, Route: "shows/{id}/refresh", Reason: endpointDropped, Deprecated: day20261013},
	{Module: consts.SocialRecommendations, Action: consts.Movies, Method: http.MethodGet, Route: "social_recommendations/movies", Reason: endpointDropped, Deprecated: day20261013},
	{Module: consts.SocialRecommendations, Action: consts.Shows, Method: http.MethodGet, Route: "social_recommendations/shows", Reason: endpointDropped, Deprecated: day20261013},
	{Module: consts.Sync, Action: consts.UpdateFavorites, Method: http.MethodPut, Route: "sync/favorites", Reason: endpointDropped, Deprecated: day20261013},
	{Module: consts.Users, Action: consts.AddSavedFilters, Method: http.MethodPost, Route: "users/saved_filters", Reason: endpointDropped, Deprecated: day20261013},
	{Module: consts.Users, Action: consts.DeleteSavedFilter, Method: http.MethodDelete, Route: "users/saved_filters/{id}", Reason: endpointDropped, Deprecated: day20261013},
}

// Notices returns a copy of every registered notice
func Notices() []Notice {
	return append([]Notice{}, notices...)
}

// Lookup returns the notice for a module action, if it is deprecated
func Lookup(module, action string) (Notice, bool) {
	for _, n := range notices {
		if n.Module == module && n.Action == action {
			return n, true
		}
	}
	return Notice{}, false
}

// LookupRoute returns the notice for an endpoint, if it is deprecated
func LookupRoute(method, route string) (Notice, bool) {
	for _, n := range notices {
		if n.Method == method && n.Route == route {
			return n, true
		}
	}
	return Notice{}, false
}
