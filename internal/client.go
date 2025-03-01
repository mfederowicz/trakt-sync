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
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// basic consts for client
const (
	Version                   = "2"
	defaultBaseURL            = "https://api.trakt.tv/"
	upgradeURL                = "https://trakt.tv/vip"
	HeaderRateLimit           = "X-RateLimit"
	HeaderRetryAfter          = "Retry-After"
	HeaderPaginationPage      = "X-Pagination-Page"
	HeaderPaginationPageCount = "X-Pagination-Page-Count"
	HeaderUpgradeURL          = "X-Upgrade-URL"

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
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}

// UpgradeRequiredError occurs when trakt.tv returns 426 user must upgrade to vip header
type UpgradeRequiredError struct {
	Response   *http.Response
	UpgradeURL *url.URL
	Message    string `json:"message"`
}

func (r *UpgradeRequiredError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}

// NotFoundError occurs when trakt.tv returns 404 error
type NotFoundError struct {
	Response *http.Response
	Message  string `json:"message"`
}

func (r *NotFoundError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}

// ConflictError occurs when trakt.tv returns 409 error
type ConflictError struct {
	Response *http.Response
	Message  string `json:"message"`
}

func (r *ConflictError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
		r.Response.Request.Method,
		uri.SanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode,
		r.Message,
	)
}

// ValidationError occurs when trakt.tv returns 422 error
type ValidationError struct {
	Response *http.Response
	Message  string      `json:"message"`
	Errors   *str.Errors `json:"errors,omitempty"`
}

func (r *ValidationError) Error() string {
	return fmt.Sprintf(consts.ErrorsPlaceholders,
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
	UpgradeURL     *url.URL
	headers        map[string]any
	common         Service
	Oauth          *OauthService
	Users          *UsersService
	Sync           *SyncService
	People         *PeopleService
	Calendars      *CalendarsService
	Certifications *CertificationsService
	Checkin        *CheckinService
	Comments       *CommentsService
	Search         *SearchService
	Lists          *ListsService
	Movies         *MoviesService
	Episodes       *EpisodesService
	Shows          *ShowsService
	Seasons        *SeasonsService
	rateMu         sync.Mutex
}

// UpdateHeaders is for update client headers map
func (c *Client) UpdateHeaders(headers map[string]any) {
	c.headers = headers
}

// HavePages checks if we have available pages to fetch
func (*Client) HavePages(page int, resp *str.Response) bool {
	_, pageHeader := resp.Header[HeaderPaginationPage]
	pages, _ := strconv.Atoi(resp.Header.Get(HeaderPaginationPageCount))
	return pageHeader && page != pages && pages > consts.ZeroValue
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
	c.Certifications = (*CertificationsService)(&c.common)
	c.Checkin = (*CheckinService)(&c.common)
	c.Comments = (*CommentsService)(&c.common)
	c.Search = (*SearchService)(&c.common)
	c.Lists = (*ListsService)(&c.common)
	c.Movies = (*MoviesService)(&c.common)
	c.Episodes = (*EpisodesService)(&c.common)
	c.Shows = (*ShowsService)(&c.common)
	c.Seasons = (*SeasonsService)(&c.common)
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

	req = c.requestSetHeaders(req, body)

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func (c *Client) requestSetHeaders(r *http.Request, body any) *http.Request {
	if body != nil {
		r.Header.Set("Content-Type", "application/json")
	}

	if c.headers["Authorization"] != nil {
		r.Header.Set("Authorization", c.headers["Authorization"].(string))
	}

	if c.headers["trakt-api-key"] != nil {
		r.Header.Set("trakt-api-key", c.headers["trakt-api-key"].(string))
	}

	return r
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

	skipResp, skipErr := c.skipCheck(ctx, req)
	if skipErr != nil {
		return skipResp, skipErr
	}

	resp, err := c.client.Do(req)

	if err != nil {
		return handleBareDoError(ctx, err)
	}

	response := c.NewResponse(resp)

	errCheck := c.CheckResponse(resp)
	if errCheck != nil {
		defer resp.Body.Close()
		switch e := errCheck.(type) {
		case *AbuseRateLimitError:
			updateRateLimitReset(c, e)
		case *UpgradeRequiredError:
			upgradeAccountRequired(c, e)
		case *NotFoundError:
			return response, errors.New(e.Error())
		case *ConflictError:
		case *ValidationError:
			response.Errors = e.Errors
			return response, errors.New("validation error")
		default:
			printer.Println("General error occurred:", errCheck)
		}
	}

	return response, nil
}

func upgradeAccountRequired(c *Client, errCheck error) {
	rerr, ok := errCheck.(*UpgradeRequiredError)
	if ok && rerr.UpgradeURL != nil {
		c.rateMu.Lock()
		c.UpgradeURL = rerr.UpgradeURL
		c.rateMu.Unlock()
	}
}

func handleBareDoError(ctx context.Context, err error) (*str.Response, error) {
	// If the context has been canceled, return its error.
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Try to sanitize the URL in the error and return the sanitized error if successful.
	if sanitizedErr := sanitizeURL(err); sanitizedErr != nil {
		return nil, sanitizedErr
	}

	// Return the original error if URL sanitization fails.
	return nil, err
}

func updateRateLimitReset(c *Client, errCheck error) {
	rerr, ok := errCheck.(*AbuseRateLimitError)
	if ok && rerr.RetryAfter != nil {
		c.rateMu.Lock()
		c.RateLimitReset = time.Now().Add(*rerr.RetryAfter)
		c.rateMu.Unlock()
	}
}
func (c *Client) skipCheck(ctx context.Context, req *http.Request) (*str.Response, error) {
	if skip := ctx.Value(skipRateLimitCheck); skip == nil {
		// don't make further requests before Retry After.
		if err := c.CheckRetryAfter(req); err != nil {
			return &str.Response{
				Response: err.Response,
			}, err
		}
	}
	return nil, nil
}

func sanitizeURL(err error) error {
	if e, ok := err.(*url.Error); ok {
		if u, err := url.Parse(e.URL); err == nil {
			e.URL = uri.SanitizeURL(u).String()
			return e
		}
	}
	return nil
}

func parseBoolResponse(err error) (bool, error) {
	if err == nil {
		return true, nil
	}

	if err, ok := err.(*NotFoundError); ok && err.Response.StatusCode == http.StatusNotFound {
		// Simply false. In this one case, we do not pass the error through.
		return false, nil
	}

	// some other real error occurred
	return false, err
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
		return c.genRateLimitError(r, errorResponse)
	case http.StatusUpgradeRequired:
		return c.genUpgradeRequiredError(r, errorResponse)
	case http.StatusNotFound:
		return c.genNotFoundError(r, errorResponse)
	case http.StatusConflict:
		return c.genConflictError(r, errorResponse)
	case http.StatusUnprocessableEntity:
		return c.genValidationError(r, errorResponse)
	default:
		return errorResponse
	}
}

func (*Client) genValidationError(r *http.Response, errorResponse *str.ErrorResponse) *ValidationError {
	validationError := &ValidationError{
		Response: errorResponse.Response,
		Message:  errorResponse.Message,
		Errors:   errorResponse.Errors,
	}
	if r.StatusCode == http.StatusUnprocessableEntity {
		return validationError
	}
	return nil
}

func (*Client) genConflictError(r *http.Response, errorResponse *str.ErrorResponse) *ConflictError {
	conflictError := &ConflictError{
		Response: errorResponse.Response,
		Message:  errorResponse.Message,
	}
	if r.StatusCode == http.StatusConflict {
		return conflictError
	}
	return nil
}

func (*Client) genNotFoundError(r *http.Response, errorResponse *str.ErrorResponse) *NotFoundError {
	notFoundError := &NotFoundError{
		Response: errorResponse.Response,
		Message:  errorResponse.Message,
	}
	if r.StatusCode == http.StatusNotFound {
		return notFoundError
	}
	return nil
}
func (c *Client) genUpgradeRequiredError(r *http.Response, errorResponse *str.ErrorResponse) *UpgradeRequiredError {
	upgradeRequiredError := &UpgradeRequiredError{
		Response: errorResponse.Response,
		Message:  errorResponse.Message,
	}
	if upgradeURL := c.ParseUpgradeUser(r); upgradeURL != nil {
		upgradeRequiredError.UpgradeURL = upgradeURL
		return upgradeRequiredError
	}
	return nil
}

func (c *Client) genRateLimitError(r *http.Response, errorResponse *str.ErrorResponse) *AbuseRateLimitError {
	abuseRateLimitError := &AbuseRateLimitError{
		Response: errorResponse.Response,
		Message:  errorResponse.Message,
	}
	if retryAfter := c.ParseRateLimit(r); retryAfter != nil {
		abuseRateLimitError.RetryAfter = retryAfter
		return abuseRateLimitError
	}
	return nil
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
func (*Client) WithContext(ctx context.Context, req *http.Request) *http.Request {
	return req.WithContext(ctx)
}

// ParseRate parses the rate related headers.
func (*Client) ParseRate(r *http.Response) str.Rate {
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
func (*Client) ParseRateLimit(r *http.Response) *time.Duration {
	// number of seconds that one should
	// wait before resuming making requests.
	if v := r.Header.Get(HeaderRetryAfter); v != "" {
		retryAfterSeconds, _ := strconv.ParseInt(v, consts.BaseInt, consts.BitSize) // Error handling is noop.
		retryAfter := time.Duration(retryAfterSeconds) * time.Second
		return &retryAfter
	}

	return nil
}

// ParseUpgradeUser parses related headers, and returns upgradeUrl.
func (*Client) ParseUpgradeUser(r *http.Response) *url.URL {
	// number of seconds that one should
	// wait before resuming making requests.
	if v := r.Header.Get(HeaderUpgradeURL); v != "" {
		u, _ := url.Parse(v)
		return u
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
