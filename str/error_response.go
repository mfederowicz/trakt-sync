// Package str used for structs
package str

import (
	"net/http"
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

func (r *ErrorResponse) Error() string {
	return Stringify(r)
}
