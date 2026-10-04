// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/stretchr/testify/assert"
)

// users -a lists sends sort_by and sort_how only for the flags given; other actions keep the global defaults.
func TestUsersListsSort(t *testing.T) {
	tests := []struct {
		name       string
		action     string
		sortBySet  bool
		sortHowSet bool
		wantBy     string
		wantHow    string
	}{
		{name: "lists drops default sort", action: consts.Lists, wantBy: "", wantHow: ""},
		{name: "lists keeps set sort", action: consts.Lists, sortBySet: true, sortHowSet: true, wantBy: "rank", wantHow: "asc"},
		{name: "lists keeps only sort_how", action: consts.Lists, sortHowSet: true, wantBy: "", wantHow: "asc"},
		{name: "other action untouched", action: consts.ListItems, wantBy: "rank", wantHow: "asc"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			options := &str.Options{Action: tt.action, SortBy: cfg.DefaultConfig().SortBy, SortHow: cfg.DefaultConfig().SortHow}
			usersListsSort(options, tt.sortBySet, tt.sortHowSet)
			assert.Equal(t, tt.wantBy, options.SortBy)
			assert.Equal(t, tt.wantHow, options.SortHow)
		})
	}
}
