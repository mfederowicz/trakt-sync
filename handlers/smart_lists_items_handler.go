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

// SmartListsItemsHandler struct for handler
type SmartListsItemsHandler struct{}

// Handle to handle smart_lists: items action
func (h SmartListsItemsHandler) Handle(options *str.Options, client *internal.Client) error {
	if err := validSmartListID(options); err != nil {
		return err
	}
	if len(options.WatchNow) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.WatchNowFilters, options.WatchNow) {
		return fmt.Errorf("watchnow '%s' is not valid, avaliable values: %v", options.WatchNow, cfg.WatchNowFilters)
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

func (h SmartListsItemsHandler) fetchItems(client *internal.Client, options *str.Options, opts *uri.SmartListItemsOptions, page int) ([]*str.UserListItem, error) {
	opts.Page = page
	list, resp, err := client.SmartLists.GetSmartListItems(client.BuildCtxFromOptions(options), &options.InternalID, opts)
	if err = smartListError(options.Action, options.InternalID, resp, err); err != nil {
		return nil, err
	}

	// Check if there are more pages
	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		nextPageItems, err := h.fetchItems(client, options, opts, page+consts.NextPageStep)
		if err != nil {
			return nil, err
		}
		list = append(list, nextPageItems...)
	}

	return list, nil
}
