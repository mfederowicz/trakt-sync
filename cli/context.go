// Package cli for basic cli functions
package cli

import (
	"context"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// ContextFromOptions builds the request context for a CLI run: it carries the user's timezone
// (options.Timezone) when it is a valid IANA name, so the client converts timestamps to it.
func ContextFromOptions(options *str.Options) context.Context {
	ctx := context.Background()

	if len(options.Timezone) > consts.ZeroValue {
		loc, err := time.LoadLocation(options.Timezone)
		if err == nil {
			ctx = trakt.WithTimezone(ctx, loc)
		}
	}
	return ctx
}
