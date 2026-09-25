// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

// UsersActivitiesHandler struct for handler
type UsersActivitiesHandler struct{}

// Handle to handle users: activities action
func (h UsersActivitiesHandler) Handle(options *str.Options, client *internal.Client) error {
	// IsValidConfigType accepts an empty value, but the route needs a type
	if len(options.Type) == consts.ZeroValue || !cfg.IsValidConfigType(cfg.SocialActivityTypes, options.Type) {
		return fmt.Errorf("set -t to one of %v for activities", cfg.SocialActivityTypes)
	}

	printer.Println("Returns recent activity of your " + options.Type)
	opts := uri.SocialActivityOptions{
		Limit:     options.PerPage,
		Extended:  options.ExtendedInfo,
		Genres:    options.Genres,
		Years:     options.Years,
		Runtimes:  options.Runtimes,
		Countries: options.Countries,
	}
	result, err := h.fetchActivities(client, options, &opts, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch %s activities error: %w", options.Type, err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	return writeResult(options, result)
}

func (h UsersActivitiesHandler) fetchActivities(client *internal.Client, options *str.Options, opts *uri.SocialActivityOptions, page int) ([]*str.SocialActivity, error) {
	opts.Page = page
	list, resp, err := client.Users.GetSocialActivity(client.BuildCtxFromOptions(options), &options.UserName, &options.Type, opts)
	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		nextPageItems, err := h.fetchActivities(client, options, opts, page+consts.NextPageStep)
		if err != nil {
			return nil, err
		}
		list = append(list, nextPageItems...)
	}

	return list, nil
}
