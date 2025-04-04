package internal

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/stretchr/testify/assert"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	baseURLPath           = "/trakt"
	clientNewRequestFatal = "client.NewRequest returned error: %v"
	firstPage             = 1
	allPages              = 10
	pagesLimit            = 2
	pagesNolimit          = 0
	havePagesErrorStr     = "HavePages: %v, want %v"
)

type TestSetup struct {
	Client    *Client
	Mux       *http.ServeMux
	ServerURL string
	Teardown  func()
}

// setup sets up a test HTTP server along with a trakt.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() *TestSetup {
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the Trakt client being tested and is
	// configured to use test server.
	client := NewClient(nil)
	uri, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = uri

	return &TestSetup{
		Client:    client,
		Mux:       mux,
		ServerURL: server.URL,
		Teardown:  server.Close,
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := &str.NewDeviceCode{ClientID: str.String("abc")}, `{"client_id":"abc"}`+"\n"
	req, _ := c.NewRequest(http.MethodGet, inURL, inBody)

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

func TestHavePages(t *testing.T) {
	t.Helper()
	testSetup := setup()
	client := testSetup.Client
	resp := &str.Response{Response: &http.Response{Header: http.Header{}}}
	resp.Header.Set(HeaderPaginationPage, strconv.Itoa(firstPage))
	resp.Header.Set(HeaderPaginationPageCount, strconv.Itoa(allPages))
	want := true
	if got := client.HavePages(firstPage, resp, pagesLimit); got != want {
		t.Errorf(havePagesErrorStr, got, want)
	}
}

func TestHavePagesNoHeaders(t *testing.T) {
	t.Helper()
	testSetup := setup()
	client := testSetup.Client
	resp := &str.Response{Response: &http.Response{Header: http.Header{}}}
	want := false
	if got := client.HavePages(firstPage, resp, pagesLimit); got != want {
		t.Errorf(havePagesErrorStr, got, want)
	}
}

func TestHavePagesNoLimit(t *testing.T) {
	t.Helper()
	testSetup := setup()
	client := testSetup.Client
	resp := &str.Response{Response: &http.Response{Header: http.Header{}}}
	resp.Header.Set(HeaderPaginationPage, strconv.Itoa(firstPage))
	resp.Header.Set(HeaderPaginationPageCount, strconv.Itoa(allPages))
	want := true
	if got := client.HavePages(firstPage, resp, pagesNolimit); got != want {
		t.Errorf(havePagesErrorStr, got, want)
	}
}

func TestHavePagesWithNoNext(t *testing.T) {
	t.Helper()
	testSetup := setup()
	client := testSetup.Client
	resp := &str.Response{Response: &http.Response{Header: http.Header{}}}
	resp.Header.Set(HeaderPaginationPage, strconv.Itoa(allPages))
	resp.Header.Set(HeaderPaginationPageCount, strconv.Itoa(allPages))
	want := false
	if got := client.HavePages(allPages, resp, pagesNolimit); got != want {
		t.Errorf(havePagesErrorStr, got, want)
	}
}

func TestBareDo_returnsOpenBody(t *testing.T) {
	testSetup := setup()
	client := testSetup.Client
	mux := testSetup.Mux
	teardown := testSetup.Teardown

	defer teardown()

	expectedBody := "Hello from the other side !"

	mux.HandleFunc("/"+consts.TestURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		printer.Fprint(w, expectedBody)
	})

	ctx := context.Background()
	req, err := client.NewRequest(http.MethodGet, consts.TestURL, nil)
	if err != nil {
		t.Fatalf(clientNewRequestFatal, err)
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
	testSetup := setup()
	client := testSetup.Client
	mux := testSetup.Mux
	teardown := testSetup.Teardown

	defer teardown()

	expectedBody := "Hello from the other side !"

	mux.HandleFunc("/test-url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Add(HeaderRetryAfter, "100")
		w.WriteHeader(http.StatusTooManyRequests)
		printer.Fprint(w, expectedBody)
	})

	ctx := context.Background()
	req, err := client.NewRequest(http.MethodGet, "test-url", nil)
	if err != nil {
		t.Fatalf(clientNewRequestFatal, err)
	}

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		t.Fatalf("client.BareDo returned error: %s", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusTooManyRequests)

	reset := client.RateLimitReset

	if reset.IsZero() {
		t.Fatalf("client.RateLimitReset is zero")
	}

	mux.HandleFunc("/"+consts.TestURLNext, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		printer.Fprint(w, "Body")
	})

	reqNext, errNext := client.NewRequest(http.MethodGet, consts.TestURLNext, nil)
	if errNext != nil {
		t.Fatalf(clientNewRequestFatal, err)
	}

	_, errBare := client.BareDo(ctx, reqNext)
	if errBare != nil {
		// Update rate limit reset.
		err, ok := errBare.(*AbuseRateLimitError)
		assert.Equal(t, ok, true)
		if !strings.Contains(err.Message, "API rate limit exceeded until") {
			t.Fatal("Rate Limit Error msg not valid")
		}
	}
}

func TestBareDo_upgrade_required(t *testing.T) {
	testSetup := setup()
	client := testSetup.Client
	mux := testSetup.Mux
	teardown := testSetup.Teardown

	defer teardown()

	expectedBody := "Hello vip!"

	mux.HandleFunc("/"+consts.TestURL, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Add(HeaderUpgradeURL, upgradeURL)
		w.WriteHeader(http.StatusUpgradeRequired)
		printer.Fprint(w, expectedBody)
	})

	ctx := context.Background()

	var emptyURL *url.URL
	assert.Equal(t, client.UpgradeURL, emptyURL)

	req, err := client.NewRequest(http.MethodGet, consts.TestURL, nil)
	if err != nil {
		t.Fatalf(clientNewRequestFatal, err)
	}

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		t.Fatalf("client.BareDo returned error: %s", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusUpgradeRequired)
	assert.Equal(t, client.UpgradeURL.String(), "https://trakt.tv/vip")
}
