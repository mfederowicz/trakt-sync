package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"github.com/mfederowicz/trakt-sync/str"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	baseURLPath = "/trakt"
)

// setup sets up a test HTTP server along with a trakt.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the Trakt client being tested and is
	// configured to use test server.
	client = NewClient(nil)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := &str.NewDeviceCode{ClientID: str.String("abc")}, `{"client_id":"abc"}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := io.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestBareDo_returnsOpenBody(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	expectedBody := "Hello from the other side !"

	mux.HandleFunc("/test-url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, expectedBody)
	})

	ctx := context.Background()
	req, err := client.NewRequest("GET", "test-url", nil)
	if err != nil {
		t.Fatalf("client.NewRequest returned error: %v", err)
	}

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		t.Fatalf("client.BareDo returned error: %v", err)
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll returned error: %v", err)
	}
	if string(got) != expectedBody {
		t.Fatalf("Expected %q, got %q", expectedBody, string(got))
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatalf("resp.Body.Close() returned error: %v", err)
	}
}

func TestBareDo_rate_limit_reset(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	expectedBody := "Hello from the other side !"

	mux.HandleFunc("/test-url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Add(HeaderRetryAfter, "100")
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprint(w, expectedBody)
	})

	ctx := context.Background()
	req, err := client.NewRequest("GET", "test-url", nil)
	if err != nil {
		t.Fatalf("client.NewRequest returned error: %v", err)
	}

	client.BareDo(ctx, req)

	reset := client.RateLimitReset

	if reset.IsZero() {
		t.Fatalf("client.RateLimitReset is zero")
	}

	mux.HandleFunc("/test-url-next", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, "Body")
	})

	reqNext, errNext := client.NewRequest("GET", "test-url-next", nil)
	if errNext != nil {
		t.Fatalf("client.NewRequest returned error: %v", err)
	}

	_, errBare := client.BareDo(ctx, reqNext)
	if errBare != nil {
		// Update rate limit reset.
		rerr, _ := errBare.(*AbuseRateLimitError)
		if !strings.Contains(rerr.Message, "API rate limit exceeded until") {
			t.Fatal("Rate Limit Error msg not valid")
		}
	}
}
