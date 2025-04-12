// Package internal used for client and services
package internal

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/mfederowicz/trakt-sync/consts"
)

// Setup sets up a test HTTP server along with a trakt.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func Setup() *SetupData {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(consts.BaseURLPath+"/", http.StripPrefix(consts.BaseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the Trakt client being tested and is
	// configured to use test server.
	client := NewClient(nil)
	uri, _ := url.Parse(server.URL + consts.BaseURLPath + "/")
	client.BaseURL = uri

	return &SetupData{
		Client:    client,
		Mux:       mux,
		ServerURL: server.URL,
		Teardown:  server.Close,
	}
}

type values map[string]string

