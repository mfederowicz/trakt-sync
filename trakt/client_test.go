package trakt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

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

// TestNewRequestHeaders checks every request carries the headers Trakt requires.
func TestNewRequestHeaders(t *testing.T) {
	tests := []struct {
		name      string
		authToken string
		clientID  string
		userAgent string
		want      map[string]string
	}{
		{
			name:      "all set",
			authToken: "token",
			clientID:  "client-id",
			userAgent: "trakt-sync/1.19.1",
			want: map[string]string{
				"Content-Type":      "application/json",
				"trakt-api-version": "2",
				"User-Agent":        "trakt-sync/1.19.1",
				"trakt-api-key":     "client-id",
				"Authorization":     "Bearer token",
			},
		},
		{
			name: "no headers",
			want: map[string]string{
				"Content-Type":      "application/json",
				"trakt-api-version": "2",
				"User-Agent":        DefaultUserAgent,
				"trakt-api-key":     "",
				"Authorization":     "",
			},
		},
		{
			name:     "empty authorization",
			clientID: "client-id",
			want: map[string]string{
				"User-Agent":    DefaultUserAgent,
				"trakt-api-key": "client-id",
				"Authorization": "",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(nil).WithAuthToken(tt.authToken).WithClientID(tt.clientID).WithUserAgent(tt.userAgent)
			req, err := c.NewRequest(http.MethodGet, "/foo", nil)
			assert.NoError(t, err)
			for name, want := range tt.want {
				assert.Equal(t, want, req.Header.Get(name), name)
			}
			if tt.want["Authorization"] == "" {
				_, present := req.Header["Authorization"]
				assert.False(t, present, "Authorization must not be sent empty")
			}
			if tt.want["trakt-api-key"] == "" {
				_, present := req.Header["Trakt-Api-Key"]
				assert.False(t, present, "trakt-api-key must not be sent empty")
			}
		})
	}
}

// TestWithMethodsReturnCopies checks With* leave the original client alone and the copy's services use the copy.
func TestWithMethodsReturnCopies(t *testing.T) {
	testSetup := Setup()
	defer testSetup.Teardown()

	testSetup.Mux.HandleFunc("/countries/movies", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer new-token", r.Header.Get("Authorization"))
		assert.Equal(t, "client-id", r.Header.Get("trakt-api-key"))
		test.SafeFprint(w, "[]")
	})

	original := testSetup.Client.WithClientID("client-id")
	original.DebugLogger = func(...any) {}
	authed := original.WithAuthToken("new-token")

	assert.NotSame(t, original, authed)
	assert.Equal(t, original.BaseURL, authed.BaseURL)
	assert.NotNil(t, authed.DebugLogger)

	req, err := original.NewRequest(http.MethodGet, "/foo", nil)
	assert.NoError(t, err)
	assert.Empty(t, req.Header.Get("Authorization"), "the original client must keep no token")

	_, _, err = authed.Countries.GetCountries(context.Background(), "movies")
	assert.NoError(t, err)
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

// TestWithTimezone checks GetTimezone returns the location set by WithTimezone, and UTC without one.
func TestWithTimezone(t *testing.T) {
	client := NewClient(nil)
	loc := time.FixedZone("UTC+2", 2*60*60)

	assert.Equal(t, loc, client.GetTimezone(WithTimezone(context.Background(), loc)))
	assert.Equal(t, time.UTC, client.GetTimezone(context.Background()))
}

// TestDoTimezoneKeepsMidnight checks a response time that is midnight in the user's timezone keeps its time,
// and a date without a time keeps its day.
func TestDoTimezoneKeepsMidnight(t *testing.T) {
	testSetup := Setup()
	client := testSetup.Client
	mux := testSetup.Mux
	teardown := testSetup.Teardown

	defer teardown()

	mux.HandleFunc("/"+consts.TestURL, func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		printer.Fprint(w, `[{"listed_at":"2026-09-30T22:00:00.000Z"},{"listed_at":"2026-10-01T04:00:00.000Z"},{"listed_at":"2026-10-01"}]`)
	})

	cases := []struct {
		name string
		loc  *time.Location
		want string
	}{
		{name: "UTC", loc: time.UTC, want: `[{"listed_at":"2026-09-30T22:00:00Z"},{"listed_at":"2026-10-01T04:00:00Z"},{"listed_at":"2026-10-01"}]`},
		{name: "east of UTC", loc: time.FixedZone("CEST", 2*60*60), want: `[{"listed_at":"2026-10-01T00:00:00+02:00"},{"listed_at":"2026-10-01T06:00:00+02:00"},{"listed_at":"2026-10-01"}]`},
		{name: "west of UTC", loc: time.FixedZone("EDT", -4*60*60), want: `[{"listed_at":"2026-09-30T18:00:00-04:00"},{"listed_at":"2026-10-01T00:00:00-04:00"},{"listed_at":"2026-10-01"}]`},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req, err := client.NewRequest(http.MethodGet, consts.TestURL, nil)
			assert.NoError(t, err)

			list := []*str.UserListItem{}
			_, err = client.Do(WithTimezone(context.Background(), tc.loc), req, &list)
			assert.NoError(t, err)

			got, err := json.Marshal(list)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, string(got))
		})
	}
}

// TestDebugLogger checks DebugLogger gets the service note and the request line, with client_secret redacted.
func TestDebugLogger(t *testing.T) {
	testSetup := Setup()
	client := testSetup.Client
	mux := testSetup.Mux
	teardown := testSetup.Teardown

	defer teardown()

	mux.HandleFunc("/countries/movies", func(w http.ResponseWriter, r *http.Request) {
		test.AssertMethod(t, r, http.MethodGet)
		printer.Fprint(w, "[]")
	})
	mux.HandleFunc("/"+consts.TestURL, func(_ http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "secret-value", r.URL.Query().Get("client_secret"))
	})

	logged := []string{}
	client.DebugLogger = func(v ...any) {
		logged = append(logged, fmt.Sprint(v...))
	}

	_, _, err := client.Countries.GetCountries(context.Background(), "movies")
	assert.NoError(t, err)

	req, err := client.NewRequest(http.MethodGet, consts.TestURL+"?client_secret=secret-value", nil)
	if err != nil {
		t.Fatalf(consts.ClientNewRequestFatal, err)
	}
	_, err = client.Do(context.Background(), req, nil)
	assert.NoError(t, err)

	assert.Equal(t, []string{
		"fetch countries url:countries/movies",
		http.MethodGet + " " + client.BaseURL.String() + "countries/movies",
		http.MethodGet + " " + client.BaseURL.String() + consts.TestURL + "?client_secret=REDACTED",
	}, logged)
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
	var rateErr *AbuseRateLimitError
	assert.ErrorAs(t, err, &rateErr)
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
	var upgradeErr *UpgradeRequiredError
	assert.ErrorAs(t, err, &upgradeErr)
	assert.Equal(t, resp.StatusCode, http.StatusUpgradeRequired)
	assert.Equal(t, client.UpgradeURL.String(), "https://trakt.tv/vip")
}

func TestDo_returnsTypedErrors(t *testing.T) {
	tests := []struct {
		name   string
		status int
		as     func(err error) bool
	}{
		{name: "400", status: http.StatusBadRequest, as: func(err error) bool { var e *BadRequestError; return errors.As(err, &e) }},
		{name: "401", status: http.StatusUnauthorized, as: func(err error) bool { var e *InvalidUserError; return errors.As(err, &e) }},
		{name: "403", status: http.StatusForbidden, as: func(err error) bool { var e *ForbiddenError; return errors.As(err, &e) }},
		{name: "404", status: http.StatusNotFound, as: func(err error) bool { var e *NotFoundError; return errors.As(err, &e) }},
		{name: "409", status: http.StatusConflict, as: func(err error) bool { var e *ConflictError; return errors.As(err, &e) }},
		{name: "420", status: 420, as: func(err error) bool { var e *UpgradeUserLimitsError; return errors.As(err, &e) }},
		{name: "500", status: http.StatusInternalServerError, as: func(err error) bool { var e *ServerError; return errors.As(err, &e) }},
		// these returned no error before (or panicked for 426/429 without their headers)
		{name: "405", status: http.StatusMethodNotAllowed, as: func(err error) bool { var e *str.ErrorResponse; return errors.As(err, &e) }},
		{name: "412", status: http.StatusPreconditionFailed, as: func(err error) bool { var e *PreconditionFailedRequestError; return errors.As(err, &e) }},
		{name: "426", status: http.StatusUpgradeRequired, as: func(err error) bool { var e *UpgradeRequiredError; return errors.As(err, &e) }},
		{name: "429", status: http.StatusTooManyRequests, as: func(err error) bool { var e *AbuseRateLimitError; return errors.As(err, &e) }},
		{name: "502", status: http.StatusBadGateway, as: func(err error) bool { var e *str.ErrorResponse; return errors.As(err, &e) }},
		{name: "503", status: http.StatusServiceUnavailable, as: func(err error) bool { var e *str.ErrorResponse; return errors.As(err, &e) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSetup := Setup()
			defer testSetup.Teardown()

			testSetup.Mux.HandleFunc("/"+consts.TestURL, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				test.SafeFprint(w, `{"message":"boom"}`)
			})

			req, err := testSetup.Client.NewRequest(http.MethodGet, consts.TestURL, nil)
			if err != nil {
				t.Fatalf(consts.ClientNewRequestFatal, err)
			}

			resp, err := testSetup.Client.Do(context.Background(), req, nil)
			if err == nil {
				t.Fatal("expected an error")
			}
			if !tt.as(err) {
				t.Errorf("error %T is not the typed error for status %d", err, tt.status)
			}
			if !strings.Contains(err.Error(), "boom") {
				t.Errorf("error %q does not carry the API message", err)
			}
			if resp == nil || resp.StatusCode != tt.status {
				t.Errorf("response status is not %d", tt.status)
			}
		})
	}
}

func TestDo_withoutResponseDoesNotPanic(t *testing.T) {
	testSetup := Setup()
	testSetup.Teardown() // closed server: the request fails before any response

	req, err := testSetup.Client.NewRequest(http.MethodGet, consts.TestURL, nil)
	if err != nil {
		t.Fatalf(consts.ClientNewRequestFatal, err)
	}

	result := new(str.Movie)
	resp, err := testSetup.Client.Do(context.Background(), req, result)
	if err == nil {
		t.Fatal("expected an error")
	}
	if resp != nil {
		t.Errorf("response is %v, want nil", resp)
	}
}
