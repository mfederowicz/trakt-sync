package str

import (
	"fmt"
	"net/http"
	"github.com/mfederowicz/trakt-sync/uri"
)

type ErrorResponse struct {
	Response *http.Response `json:"-"`       // HTTP response that caused this error
	Message  string         `json:"message"` // error message
}

func (r *ErrorResponse) Error() string {
	if r.Response != nil && r.Response.Request != nil {
		return fmt.Sprintf("%v %v: %d %v",
			r.Response.Request.Method, uri.SanitizeURL(r.Response.Request.URL),
			r.Response.StatusCode, r.Message)
	}

	if r.Response != nil {
		return fmt.Sprintf("%d %v", r.Response.StatusCode, r.Message)
	}

	return fmt.Sprintf("%v", r.Message)
}
