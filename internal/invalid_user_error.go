// Package internal used for client and services
package internal

import (
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/uri"
)

// InvalidUserError occurs when trakt.tv returns 401 error
type InvalidUserError struct {
	Response  *http.Response
	Message   string `json:"message"`
	ErrorCode string `json:"error_code,omitempty"` // set by Plex routes, e.g. bad_auth
	Guidance  string `json:"guidance,omitempty"`
}

func (r *InvalidUserError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}
