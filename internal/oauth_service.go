// Package internal used for client and services
package internal

import (
	"context"
	"net/http"

	"github.com/mfederowicz/trakt-sync/str"
)

// OauthService handles communication with the oauth related
// methods of the Trakt API.
//
// API docs: https://trakt.docs.apiary.io/#reference/authentication-oauth
type OauthService Service

// GenerateNewDeviceCodes Generate new codes to start the device authentication process.
//
// API docs: https://trakt.docs.apiary.io/#reference/authentication-devices/device-code/generate-new-device-codes
func (o *OauthService) GenerateNewDeviceCodes(ctx context.Context, code *str.NewDeviceCode) (*str.DeviceCode, *str.Response, error) {
	u := "oauth/device/code"
	req, err := o.client.NewRequest("POST", u, code)
	if err != nil {
		return nil, nil, err
	}

	d := new(str.DeviceCode)
	resp, err := o.client.Do(ctx, req, d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// PoolForTheAccessToken Use the device_code and poll at the interval (in seconds) to check if the user has authorized you app.
//
// API docs: https://trakt.docs.apiary.io/#reference/authentication-devices/get-token/poll-for-the-access_token
func (o *OauthService) PoolForTheAccessToken(ctx context.Context, deviceToken *str.NewDeviceToken) (*str.DeviceToken, *str.Response, error) {
	u := "oauth/device/token"
	req, err := o.client.NewRequest("POST", u, deviceToken)
	if err != nil {
		return nil, nil, err
	}

	d := new(str.DeviceToken)
	resp, err := o.client.Do(ctx, req, d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}

// ExchangeRefreshTokenForAccessToken Use the refresh_token to get a new access_token
// without asking the user to re-authenticate. The access_token is valid for 3 months
// before it needs to be refreshed again.
//
// API docs: https://trakt.docs.apiary.io/#reference/authentication-oauth/get-token/exchange-refresh_token-for-access_token
func (o *OauthService) ExchangeRefreshTokenForAccessToken(ctx context.Context, deviceToken *str.CurrentDeviceToken) (*str.DeviceToken, *str.Response, error) {
	u := "oauth/token"
	req, err := o.client.NewRequest(http.MethodPost, u, deviceToken)
	if err != nil {
		return nil, nil, err
	}

	d := new(str.DeviceToken)
	resp, err := o.client.Do(ctx, req, d)
	if err != nil {
		return nil, resp, err
	}

	return d, resp, nil
}
