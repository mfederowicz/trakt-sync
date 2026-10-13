// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// UsersYearInReviewPersonaHandler struct for handler
type UsersYearInReviewPersonaHandler struct{}

// Handle to handle users: year_in_review_persona action
func (UsersYearInReviewPersonaHandler) Handle(options *str.Options, client *trakt.Client) error {
	if options.Year < consts.PersonaMinYear {
		return fmt.Errorf("set -year %d or later for %s ie: -year %d", consts.PersonaMinYear, options.Action, consts.PersonaMinYear)
	}

	printer.Printf("Returns year in review persona %d for: %s\n", options.Year, options.UserName)
	result, resp, err := client.Users.GetYearInReviewPersona(cli.ContextFromOptions(options), options.UserName, options.Year)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no %s for %d (user %s)", options.Action, options.Year, options.UserName)
	}
	// the contract answers 403 vip_required when a non-VIP user asks for their own persona
	var forbidden *trakt.ForbiddenError
	if errors.As(err, &forbidden) {
		return fmt.Errorf("%s is VIP only: %w", options.Action, err)
	}
	if apiErr := notOpenToAPIApps(options.Action, err); apiErr != nil {
		return apiErr
	}
	if err != nil {
		return fmt.Errorf("fetch %s error: %w", options.Action, err)
	}
	if result == nil || result.Persona == nil {
		return fmt.Errorf("no %s for %d (user %s)", options.Action, options.Year, options.UserName)
	}

	return writeResult(options, result)
}
