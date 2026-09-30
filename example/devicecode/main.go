// Package main runs the OAuth device flow and prints the access token as JSON. The user opens the
// verification URL, enters the code, and approves the app:
//
//	TRAKT_CLIENT_ID=... TRAKT_CLIENT_SECRET=... go run ./devicecode
//
// Store the printed token somewhere safe; it is a credential for the user's account.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

func main() {
	clientID := os.Getenv("TRAKT_CLIENT_ID")
	clientSecret := os.Getenv("TRAKT_CLIENT_SECRET")
	client := trakt.NewClient(nil).WithClientID(clientID).WithUserAgent("trakt-sync-example/1.0")
	ctx := context.Background()

	code, _, err := client.Oauth.GenerateNewDeviceCodes(ctx, &str.NewDeviceCode{ClientID: &clientID})
	if err != nil {
		log.Fatalf("device code: %v", err)
	}
	fmt.Printf("Open %s and enter the code %s\n", code.VerificationURL, code.UserCode)

	token, err := waitForToken(ctx, client, code, clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}

	out, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		log.Fatalf("encode token: %v", err)
	}
	fmt.Println(string(out))
}

// waitForToken polls every code.Interval seconds until the user approves the app or the code expires.
// Trakt answers 400 while the approval is pending and 429 when polling too fast; any other error ends the flow.
func waitForToken(ctx context.Context, client *trakt.Client, code *str.DeviceCode, clientID, clientSecret string) (*str.DeviceToken, error) {
	request := &str.NewDeviceToken{Code: &code.DeviceCode, ClientID: &clientID, ClientSecret: &clientSecret}
	interval := time.Duration(code.Interval) * time.Second
	deadline := time.Now().Add(time.Duration(code.ExpiresIn) * time.Second)

	for time.Now().Before(deadline) {
		time.Sleep(interval)

		token, resp, err := client.Oauth.PollForAccessToken(ctx, request)
		if err == nil {
			return token, nil
		}
		var pending *trakt.BadRequestError
		if errors.As(err, &pending) && resp != nil && resp.StatusCode == http.StatusBadRequest {
			continue
		}
		var slowDown *trakt.AbuseRateLimitError
		if errors.As(err, &slowDown) {
			interval += time.Second
			continue
		}
		return nil, fmt.Errorf("device token: %w", err)
	}
	return nil, errors.New("device code expired before it was approved")
}
