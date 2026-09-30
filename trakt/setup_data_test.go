// Package trakt is a Go client for the Trakt API
package trakt

import "net/http"

// SetupData comment
type SetupData struct {
	Client    *Client
	Mux       *http.ServeMux
	ServerURL string
	Teardown  func()
}
