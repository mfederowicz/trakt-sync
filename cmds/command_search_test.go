// Package cmds used for commands modules
package cmds

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeSearchAction(t *testing.T) {
	tests := []struct {
		name   string
		action string
		want   string
	}{
		{name: "legacy text-query", action: "text-query", want: "text_query"},
		{name: "legacy id-lookup", action: "id-lookup", want: "id_lookup"},
		{name: "text_query unchanged", action: consts.TextQuery, want: consts.TextQuery},
		{name: "id_lookup unchanged", action: consts.IDLookup, want: consts.IDLookup},
		{name: "unknown unchanged", action: "bogus", want: "bogus"},
		{name: "empty unchanged", action: consts.EmptyString, want: consts.EmptyString},
		{name: "default config action unchanged", action: cfg.DefaultConfig().Action, want: cfg.DefaultConfig().Action},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, normalizeSearchAction(tt.action))
		})
	}
}
