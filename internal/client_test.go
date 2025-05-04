package internal

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

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

func TestHavePages(t *testing.T) {
	t.Helper()
	testSetup := Setup()
	client := testSetup.Client
	resp := &str.Response{Response: &http.Response{Header: http.Header{}}}
	resp.Header.Set(HeaderPaginationPage, strconv.Itoa(consts.FirstPage))
	resp.Header.Set(HeaderPaginationPageCount, strconv.Itoa(consts.AllPages))
	want := true
	if got := client.HavePages(consts.FirstPage, resp, consts.PagesLimit); got != want {
		t.Errorf(consts.HavePagesErrorStr, got, want)
	}
}

func TestHavePagesNoHeaders(t *testing.T) {
	t.Helper()
	testSetup := Setup()
	client := testSetup.Client
	resp := &str.Response{Response: &http.Response{Header: http.Header{}}}
	want := false
	if got := client.HavePages(consts.FirstPage, resp, consts.PagesLimit); got != want {
		t.Errorf(consts.HavePagesErrorStr, got, want)
	}
}

func TestHavePagesNoLimit(t *testing.T) {
	t.Helper()
	testSetup := Setup()
	client := testSetup.Client
	resp := &str.Response{Response: &http.Response{Header: http.Header{}}}
	resp.Header.Set(HeaderPaginationPage, strconv.Itoa(consts.FirstPage))
	resp.Header.Set(HeaderPaginationPageCount, strconv.Itoa(consts.AllPages))
	want := true
	if got := client.HavePages(consts.FirstPage, resp, consts.PagesNoLimit); got != want {
		t.Errorf(consts.HavePagesErrorStr, got, want)
	}
}

func TestHavePagesWithNoNext(t *testing.T) {
	t.Helper()
	testSetup := Setup()
	client := testSetup.Client
	resp := &str.Response{Response: &http.Response{Header: http.Header{}}}
	resp.Header.Set(HeaderPaginationPage, strconv.Itoa(consts.AllPages))
	resp.Header.Set(HeaderPaginationPageCount, strconv.Itoa(consts.AllPages))
	want := false
	if got := client.HavePages(consts.AllPages, resp, consts.PagesNoLimit); got != want {
		t.Errorf(consts.HavePagesErrorStr, got, want)
	}
}

func TestBareDo_returnsOpenBody(t *testing.T) {
	testSetup := Setup()
	client := testSetup.Client
	mux := testSetup.Mux
	teardown := testSetup.Teardown

	defer teardown()

	expectedBody := "Hello from the other side !"

	mux.HandleFunc("/"+consts.TestURL, func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		printer.Fprint(w, expectedBody)
	})

	ctx := context.Background()
	req, err := client.NewRequest(http.MethodGet, consts.TestURL, nil)
	if err != nil {
		t.Fatalf(consts.ClientNewRequestFatal, err)
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
	testSetup := Setup()
	client := testSetup.Client
	mux := testSetup.Mux
	teardown := testSetup.Teardown

	defer teardown()

	expectedBody := "Hello from the other side !"

	mux.HandleFunc("/test-url", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		w.Header().Add(HeaderRetryAfter, "100")
		w.WriteHeader(http.StatusTooManyRequests)
		printer.Fprint(w, expectedBody)
	})

	ctx := context.Background()
	req, err := client.NewRequest(http.MethodGet, "test-url", nil)
	if err != nil {
		t.Fatalf(consts.ClientNewRequestFatal, err)
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
		test.AssertMethod(t, r, http.MethodGet)
		printer.Fprint(w, "Body")
	})

	reqNext, errNext := client.NewRequest(http.MethodGet, consts.TestURLNext, nil)
	if errNext != nil {
		t.Fatalf(consts.ClientNewRequestFatal, err)
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
	testSetup := Setup()
	client := testSetup.Client
	mux := testSetup.Mux
	teardown := testSetup.Teardown

	defer teardown()

	expectedBody := "Hello vip!"

	mux.HandleFunc("/"+consts.TestURL, func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		w.Header().Add(HeaderUpgradeURL, upgradeURL)
		w.WriteHeader(http.StatusUpgradeRequired)
		printer.Fprint(w, expectedBody)
	})

	ctx := context.Background()

	var emptyURL *url.URL
	assert.Equal(t, client.UpgradeURL, emptyURL)

	req, err := client.NewRequest(http.MethodGet, consts.TestURL, nil)
	if err != nil {
		t.Fatalf(consts.ClientNewRequestFatal, err)
	}

	resp, err := client.BareDo(ctx, req)
	if err != nil {
		t.Fatalf("client.BareDo returned error: %s", err)
	}
	assert.Equal(t, resp.StatusCode, http.StatusUpgradeRequired)
	assert.Equal(t, client.UpgradeURL.String(), "https://trakt.tv/vip")
}
