// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// UsersMonthInReviewHandler struct for handler
type UsersMonthInReviewHandler struct{}

// Handle to handle users: month_in_review action
func (UsersMonthInReviewHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.Year <= consts.ZeroValue {
		return errors.New("set -year for month_in_review ie: -year 2025 -month 8")
	}
	if options.Month < consts.OneValue || options.Month > consts.MonthsInYear {
		return errors.New("set -month 1-12 for month_in_review ie: -year 2025 -month 8")
	}

	period := fmt.Sprintf(consts.YearMonthFormat, options.Year, options.Month)
	printer.Println("Returns month in review " + period + " for: " + options.UserName)
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Users.GetMonthInReview(client.BuildCtxFromOptions(options), &options.UserName, options.Year, options.Month, &opts)
	return writeReview(options, period, result, resp, err)
}

// UsersYearInReviewHandler struct for handler
type UsersYearInReviewHandler struct{}

// Handle to handle users: year_in_review action
func (UsersYearInReviewHandler) Handle(options *str.Options, client *internal.Client) error {
	if options.Year <= consts.ZeroValue {
		return errors.New("set -year for year_in_review ie: -year 2025")
	}

	period := fmt.Sprint(options.Year)
	printer.Println("Returns year in review " + period + " for: " + options.UserName)
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, resp, err := client.Users.GetYearInReview(client.BuildCtxFromOptions(options), &options.UserName, options.Year, &opts)
	return writeReview(options, period, result, resp, err)
}

// writeReview maps a review response to a readable error, or writes it; an empty {} review counts as not found.
func writeReview(options *str.Options, period string, result *str.Review, resp *str.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("no %s for %s (user %s)", options.Action, period, options.UserName)
	}
	if apiErr := notOpenToAPIApps(options.Action, err); apiErr != nil {
		return apiErr
	}
	if err != nil {
		return fmt.Errorf("fetch %s error: %w", options.Action, err)
	}
	if result == nil || (result.Stats == nil && result.Images == nil) {
		return fmt.Errorf("no %s for %s (user %s)", options.Action, period, options.UserName)
	}

	return writeResult(options, result)
}
