// Package handlers used to handle module actions
package handlers

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/stretchr/testify/assert"
)

// The comment routes of watchlist, favorites and lists accept only the sorts of
// listCommentsSortParamsSchema; live, every other value answers 400 invalid_sort.
func TestUsersCommentsSort(t *testing.T) {
	c := &CommonLogic{}
	defaults := cfg.DefaultConfig()
	actions := []string{"watchlist_comments", "favorites_comments", "list_comments"}

	for _, action := range actions {
		for _, sort := range []string{defaults.UsersSort, "likes", "replies", "newest", "oldest"} {
			options := &str.Options{Module: "users", Action: action, Type: defaults.UsersType, Sort: sort}
			assert.NoError(t, c.CheckSortAndTypes(options), "%s -s %s", action, sort)
		}
		for _, sort := range []string{"likes_30", "replies_30", "plays", "rating", "added"} {
			options := &str.Options{Module: "users", Action: action, Type: defaults.UsersType, Sort: sort}
			assert.EqualError(t, c.CheckSortAndTypes(options), "not found sort for module 'users'", "%s -s %s", action, sort)
		}
	}
}
