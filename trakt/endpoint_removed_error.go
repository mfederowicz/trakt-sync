package trakt

import "github.com/mfederowicz/trakt-sync/deprecation"

// EndpointRemovedError is returned, without sending a request, by a method whose Trakt endpoint is locked or removed
type EndpointRemovedError struct {
	// Notice holds the endpoint and its deprecation history
	Notice deprecation.Notice
}

func (r *EndpointRemovedError) Error() string {
	return r.Notice.RouteMessage()
}

// lookupRoute finds the deprecation notice of an endpoint; a variable so tests can lock a route
var lookupRoute = deprecation.LookupRoute

// lockedEndpoint returns an *EndpointRemovedError when requests to the endpoint are no longer sent
func lockedEndpoint(method, route string) error {
	notice, ok := lookupRoute(method, route)
	if !ok || !notice.Blocked() {
		return nil
	}
	return &EndpointRemovedError{Notice: notice}
}
