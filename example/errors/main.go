// Package main shows how the client reports API errors. Every non-2xx response comes back as an error;
// the common statuses have their own types, which errors.As can pick out:
//
//	TRAKT_CLIENT_ID=... go run ./errors
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

func main() {
	client := trakt.NewClient(nil).
		WithClientID(os.Getenv("TRAKT_CLIENT_ID")).
		WithUserAgent("trakt-sync-example/1.0")
	// print every request; the client itself never writes to stdout
	client.DebugLogger = func(v ...any) { log.Println(v...) }

	slug := "no-such-movie-slug-12345"
	_, _, err := client.Movies.GetMovie(context.Background(), &slug, &uri.ListOptions{})
	fmt.Println("result:", explain(err))
}

// explain maps an error from the client to a message for the user.
func explain(err error) string {
	var notFound *trakt.NotFoundError
	var rateLimit *trakt.AbuseRateLimitError
	var upgrade *trakt.UpgradeRequiredError
	var invalidUser *trakt.InvalidUserError
	var other *str.ErrorResponse

	switch {
	case err == nil:
		return "ok"
	case errors.As(err, &notFound):
		return "not found: " + notFound.Error()
	case errors.As(err, &rateLimit):
		// later requests are refused locally until the reset time; wait before retrying
		if rateLimit.RetryAfter != nil {
			return fmt.Sprintf("rate limited, retry after %s", *rateLimit.RetryAfter)
		}
		return "rate limited"
	case errors.As(err, &upgrade):
		if upgrade.UpgradeURL != nil {
			return "Trakt VIP required, see " + upgrade.UpgradeURL.String()
		}
		return "Trakt VIP required"
	case errors.As(err, &invalidUser):
		return "missing or invalid access token (401)"
	case errors.As(err, &other):
		// statuses without their own type, e.g. 502 or 503
		return "API error: " + other.Error()
	default:
		// network errors, canceled contexts, ...
		return "request failed: " + err.Error()
	}
}
