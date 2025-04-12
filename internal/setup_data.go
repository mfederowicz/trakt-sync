// Package internal abc
package internal

import "net/http"

// SetupData comment
type SetupData struct {
	Client    *Client
	Mux       *http.ServeMux
	ServerURL string
	Teardown  func()
}

