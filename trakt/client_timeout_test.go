package trakt

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewClientTimeout(t *testing.T) {
	const custom = 5 * time.Second
	tests := []struct {
		name       string
		httpClient *http.Client
		want       time.Duration
	}{
		{name: "nil gets the default timeout", want: DefaultTimeout},
		{name: "a client without a timeout keeps none", httpClient: &http.Client{}, want: 0},
		{name: "a client keeps its own timeout", httpClient: &http.Client{Timeout: custom}, want: custom},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.httpClient)
			if got := c.client.Timeout; got != tt.want {
				t.Errorf("timeout is %v, want %v", got, tt.want)
			}
			if got := c.WithClientID("id").client.Timeout; got != tt.want {
				t.Errorf("timeout of a copy is %v, want %v", got, tt.want)
			}
		})
	}
}

// A response that never comes ends with an error once the timeout passes; the token in the URL is not in the message.
func TestClientStalledResponseTimesOut(t *testing.T) {
	setup := Setup()
	defer setup.Teardown()

	release := make(chan struct{})
	defer close(release)
	setup.Mux.HandleFunc("/movies/trending", func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	})
	setup.Client.client.Timeout = 50 * time.Millisecond

	done := make(chan error, 1)
	go func() {
		_, _, err := setup.Client.Movies.GetTrendingMovies(context.Background(), nil)
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "Client.Timeout") {
			t.Fatalf("error is %v, want a client timeout", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("request did not time out")
	}
}
