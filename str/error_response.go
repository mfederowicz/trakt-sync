package str

import (
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/uri"
)

// ErrorResponse represents reponse with message
type ErrorResponse struct {
	Response  *http.Response `json:"-"`                    // HTTP response that caused this error
	Message   string         `json:"message,omitempty"`    // error message
	Errors    *Errors        `json:"errors,omitempty"`     // errors object
	ErrorCode string         `json:"error_code,omitempty"` // machine-readable code (Plex routes)
	Guidance  string         `json:"guidance,omitempty"`   // what to do about the error (Plex routes)

}

func (r ErrorResponse) String() string {
	return Stringify(r)
}

// Error formats the error like the typed client errors: method, URL, status and API message.
func (r *ErrorResponse) Error() string {
	if r.Response == nil || r.Response.Request == nil || r.Response.Request.URL == nil {
		return Stringify(r)
	}
	u := *r.Response.Request.URL
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(&u),
		r.Response.StatusCode,
		r.Message,
	)
}
