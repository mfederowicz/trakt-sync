// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// YounifyHandler interface to handle younify module action
type YounifyHandler interface {
	Handle(options *str.Options, client *trakt.Client) error
}

// validServiceID checks the streaming service id needed by connect, refresh and disconnect.
func validServiceID(options *str.Options) error {
	if len(options.ServiceID) == consts.ZeroValue {
		return errors.New(consts.EmptyServiceIDMsg)
	}
	return nil
}

// younifyError maps a younify response to a readable error: not open to API apps (401 / 403), unknown service (404),
// VIP gating (422 / 426) and a rejected return_url (400).
func younifyError(action string, serviceID string, resp *str.Response, err error) error {
	// Trakt answered 401 on 2026-09-25 and 403 on 2026-10-31 while other OAuth routes worked,
	// so neither is an expired token of this app
	if resp != nil && (resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden) {
		return fmt.Errorf(consts.YounifyUnauthorizedMsg, action, resp.StatusCode, err)
	}

	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found streaming service for:%s", serviceID)
	}

	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}

	// Client.Do turns a 422 into a plain "validation error", so check the status
	if resp != nil && resp.StatusCode == http.StatusUnprocessableEntity {
		return fmt.Errorf("streaming service %s is not connectable on your plan (VIP only?): %w", serviceID, err)
	}

	var badRequest *trakt.BadRequestError
	if errors.As(err, &badRequest) {
		return fmt.Errorf("younify %s rejected (return_url must be trakt://... or https://*.trakt.tv): %w", action, err)
	}

	if err != nil {
		return fmt.Errorf("younify %s error: %w", action, err)
	}

	return nil
}
