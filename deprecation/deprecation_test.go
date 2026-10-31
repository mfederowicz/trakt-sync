// Package deprecation keeps the history of module actions whose Trakt endpoint is going away or gone
package deprecation

import (
	"net/http"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/stretchr/testify/assert"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		name   string
		module string
		action string
		method string
		route  string
		want   bool
	}{
		{name: "movies refresh", module: consts.Movies, action: consts.Refresh, method: http.MethodPost, route: "movies/{id}/refresh", want: true},
		{name: "shows refresh", module: consts.Shows, action: consts.Refresh, method: http.MethodPost, route: "shows/{id}/refresh", want: true},
		{name: "people refresh", module: consts.People, action: consts.Refresh, method: http.MethodPost, route: "people/{id}/refresh", want: true},
		{name: "people lists", module: consts.People, action: consts.Lists, method: http.MethodGet, route: "people/{id}/lists/{type}/{sort}", want: true},
		{name: "search remove_recent", module: consts.Search, action: consts.RemoveRecent, method: http.MethodPost, route: "search/recent/remove", want: true},
		{name: "users add_saved_filters", module: consts.Users, action: consts.AddSavedFilters, method: http.MethodPost, route: "users/saved_filters", want: true},
		{name: "users delete_saved_filter", module: consts.Users, action: consts.DeleteSavedFilter, method: http.MethodDelete, route: "users/saved_filters/{id}", want: true},
		{name: "sync update_favorites", module: consts.Sync, action: consts.UpdateFavorites, method: http.MethodPut, route: "sync/favorites", want: true},
		{name: "social movies", module: consts.SocialRecommendations, action: consts.Movies, method: http.MethodGet, route: "social_recommendations/movies", want: true},
		{name: "social shows", module: consts.SocialRecommendations, action: consts.Shows, method: http.MethodGet, route: "social_recommendations/shows", want: true},
		{name: "younify refresh", module: consts.Younify, action: consts.Refresh, method: http.MethodPost, route: "younify/users/refresh/{service_id}"},
		{name: "movies refresh_justwatch", module: consts.Movies, action: consts.RefreshJustwatch, method: http.MethodPost, route: "movies/{id}/refresh/justwatch"},
		{name: "movies without action", module: consts.Movies, action: consts.EmptyString, method: http.MethodGet, route: "movies/{id}"},
		{name: "sync add_to_favorites, same path", module: consts.Sync, action: consts.AddToFavorites, method: http.MethodPost, route: "sync/favorites"},
		{name: "users saved_filters, same prefix", module: consts.Users, action: consts.SavedFilters, method: http.MethodGet, route: "users/saved_filters/{section}"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			byAction, ok := Lookup(tt.module, tt.action)
			assert.Equal(t, tt.want, ok)
			byRoute, ok := LookupRoute(tt.method, tt.route)
			assert.Equal(t, tt.want, ok)
			if tt.want {
				assert.Equal(t, byAction, byRoute)
				assert.NotEmpty(t, byAction.Deprecated)
				assert.NotEmpty(t, byAction.Reason)
			}
		})
	}
}

func TestNoticeStages(t *testing.T) {
	base := Notice{Module: consts.Movies, Action: consts.Refresh, Method: http.MethodPost, Route: "movies/{id}/refresh",
		Reason: "Trakt removed this endpoint", Deprecated: "2026-10-13"}
	locked := base
	locked.Locked = "2026-11-01"
	removed := locked
	removed.Removed = "2026-12-01"

	tests := []struct {
		name        string
		notice      Notice
		wantStage   Stage
		wantBlocked bool
		wantMessage string
		wantRoute   string
	}{
		{name: "deprecated", notice: base, wantStage: StageDeprecated,
			wantMessage: "movies -a refresh is deprecated since 2026-10-13: Trakt removed this endpoint",
			wantRoute:   "POST movies/{id}/refresh is deprecated since 2026-10-13: Trakt removed this endpoint"},
		{name: "locked", notice: locked, wantStage: StageLocked, wantBlocked: true,
			wantMessage: "movies -a refresh is locked since 2026-11-01: Trakt removed this endpoint; no request was sent",
			wantRoute:   "POST movies/{id}/refresh is locked since 2026-11-01: Trakt removed this endpoint; no request was sent"},
		{name: "removed", notice: removed, wantStage: StageRemoved, wantBlocked: true,
			wantMessage: "movies -a refresh is removed since 2026-12-01: Trakt removed this endpoint; no request was sent",
			wantRoute:   "POST movies/{id}/refresh is removed since 2026-12-01: Trakt removed this endpoint; no request was sent"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantStage, tt.notice.Stage())
			assert.Equal(t, tt.wantBlocked, tt.notice.Blocked())
			assert.Equal(t, tt.wantMessage, tt.notice.Message())
			assert.Equal(t, tt.wantRoute, tt.notice.RouteMessage())
		})
	}
}

func TestNoticesIsACopy(t *testing.T) {
	list := Notices()
	assert.Len(t, list, len(notices))
	list[consts.ZeroValue].Module = "changed"
	assert.NotEqual(t, "changed", notices[consts.ZeroValue].Module)
}
