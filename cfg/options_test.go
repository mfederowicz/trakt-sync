package cfg

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/stretchr/testify/assert"
)

func TestIsValidConfigTypeSlice(t *testing.T) {
	t.Helper()
	got := IsValidConfigTypeSlice([]string{"movie", "show"}, []string{"xxx"})
	if bool(got) {
		t.Fatalf("Expected %v, got %v", false, bool(got))
	}
}

func TestModuleConfigTypeComments(t *testing.T) {
	t.Helper()

	got := ModuleActionConfig["comments:trending"].Type

	assert.Equal(t, got, []string{"all", "movies", "shows", "seasons", "episodes", "lists"})

	got = ModuleActionConfig["comments:recent"].Type

	assert.Equal(t, got, []string{"all", "movies", "shows", "seasons", "episodes", "lists"})
}

func TestModuleConfigTypeUsers(t *testing.T) {
	t.Helper()

	got := ModuleActionConfig["users:watched"].Type

	assert.Equal(t, got, []string{"movies", "shows"})
}

func TestGetOutputForModuleListsTrendingPopular(t *testing.T) {
	tests := []struct {
		action, listType, want string
	}{
		{action: "trending", want: "export_lists_trending.json"},
		{action: "popular", want: "export_lists_popular.json"},
		{action: "trending", listType: "personal", want: "export_lists_trending_personal.json"},
		{action: "popular", listType: "official", want: "export_lists_popular_official.json"},
	}

	for _, tt := range tests {
		options := &str.Options{Module: "lists", Action: tt.action, Type: tt.listType}
		assert.Equal(t, tt.want, GetOutputForModule(options))
	}
}

func TestGetOutputForModuleSearchExactTrending(t *testing.T) {
	assert.Equal(t, "export_search_exact_query_movie.json", GetOutputForModule(&str.Options{Module: "search", Action: "exact_query", SearchType: str.Slice{"movie"}}))
	assert.Equal(t, "export_search_trending_people.json", GetOutputForModule(&str.Options{Module: "search", Action: "trending", SearchType: str.Slice{"people"}}))
}

func TestGetOutputForModuleSyncMinimalCollection(t *testing.T) {
	assert.Equal(t, "export_sync_minimal_collection_shows.json", GetOutputForModule(&str.Options{Module: "sync", Action: "get_minimal_collection", Type: "shows"}))
}

func TestGetOutputForModuleSyncProgress(t *testing.T) {
	assert.Equal(t, "export_sync_up_next.json", GetOutputForModule(&str.Options{Module: "sync", Action: "get_up_next"}))
	assert.Equal(t, "export_sync_up_next_nitro.json", GetOutputForModule(&str.Options{Module: "sync", Action: "get_up_next_nitro"}))
	assert.Equal(t, "export_sync_watched_progress.json", GetOutputForModule(&str.Options{Module: "sync", Action: "get_watched_progress"}))
}

func TestGetOutputForModuleSocialRecommendations(t *testing.T) {
	assert.Equal(t, "export_social_recommendations_movies.json", GetOutputForModule(&str.Options{Module: "social_recommendations", Action: "movies"}))
	assert.Equal(t, "export_social_recommendations_shows.json", GetOutputForModule(&str.Options{Module: "social_recommendations", Action: "shows"}))
}

func TestGetOutputForModuleUsersPlex(t *testing.T) {
	assert.Equal(t, "export_users_plex_settings.json", GetOutputForModule(&str.Options{Module: "users", Action: "plex_settings"}))
	assert.Equal(t, "export_users_plex_servers.json", GetOutputForModule(&str.Options{Module: "users", Action: "plex_servers"}))
	assert.Equal(t, "export_users_plex_server_abc.json", GetOutputForModule(&str.Options{Module: "users", Action: "plex_server", ID: "abc"}))
	assert.Equal(t, "export_users_plex_connect_results.json", GetOutputForModule(&str.Options{Module: "users", Action: "plex_connect"}))
}

func TestGetOutputForModuleUsersDataSyncs(t *testing.T) {
	assert.Equal(t, "export_users_data_syncs.json", GetOutputForModule(&str.Options{Module: "users", Action: "data_syncs"}))
	assert.Equal(t, "export_users_data_syncs_plex.json", GetOutputForModule(&str.Options{Module: "users", Action: "data_syncs", Type: "plex"}))
	assert.Equal(t, "export_users_data_sync_157.json", GetOutputForModule(&str.Options{Module: "users", Action: "data_sync", ID: "157"}))
	assert.Equal(t, "export_users_data_sync_skipped_157.json", GetOutputForModule(&str.Options{Module: "users", Action: "data_sync_skipped", ID: "157"}))
}

func TestGetOutputForModuleUsersAddSavedFilters(t *testing.T) {
	assert.Equal(t, "users_add_saved_filters_results.json", GetOutputForModule(&str.Options{Module: "users", Action: "add_saved_filters"}))
}

func TestGetOutputForModuleUsersSmartLists(t *testing.T) {
	assert.Equal(t, "export_users_smart_lists.json", GetOutputForModule(&str.Options{Module: "users", Action: "smart_lists"}))
	assert.Equal(t, "export_users_smart_list_sci-fi.json", GetOutputForModule(&str.Options{Module: "users", Action: "smart_list", ID: "sci-fi"}))
	assert.Equal(t, "users_add_smart_list_results.json", GetOutputForModule(&str.Options{Module: "users", Action: "add_smart_list"}))
	assert.Equal(t, "export_users_update_smart_list_results.json", GetOutputForModule(&str.Options{Module: "users", Action: "update_smart_list"}))
}

func TestGetOutputForModuleYounify(t *testing.T) {
	assert.Equal(t, "export_younify_connections.json", GetOutputForModule(&str.Options{Module: "younify", Action: "connections"}))
	assert.Equal(t, "export_younify_connect_results.json", GetOutputForModule(&str.Options{Module: "younify", Action: "connect"}))
}

func TestGetOutputForModuleSmartLists(t *testing.T) {
	assert.Equal(t, "export_smart_lists_summary_top-sci-fi.json", GetOutputForModule(&str.Options{Module: "smart_lists", Action: "summary", InternalID: "top-sci-fi"}))
	assert.Equal(t, "export_smart_lists_items_top-sci-fi.json", GetOutputForModule(&str.Options{Module: "smart_lists", Action: "items", InternalID: "top-sci-fi"}))
}

func TestGetOutputForModuleWatchNowSources(t *testing.T) {
	assert.Equal(t, "export_watchnow_sources.json", GetOutputForModule(&str.Options{Module: "watchnow", Action: "sources"}))
	assert.Equal(t, "export_watchnow_sources_us.json", GetOutputForModule(&str.Options{Module: "watchnow", Action: "sources", Country: "us"}))
}

func TestGetOutputForModuleTeamMembers(t *testing.T) {
	assert.Equal(t, "export_team_members.json", GetOutputForModule(&str.Options{Module: "team", Action: "members"}))
}

func TestGetOutputForModuleMoviesHotStreaming(t *testing.T) {
	assert.Equal(t, "export_movies_hot.json", GetOutputForModule(&str.Options{Module: "movies", Action: "hot"}))
	assert.Equal(t, "export_movies_streaming_weekly.json", GetOutputForModule(&str.Options{Module: "movies", Action: "streaming", Period: DefaultConfig().MoviesPeriod}))
}

func TestGetOutputForModuleMoviesSentiments(t *testing.T) {
	assert.Equal(t, "export_movies_sentiments_tron-legacy-2010.json", GetOutputForModule(&str.Options{Module: "movies", Action: "sentiments", InternalID: "tron-legacy-2010"}))
}

func TestGetOutputForModuleMoviesWatchNow(t *testing.T) {
	assert.Equal(t, "export_movies_watchnow_tron-legacy-2010.json", GetOutputForModule(&str.Options{Module: "movies", Action: "watchnow", InternalID: "tron-legacy-2010"}))
	assert.Equal(t, "export_movies_justwatch_links_tron-legacy-2010.json", GetOutputForModule(&str.Options{Module: "movies", Action: "justwatch_links", InternalID: "tron-legacy-2010"}))
}

func TestGetOutputForModuleShowsSentiments(t *testing.T) {
	assert.Equal(t, "export_shows_sentiments_the-sopranos.json", GetOutputForModule(&str.Options{Module: "shows", Action: "sentiments", InternalID: "the-sopranos"}))
}

func TestGetOutputForModuleShowsWatchNow(t *testing.T) {
	assert.Equal(t, "export_shows_watchnow_the-sopranos.json", GetOutputForModule(&str.Options{Module: "shows", Action: "watchnow", InternalID: "the-sopranos"}))
	assert.Equal(t, "export_shows_justwatch_links_the-sopranos.json", GetOutputForModule(&str.Options{Module: "shows", Action: "justwatch_links", InternalID: "the-sopranos"}))
}

func TestGetOutputForModuleSeasonsEpisodesWatchNow(t *testing.T) {
	assert.Equal(t, "export_episodes_watchnow_the-sopranos.json", GetOutputForModule(&str.Options{Module: "episodes", Action: "watchnow", InternalID: "the-sopranos"}))
	assert.Equal(t, "export_seasons_justwatch_links_the-sopranos.json", GetOutputForModule(&str.Options{Module: "seasons", Action: "justwatch_links", InternalID: "the-sopranos"}))
}

// TestGetOutputForModuleWriteResults keeps the result file names the sync/users write handlers used before -o applied to them.
func TestGetOutputForModuleWriteResults(t *testing.T) {
	tests := []struct {
		module, action, want string
	}{
		{module: "sync", action: "add_to_history", want: "sync_add_to_history_results.json"},
		{module: "sync", action: "remove_from_history", want: "sync_remove_from_history_results.json"},
		{module: "sync", action: "add_to_ratings", want: "sync_add_to_ratings_results.json"},
		{module: "sync", action: "remove_from_ratings", want: "sync_remove_from_ratings_results.json"},
		{module: "sync", action: "add_to_watchlist", want: "sync_add_to_watchlist_results.json"},
		{module: "sync", action: "remove_from_watchlist", want: "sync_remove_from_watchlist_results.json"},
		{module: "sync", action: "reorder_watchlist", want: "sync_reorder_watchlist_results.json"},
		{module: "sync", action: "add_to_favorites", want: "sync_add_to_favorites_results.json"},
		{module: "sync", action: "remove_from_favorites", want: "sync_remove_from_favorites_results.json"},
		{module: "sync", action: "reorder_favorites", want: "sync_reorder_favorites_results.json"},
		{module: "users", action: "add_list", want: "users_add_list_results.json"},
		{module: "users", action: "reorder_lists", want: "users_reorder_lists_results.json"},
		{module: "users", action: "add_list_items", want: "users_add_list_items_results.json"},
		{module: "users", action: "remove_list_items", want: "users_remove_list_items_results.json"},
		{module: "users", action: "reorder_list_items", want: "users_reorder_list_items_results.json"},
		{module: "users", action: "add_hidden_items", want: "users_add_hidden_items_results.json"},
		{module: "users", action: "remove_hidden_items", want: "users_remove_hidden_items_results.json"},
		{module: "users", action: "update_list", want: "export_users_update_list_results.json"},
		{module: "users", action: "watching", want: "export_users_watching_results.json"},
	}
	for _, tt := range tests {
		t.Run(tt.module+" "+tt.action, func(t *testing.T) {
			assert.Equal(t, tt.want, GetOutputForModule(&str.Options{Module: tt.module, Action: tt.action, Type: "movies"}))
		})
	}
}
