// Package cli for basic cli functions
package cli

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
)

// headers builds canonical header keys, as a real response has them.
func headers(kv ...string) http.Header {
	h := http.Header{}
	for i := 0; i+1 < len(kv); i += 2 {
		h.Set(kv[i], kv[i+1])
	}
	return h
}

func limitsError(h http.Header) error {
	return &internal.UpgradeUserLimitsError{Response: &http.Response{
		StatusCode: 420,
		Header:     h,
		Request:    httptest.NewRequest(http.MethodPost, "https://api.trakt.tv/users/sean/lists", nil),
	}}
}

func TestHandleVIPResponse(t *testing.T) {
	tests := []struct {
		name       string
		resp       *str.Response
		err        error
		browserErr error
		wantOpened string
		wantErr    string
	}{
		{
			name:       "426 opens X-Upgrade-URL",
			resp:       &str.Response{Response: &http.Response{StatusCode: http.StatusUpgradeRequired, Header: headers(internal.HeaderUpgradeURL, "https://trakt.tv/vip?from=api")}},
			wantOpened: "https://trakt.tv/vip?from=api",
			wantErr:    "trakt vip required, browser opened: https://trakt.tv/vip?from=api",
		},
		{
			name:       "426 without header opens the default VIP page",
			resp:       &str.Response{Response: &http.Response{StatusCode: http.StatusUpgradeRequired, Header: http.Header{}}},
			wantOpened: internal.DefaultUpgradeURL,
			wantErr:    "trakt vip required, browser opened: " + internal.DefaultUpgradeURL,
		},
		{
			name:       "420 for a non-VIP user opens X-Upgrade-URL",
			err:        limitsError(headers(internal.HeaderVIPUser, "false", internal.HeaderAccountLimit, "2", internal.HeaderUpgradeURL, "https://trakt.tv/vip")),
			wantOpened: "https://trakt.tv/vip",
			wantErr:    "account limit exceeded (limit: 2), browser opened: https://trakt.tv/vip",
		},
		{
			name:    "420 for a VIP user only reports the limit",
			err:     limitsError(headers(internal.HeaderVIPUser, "true", internal.HeaderAccountLimit, "100")),
			wantErr: "account limit exceeded (limit: 100)",
		},
		{
			name:    "wrapped 420 is still handled",
			err:     fmt.Errorf("add personal list error: %w", limitsError(headers(internal.HeaderVIPUser, "true"))),
			wantErr: "account limit exceeded",
		},
		{
			name:       "browser failure keeps the URL in the error",
			resp:       &str.Response{Response: &http.Response{StatusCode: http.StatusUpgradeRequired, Header: http.Header{}}},
			browserErr: errors.New("no display"),
			wantOpened: internal.DefaultUpgradeURL,
			wantErr:    "trakt vip required, open " + internal.DefaultUpgradeURL + " to upgrade (browser error: no display)",
		},
		{name: "404 is not VIP related", resp: &str.Response{Response: &http.Response{StatusCode: http.StatusNotFound}}, err: errors.New("not found")},
		{name: "nil response and error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opened := ""
			openBrowser = func(url string) error {
				opened = url
				return tt.browserErr
			}
			t.Cleanup(func() { openBrowser = OpenBrowser })

			err := HandleVIPResponse(tt.resp, tt.err)
			if opened != tt.wantOpened {
				t.Errorf("opened %q, want %q", opened, tt.wantOpened)
			}
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("error is %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
			}
		})
	}
}
