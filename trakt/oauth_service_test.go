package trakt

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
	"github.com/stretchr/testify/assert"
)

// TestNewClientHosts checks the default hosts: OAuth on auth.trakt.tv, everything else on api.trakt.tv.
func TestNewClientHosts(t *testing.T) {
	c := NewClient(nil)
	assert.Equal(t, "https://auth.trakt.tv/", c.AuthURL.String())
	assert.Equal(t, "https://api.trakt.tv/", c.BaseURL.String())
	assert.Equal(t, c.AuthURL, c.WithAuthToken("token").AuthURL, "a copy keeps AuthURL")
}

// TestOauthServiceUsesAuthURL checks every OAuth call goes to AuthURL and none reaches BaseURL.
func TestOauthServiceUsesAuthURL(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		call func(c *Client) (*str.Response, error)
	}{
		{name: "device code", path: "/oauth/device/code", body: `{"device_code":"d","user_code":"u","verification_url":"https://auth.trakt.tv/activate"}`,
			call: func(c *Client) (*str.Response, error) {
				_, resp, err := c.Oauth.GenerateNewDeviceCodes(context.Background(), &str.NewDeviceCode{ClientID: str.String("id")})
				return resp, err
			}},
		{name: "device token", path: "/oauth/device/token", body: `{"access_token":"a"}`,
			call: func(c *Client) (*str.Response, error) {
				_, resp, err := c.Oauth.PollForAccessToken(context.Background(), &str.NewDeviceToken{Code: str.String("d")})
				return resp, err
			}},
		{name: "refresh token", path: "/oauth/token", body: `{"access_token":"a"}`,
			call: func(c *Client) (*str.Response, error) {
				_, resp, err := c.Oauth.ExchangeRefreshTokenForAccessToken(context.Background(), &str.CurrentDeviceToken{RefreshToken: str.String("r")})
				return resp, err
			}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var apiCalls atomic.Int32
			api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				apiCalls.Add(1)
				w.WriteHeader(http.StatusNotFound)
			}))
			defer api.Close()
			auth := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				test.AssertMethod(t, r, http.MethodPost)
				assert.Equal(t, tt.path, r.URL.Path)
				test.SafeFprint(w, tt.body)
			}))
			defer auth.Close()

			c := NewClient(nil)
			c.BaseURL, _ = url.Parse(api.URL + "/")
			c.AuthURL, _ = url.Parse(auth.URL + "/")

			resp, err := tt.call(c)
			test.AssertNilError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Zero(t, apiCalls.Load(), "an OAuth call reached BaseURL")
		})
	}
}

// TestOauthServiceAuthURLTrailingSlash checks an AuthURL without a trailing slash is refused before any request.
func TestOauthServiceAuthURLTrailingSlash(t *testing.T) {
	c := NewClient(nil)
	c.AuthURL, _ = url.Parse("https://auth.example.com/v1")
	_, _, err := c.Oauth.GenerateNewDeviceCodes(context.Background(), new(str.NewDeviceCode))
	assert.EqualError(t, err, `AuthURL must have a trailing slash, but "https://auth.example.com/v1" does not`)
}
