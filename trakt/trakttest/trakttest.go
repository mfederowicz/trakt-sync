// Package trakttest provides a mock Trakt API server and a client wired to it, for tests outside the trakt package.
package trakttest

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// SetupData holds the mock server and the client that talks to it.
type SetupData struct {
	Client    *trakt.Client
	Mux       *http.ServeMux
	ServerURL string
	Teardown  func()
}

// Setup sets up a test HTTP server along with a trakt.Client that is
// configured to talk to that test server. Tests should register handlers on
// Mux which provide mock responses for the API method being tested.
func Setup() *SetupData {
	mux := http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(consts.BaseURLPath+"/", http.StripPrefix(consts.BaseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)

	client := trakt.NewClient(nil)
	uri, _ := url.Parse(server.URL + consts.BaseURLPath + "/")
	client.BaseURL = uri

	return &SetupData{
		Client:    client,
		Mux:       mux,
		ServerURL: server.URL,
		Teardown:  server.Close,
	}
}
