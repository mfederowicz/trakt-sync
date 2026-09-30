// Package cli for basic cli functions
package cli

import (
	"context"
	"testing"
	"time"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/stretchr/testify/assert"
)

// TestContextFromOptions checks the context carries a valid options timezone and falls back to UTC otherwise.
func TestContextFromOptions(t *testing.T) {
	warsaw, err := time.LoadLocation("Europe/Warsaw")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		timezone string
		want     *time.Location
	}{
		{name: "valid timezone", timezone: "Europe/Warsaw", want: warsaw},
		{name: "empty timezone", timezone: "", want: time.UTC},
		{name: "invalid timezone", timezone: "Mars/Olympus", want: time.UTC},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := ContextFromOptions(&str.Options{Timezone: tt.timezone})
			assert.Equal(t, tt.want.String(), trakt.NewClient(nil).GetTimezone(ctx).String())
		})
	}
}

// TestContextFromOptionsIsBackground checks the context has no deadline or cancel, like before.
func TestContextFromOptionsIsBackground(t *testing.T) {
	ctx := ContextFromOptions(&str.Options{})
	_, hasDeadline := ctx.Deadline()
	assert.False(t, hasDeadline)
	assert.Nil(t, ctx.Done())
	assert.Equal(t, context.Background(), ctx)
}
