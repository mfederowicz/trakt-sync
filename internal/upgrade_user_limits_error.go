// Package internal used for client and services
package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/uri"
)

// UpgradeUserLimitsError occurs when trakt.tv returns 420 requests header
type UpgradeUserLimitsError struct {
	Response   *http.Response
	RetryAfter *time.Duration
	Message    string `json:"message"`
}

func (r *UpgradeUserLimitsError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}
