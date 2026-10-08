// Package cmds used for commands modules
package cmds

import (
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeCalendarsAction(t *testing.T) {
	for _, current := range calendarsActions {
		current := current
		legacy := strings.ReplaceAll(current, "_", "-")
		t.Run(legacy, func(t *testing.T) {
			assert.Equal(t, current, normalizeCalendarsAction(legacy))
			assert.Equal(t, current, normalizeCalendarsAction(current))
		})
	}
	for _, action := range []string{"bogus", "my-bogus", consts.EmptyString, cfg.DefaultConfig().Action} {
		assert.Equal(t, action, normalizeCalendarsAction(action), "%q must stay as it is", action)
	}
}

// TestCalendarsActionNames checks the list has all 20 actions and none still uses hyphens.
func TestCalendarsActionNames(t *testing.T) {
	assert.Len(t, calendarsActions, 20)
	for _, action := range calendarsActions {
		assert.NotContains(t, action, "-", "%s should use underscores", action)
	}
}

// TestCalendarsOutputForLegacyName checks an old action name still gets its usual export file name.
func TestCalendarsOutputForLegacyName(t *testing.T) {
	options := &str.Options{Module: consts.Calendars, Action: normalizeCalendarsAction("my-new-shows"), StartDate: "2026-10-01", Days: 7}
	assert.Equal(t, "export_calendars_new_shows_20261001_7.json", cfg.GetOutputForModule(options))
}
