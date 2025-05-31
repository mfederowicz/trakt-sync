// Package internal used for client and services
package internal

import (
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/uri"
)

// PreconditionFailedRequestError occurs when trakt.tv returns 412 error
type PreconditionFailedRequestError struct {
	Response *http.Response
	Message  string `json:"message"`
}

func (r *PreconditionFailedRequestError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}
