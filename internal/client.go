// Package internal used for client and services
package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// basic consts for client
const (
	Version                   = "2"
	defaultBaseURL            = "https://api.trakt.tv/"
	HeaderRateLimit           = "X-RateLimit"
	HeaderRetryAfter          = "Retry-After"
	HeaderPaginationPageCount = "X-Pagination-Page-Count"

	skipRateLimitCheck requestContext = iota
	emptyLimit                        = ""
)

var errNonNilContext = errors.New("context must be non-nil")
var emptyReader = strings.NewReader("")

type requestContext uint8

// AbuseRateLimitError occurs when trakt.tv returns 429 too many requests header
type AbuseRateLimitError struct {
	Response   *http.Response
	RetryAfter *time.Duration
	Message    string `json:"message"`
}

func (r *AbuseRateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}

// RequestOption represents an option that can modify an http.Request.
type RequestOption func(req *http.Request)

// A Client manages communication with the trakt.tv API.
type Client struct {
	RateLimitReset time.Time
	client         *http.Client
	BaseURL        *url.URL
	headers        map[string]any
	common         Service
	Oauth          *OauthService
	Users          *UsersService
	Sync           *SyncService
	People         *PeopleService
	Calendars      *CalendarsService
	Search         *SearchService
	rateMu         sync.Mutex
}

// UpdateHeaders is for update client headers map
func (c *Client) UpdateHeaders(headers map[string]any) {
	c.headers = headers
}

// HavePages checks if we have available pages to fetch
func (c *Client) HavePages(page int, pages int) bool {
	return page != pages && pages > consts.ZeroValue
}

// initialize sets default values and initializes services.
func (c *Client) initialize() {
	if c.client == nil {
		c.client = &http.Client{}
	}
	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse(defaultBaseURL)
	}

	c.common.client = c
	c.Oauth = (*OauthService)(&c.common)
	c.Users = (*UsersService)(&c.common)
	c.Sync = (*SyncService)(&c.common)
	c.People = (*PeopleService)(&c.common)
	c.Calendars = (*CalendarsService)(&c.common)
	c.Search = (*SearchService)(&c.common)
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body any, opts ...RequestOption) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.headers["Authorization"] != nil {
		req.Header.Set("Authorization", c.headers["Authorization"].(string))
	}

	if c.headers["trakt-api-key"] != nil {
		req.Header.Set("trakt-api-key", c.headers["trakt-api-key"].(string))
	}

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred
func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*str.Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:

		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}

// BareDo sends an API request and lets you handle the api response.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*str.Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	req = c.WithContext(ctx, req)

	if skip := ctx.Value(skipRateLimitCheck); skip == nil {
		// don't make further requests before Retry After.
		if err := c.CheckRetryAfter(req); err != nil {
			return &str.Response{
				Response: err.Response,
			}, err
		}
	}

	resp, err := c.client.Do(req)

	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if u, err := url.Parse(e.URL); err == nil {
				e.URL = uri.SanitizeURL(u).String()
				return nil, e
			}
		}

		return nil, err
	}

	response := c.NewResponse(resp)

	errCheck := c.CheckResponse(resp)
	if errCheck != nil {
		defer resp.Body.Close()

		// Update rate limit reset.
		rerr, ok := errCheck.(*AbuseRateLimitError)
		if ok && rerr.RetryAfter != nil {
			c.rateMu.Lock()
			c.RateLimitReset = time.Now().Add(*rerr.RetryAfter)
			c.rateMu.Unlock()
		}
	}

	return response, nil
}

// CheckResponse checks if api response have errors.
func (c *Client) CheckResponse(r *http.Response) error {
	if c := r.StatusCode; http.StatusOK <= c && c <= consts.MaxAcceptedStatus {
		return nil
	}

	errorResponse := &str.ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			// reset the response as if this never happened
			errorResponse = &str.ErrorResponse{Response: r}
		}
	}

	r.Body = io.NopCloser(bytes.NewBuffer(data))
	switch r.StatusCode {
	case http.StatusTooManyRequests:
		abuseRateLimitError := &AbuseRateLimitError{
			Response: errorResponse.Response,
			Message:  errorResponse.Message,
		}
		if retryAfter := c.ParseRateLimit(r); retryAfter != nil {
			abuseRateLimitError.RetryAfter = retryAfter
			return abuseRateLimitError
		}
		return nil

	default:
		return errorResponse
	}
}

// CheckRetryAfter check Retry After header.
func (c *Client) CheckRetryAfter(req *http.Request) *AbuseRateLimitError {
	c.rateMu.Lock()
	reset := c.RateLimitReset
	c.rateMu.Unlock()
	if !reset.IsZero() && time.Now().Before(reset) {
		// Create a fake response.
		resp := &http.Response{
			Status:     http.StatusText(http.StatusForbidden),
			StatusCode: http.StatusForbidden,
			Request:    req,
			Header:     make(http.Header),
			Body:       io.NopCloser(emptyReader),
		}

		retryAfter := time.Until(reset)
		return &AbuseRateLimitError{
			Response:   resp,
			Message:    fmt.Sprintf("API rate limit exceeded until %v, not making remote request.", reset),
			RetryAfter: &retryAfter,
		}
	}

	return nil
}

// WithContext pass context to request
func (c *Client) WithContext(ctx context.Context, req *http.Request) *http.Request {
	return req.WithContext(ctx)
}

// ParseRate parses the rate related headers.
func (c *Client) ParseRate(r *http.Response) str.Rate {
	var rate str.Rate
	if limit := r.Header.Get(HeaderRateLimit); limit != emptyLimit {
		rate.Limit, _ = strconv.Atoi(limit)
	}

	return rate
}

// NewResponse creates a new Response for the provided http.Response.
// r must not be nil.
func (c *Client) NewResponse(r *http.Response) *str.Response {
	response := &str.Response{Response: r}
	response.Rate = c.ParseRate(r)
	return response
}

// ParseRateLimit parses related headers, and returns the time to retry after.
func (c *Client) ParseRateLimit(r *http.Response) *time.Duration {
	// number of seconds that one should
	// wait before resuming making requests.
	if v := r.Header.Get(HeaderRetryAfter); v != "" {
		retryAfterSeconds, _ := strconv.ParseInt(v, consts.BaseInt, consts.BitSize) // Error handling is noop.
		retryAfter := time.Duration(retryAfterSeconds) * time.Second
		return &retryAfter
	}

	return nil
}

// NewClient returns a new API client. If a nil httpClient is
// provided, a new http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient2 := *httpClient
	c := &Client{client: &httpClient2}
	c.initialize()
	return c
}
