// Package str used for structs
package str

import "net/http"

// Response http object
type Response struct {
	*http.Response
	// Explicitly specify the Rate type so Rate's String() receiver doesn't
	// propagate to Response.
	Rate   Rate
	Errors *Errors
}
