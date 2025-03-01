// Package str used for structs
package str

import (
	"net/http"
)

// ErrorResponse represents reponse with message
type ErrorResponse struct {
	Response *http.Response `json:"-"`                 // HTTP response that caused this error
	Message  string         `json:"message,omitempty"` // error message
	Errors   *Errors        `json:"errors,omitempty"`  // errors object

}

func (r ErrorResponse) String() string {
	return Stringify(r)
}

func (r *ErrorResponse) Error() string {
	return Stringify(r)
}
