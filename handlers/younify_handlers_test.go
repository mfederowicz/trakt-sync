// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/test"
)

func TestYounifyHandlers(t *testing.T) {
	tests := []struct {
		name       string
		handler    Handler
		options    str.Options
		method     string
		path       string
		status     int
		body       string
		wantBody   string
		wantErr    string
		wantNoAPI  bool
		wantOutput bool
	}{
		{name: "connections", handler: YounifyConnectionsHandler{}, options: str.Options{Action: consts.Connections},
			method: http.MethodGet, path: "/younify/connections", status: http.StatusOK, body: `[{"id":"netflix","connected":true}]`, wantOutput: true},
		{name: "connections empty", handler: YounifyConnectionsHandler{}, options: str.Options{Action: consts.Connections},
			method: http.MethodGet, path: "/younify/connections", status: http.StatusOK, body: `[]`, wantErr: consts.EmptyResult},
		{name: "connect default return_url", handler: YounifyConnectHandler{}, options: str.Options{Action: consts.Connect, ServiceID: "netflix"},
			method: http.MethodPost, path: "/younify/connect", status: http.StatusOK, body: `{"url":"https://auth.example/x"}`,
			wantBody: `{"service_id":"netflix","return_url":"https://trakt.tv"}`, wantOutput: true},
		{name: "connect own return_url", handler: YounifyConnectHandler{}, options: str.Options{Action: consts.Connect, ServiceID: "netflix", ReturnURL: "trakt://settings"},
			method: http.MethodPost, path: "/younify/connect", status: http.StatusOK, body: `{"url":"https://auth.example/x"}`,
			wantBody: `{"service_id":"netflix","return_url":"trakt://settings"}`, wantOutput: true},
		{name: "connect not connectable", handler: YounifyConnectHandler{}, options: str.Options{Action: consts.Connect, ServiceID: "hulu"},
			method: http.MethodPost, path: "/younify/connect", status: http.StatusUnprocessableEntity, body: `{}`, wantErr: "hulu is not connectable on your plan"},
		{name: "connect bad return_url", handler: YounifyConnectHandler{}, options: str.Options{Action: consts.Connect, ServiceID: "netflix", ReturnURL: "https://example.com"},
			method: http.MethodPost, path: "/younify/connect", status: http.StatusBadRequest, body: `{}`, wantErr: "return_url must be trakt://"},
		{name: "connect empty url", handler: YounifyConnectHandler{}, options: str.Options{Action: consts.Connect, ServiceID: "netflix"},
			method: http.MethodPost, path: "/younify/connect", status: http.StatusOK, body: `{}`, wantErr: "no web auth URL"},
		{name: "connect without service", handler: YounifyConnectHandler{}, options: str.Options{Action: consts.Connect},
			wantErr: consts.EmptyServiceIDMsg, wantNoAPI: true},
		{name: "refresh", handler: YounifyRefreshHandler{}, options: str.Options{Action: consts.Refresh, ServiceID: "netflix"},
			method: http.MethodPost, path: "/younify/users/refresh/netflix", status: http.StatusNoContent},
		{name: "refresh all data", handler: YounifyRefreshHandler{}, options: str.Options{Action: consts.Refresh, ServiceID: "netflix", AllData: true},
			method: http.MethodPost, path: "/younify/users/refresh/netflix/all_data", status: http.StatusNoContent},
		{name: "refresh unknown service", handler: YounifyRefreshHandler{}, options: str.Options{Action: consts.Refresh, ServiceID: "nope"},
			method: http.MethodPost, path: "/younify/users/refresh/nope", status: http.StatusNotFound, body: `{}`, wantErr: "not found streaming service for:nope"},
		{name: "refresh without service", handler: YounifyRefreshHandler{}, options: str.Options{Action: consts.Refresh},
			wantErr: consts.EmptyServiceIDMsg, wantNoAPI: true},
		{name: "disconnect", handler: YounifyDisconnectHandler{}, options: str.Options{Action: consts.Disconnect, ServiceID: "netflix"},
			method: http.MethodDelete, path: "/younify/users/services/netflix", status: http.StatusNoContent},
		{name: "disconnect without service", handler: YounifyDisconnectHandler{}, options: str.Options{Action: consts.Disconnect},
			wantErr: consts.EmptyServiceIDMsg, wantNoAPI: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := setup(t)
			defer s.Teardown()

			calls := 0
			s.Mux.HandleFunc("/younify/", func(w http.ResponseWriter, r *http.Request) {
				calls++
				test.AssertMethod(t, r, tt.method)
				if r.URL.Path != tt.path {
					t.Errorf("path is %q, want %q", r.URL.Path, tt.path)
				}
				if tt.wantBody != "" {
					got := map[string]string{}
					test.AssertNilError(t, json.NewDecoder(r.Body).Decode(&got))
					want := map[string]string{}
					test.AssertNilError(t, json.Unmarshal([]byte(tt.wantBody), &want))
					test.AssertNoDiff(t, want, got)
				}
				w.WriteHeader(tt.status)
				test.SafeFprint(w, tt.body)
			})

			options := tt.options
			options.Output = filepath.Join(t.TempDir(), "out.json")
			err := tt.handler.Handle(&options, s.Client)
			if tt.wantNoAPI && calls != 0 {
				t.Error("API was called with invalid options")
			}
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error is %v, want it to contain %q", err, tt.wantErr)
				}
				if _, statErr := os.Stat(options.Output); !os.IsNotExist(statErr) {
					t.Errorf("output file written on error: %v", statErr)
				}
				return
			}
			test.AssertNilError(t, err)
			if calls != 1 {
				t.Errorf("API calls = %d, want 1", calls)
			}
			_, statErr := os.Stat(options.Output)
			if tt.wantOutput && statErr != nil {
				t.Errorf("output file was not written: %v", statErr)
			}
			if !tt.wantOutput && !os.IsNotExist(statErr) {
				t.Errorf("unexpected output file: %v", statErr)
			}
		})
	}
}
