// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
)

// parentalRange is a parental guide severity range min-max, from 0 (none) to 3 (severe).
var parentalRange = regexp.MustCompile(`^[0-3]-[0-3]$`)

// watchNowCountry is the two-letter lowercase region of the watchnow filter.
var watchNowCountry = regexp.MustCompile(`^[a-z]{2}$`)

// SmartListsItemsHandler struct for handler
type SmartListsItemsHandler struct{}

// Handle to handle smart_lists: items action
func (h SmartListsItemsHandler) Handle(options *str.Options, client *trakt.Client) error {
	if err := validSmartListID(options); err != nil {
		return err
	}
	if len(options.WatchNow) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.WatchNowFilters, options.WatchNow) {
		return fmt.Errorf("watchnow '%s' is not valid, avaliable values: %v", options.WatchNow, cfg.WatchNowFilters)
	}
	if err := validSmartListItemsFilters(options); err != nil {
		return err
	}
	if err := smartListItemsSort(options); err != nil {
		return err
	}

	printer.Println("Returns the items of smart list: " + options.InternalID)
	opts := uri.SmartListItemsOptions{
		Limit:             options.PerPage,
		Extended:          options.ExtendedInfo,
		WatchNow:          options.WatchNow,
		Genres:            options.Genres,
		Subgenres:         options.Subgenres,
		Years:             options.Years,
		Ratings:           options.Ratings,
		Runtimes:          options.Runtimes,
		Countries:         options.Countries,
		Certifications:    options.Certifications,
		IgnoreWatched:     options.IgnoreWatched,
		IgnoreWatchlisted: options.IgnoreWatchlisted,
		WatchNowCountry:   options.WatchNowCountry,

		ParentalAlcohol:        options.Parental.Alcohol,
		ParentalFrightening:    options.Parental.Frightening,
		ParentalIncludeUnrated: options.Parental.IncludeUnrated,
		ParentalNudity:         options.Parental.Nudity,
		ParentalProfanity:      options.Parental.Profanity,
		ParentalViolence:       options.Parental.Violence,
	}
	result, err := h.fetchItems(client, options, &opts, consts.DefaultPage)
	if err != nil {
		return err
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	return writeSmartList(options, result)
}

func (h SmartListsItemsHandler) fetchItems(client *trakt.Client, options *str.Options, opts *uri.SmartListItemsOptions, page int) ([]*str.UserListItem, error) {
	opts.Page = page
	list, resp, err := h.fetchPage(client, options, opts)
	if err = smartListError(options.Action, options.InternalID, resp, err); err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		waitPageDelay()
		nextPageItems, err := h.fetchItems(client, options, opts, page+consts.NextPageStep)
		if err != nil {
			return nil, err
		}
		list = append(list, nextPageItems...)
	}

	return list, nil
}

// fetchPage reads one page from the plain items route, or from the typed and sorted one when smartListItemsSort set a type.
func (SmartListsItemsHandler) fetchPage(client *trakt.Client, options *str.Options, opts *uri.SmartListItemsOptions) ([]*str.UserListItem, *str.Response, error) {
	ctx := cli.ContextFromOptions(options)
	if len(options.Type) == consts.ZeroValue {
		return client.SmartLists.GetSmartListItems(ctx, options.InternalID, opts)
	}

	return client.SmartLists.GetSmartListItemsByTypeAndSort(ctx, options.InternalID, options.Type, options.SortBy, options.SortHow, opts)
}

// smartListItemsSort checks -t, -sort_by and -sort_how; when at least one is given the others get the API defaults,
// because the typed and sorted route needs all three. With none given they stay empty and the plain route is used.
func smartListItemsSort(options *str.Options) error {
	if len(options.Type) == consts.ZeroValue && len(options.SortBy) == consts.ZeroValue && len(options.SortHow) == consts.ZeroValue {
		return nil
	}
	if len(options.Type) == consts.ZeroValue {
		options.Type = consts.ActionTypeAll
	}
	if len(options.SortBy) == consts.ZeroValue {
		options.SortBy = consts.SortRank
	}
	if len(options.SortHow) == consts.ZeroValue {
		options.SortHow = consts.SortDesc
		if options.SortBy == consts.SortRank || options.SortBy == consts.SortTitle {
			options.SortHow = consts.SortAsc
		}
	}

	return checkEnums([]enumCheck{
		{name: "type", value: &options.Type, valid: cfg.SmartListItemsTypes},
		{name: "sort_by", value: &options.SortBy, valid: cfg.SmartListItemsSortBy},
		{name: "sort_how", value: &options.SortHow, valid: cfg.SmartListItemsSortHow},
	})
}

// validSmartListItemsFilters checks the format of -watchnow_country and the parental guide ranges.
func validSmartListItemsFilters(options *str.Options) error {
	if len(options.WatchNowCountry) > consts.ZeroValue && !watchNowCountry.MatchString(options.WatchNowCountry) {
		return fmt.Errorf("watchnow_country '%s' is not valid, want a 2 character lowercase code ie: us", options.WatchNowCountry)
	}

	ranges := []struct{ name, value string }{
		{"parental_alcohol", options.Parental.Alcohol},
		{"parental_frightening", options.Parental.Frightening},
		{"parental_nudity", options.Parental.Nudity},
		{"parental_profanity", options.Parental.Profanity},
		{"parental_violence", options.Parental.Violence},
	}
	for _, r := range ranges {
		if len(r.value) > consts.ZeroValue && !parentalRange.MatchString(r.value) {
			return fmt.Errorf("%s '%s' is not valid, want a range min-max from 0 to 3 ie: 0-1", r.name, r.value)
		}
	}

	return nil
}
