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

func TestGetOutputForModuleMoviesHotStreaming(t *testing.T) {
	assert.Equal(t, "export_movies_hot.json", GetOutputForModule(&str.Options{Module: "movies", Action: "hot"}))
	assert.Equal(t, "export_movies_streaming_weekly.json", GetOutputForModule(&str.Options{Module: "movies", Action: "streaming", Period: DefaultConfig().MoviesPeriod}))
}

func TestGetOutputForModuleMoviesSentiments(t *testing.T) {
	assert.Equal(t, "export_movies_sentiments_tron-legacy-2010.json", GetOutputForModule(&str.Options{Module: "movies", Action: "sentiments", InternalID: "tron-legacy-2010"}))
}
