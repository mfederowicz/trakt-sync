// Package internal used for client and services
package internal

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/uri"
)

// UpgradeRequiredError occurs when trakt.tv returns 426 user must upgrade to vip header
type UpgradeRequiredError struct {
	Response   *http.Response
	UpgradeURL *url.URL
	Message    string `json:"message"`
}

func (r *UpgradeRequiredError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}
