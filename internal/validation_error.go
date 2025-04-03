// Package internal used for client and services
package internal

import (
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// ValidationError occurs when trakt.tv returns 422 error
type ValidationError struct {
	Response *http.Response
	Message  string      `json:"message"`
	Errors   *str.Errors `json:"errors,omitempty"`
}

func (r *ValidationError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}
