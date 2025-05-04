// Package internal used for client and services
package internal

import (
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/uri"
)

// ServerError occurs when trakt.tv returns 500 error
type ServerError struct {
	Response *http.Response
	Message  string `json:"message"`
}

func (r *ServerError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}
