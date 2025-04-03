// Package internal used for client and services
package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/uri"
)

// AbuseRateLimitError occurs when trakt.tv returns 429 too many requests header
type AbuseRateLimitError struct {
	Response   *http.Response
	RetryAfter *time.Duration
	Message    string `json:"message"`
}

func (r *AbuseRateLimitError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}

