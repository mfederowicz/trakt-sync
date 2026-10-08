// Package handlers used to handle module actions
package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// unknownValueError is returned after a types or items usage, so a wrong or missing -t / -item exits with status 1.
func unknownValueError(flagName string, value string) error {
	if value == consts.EmptyString {
		return fmt.Errorf("no %s given", flagName)
	}
	return fmt.Errorf("unknown %s %q", flagName, value)
}

// episodeTitle returns the title of an episode for messages; the API sends none for many upcoming episodes.
func episodeTitle(episode *str.Episode) string {
	if episode.Title == nil {
		return consts.NoEpisodeTitle
	}
	return *episode.Title
}

func isMovieType(stype string) bool {
	switch stype {
	case consts.Movie, consts.Movies:
		return true
	default:
		return false
	}
}

func isShowType(stype string) bool {
	switch stype {
	case consts.Show, consts.Shows:
		return true
	default:
		return false
	}
}

func isSeasonType(stype string) bool {
	switch stype {
	case consts.Season, consts.Seasons:
		return true
	default:
		return false
	}
}

func isEpisodeType(stype string) bool {
	switch stype {
	case consts.Episode, consts.Episodes:
		return true
	default:
		return false
	}
}

func isPeopleType(stype string) bool {
	switch stype {
	case consts.People:
		return true
	default:
		return false
	}
}

func isUserType(stype string) bool {
	switch stype {
	case consts.User:
		return true
	default:
		return false
	}
}

// Media interface for helpers
type Media interface {
	str.Movie | str.Show | str.Episode | str.Season
}

// onlyIDs is a helper function to extract each type objects with only ids
func onlyIDs[T Media](items []str.ExportlistItem) []T {
	result := make([]T, 0, len(items))

	var zero T

	switch any(zero).(type) {
	case str.Movie:
		for _, item := range items {
			result = append(result, any(str.Movie{IDs: item.IDs}).(T))
		}
	case str.Show:
		for _, item := range items {
			if item.Seasons != nil && len(*item.Seasons) > 0 {
				updatedSeasons := SeasonsWithEpisodeNumbersOnly(item.Seasons)
				result = append(result, any(str.Show{IDs: item.IDs, Seasons: updatedSeasons}).(T))
			} else {
				result = append(result, any(str.Show{IDs: item.IDs}).(T))
			}
		}
	case str.Episode:
		for _, item := range items {
			result = append(result, any(str.Episode{IDs: item.IDs}).(T))
		}
	case str.Season:
		for _, item := range items {
			result = append(result, any(str.Season{IDs: item.IDs}).(T))
		}
	default:
		panic("unsupported type")
	}

	return result
}

// Ptr is a helper routine that allocates a new T value
// to store v and returns a pointer to it.
func Ptr[T any](v T) *T {
	return &v
}

// pageDelay is the pause between two API calls of one action, such as the next page of a list; tests set it to zero.
var pageDelay = time.Duration(consts.SleepNumberOfSeconds) * time.Second

// waitPageDelay pauses before the next API call of the same action.
func waitPageDelay() {
	time.Sleep(pageDelay)
}

// mediaFilters builds the media filters of a list route from the filter flags.
func mediaFilters(options *str.Options) uri.MediaFilters {
	return uri.MediaFilters{
		WatchNow:       options.WatchNow,
		Genres:         options.Genres,
		Subgenres:      options.Subgenres,
		Years:          options.Years,
		Ratings:        options.Ratings,
		StartDate:      options.MediaStartDate,
		EndDate:        options.MediaEndDate,
		Runtimes:       options.Runtimes,
		Countries:      options.Countries,
		Certifications: options.Certifications,
		Languages:      options.Languages,
		ImdbRatings:    options.ImdbRatings,
		RtMeters:       options.RtMeters,
		RtUserMeters:   options.RtUserMeters,
	}
}

// checkMediaFilters reports a -watchnow value the API does not know.
func checkMediaFilters(options *str.Options) error {
	if len(options.WatchNow) > consts.ZeroValue && !cfg.IsValidConfigType(cfg.WatchNowFilters, options.WatchNow) {
		return fmt.Errorf("watchnow '%s' is not valid, available values: %v", options.WatchNow, cfg.WatchNowFilters)
	}
	return nil
}

// showStatus builds the status filter of a shows list route from the -status flag.
func showStatus(options *str.Options) []string {
	if len(options.ShowStatus) == consts.ZeroValue {
		return nil
	}
	return strings.Split(options.ShowStatus, consts.SeparatorString)
}

// checkShowFilters reports a -watchnow or -status value the API does not know.
func checkShowFilters(options *str.Options) error {
	if err := checkMediaFilters(options); err != nil {
		return err
	}
	for _, status := range showStatus(options) {
		if !slices.Contains(uri.StatusOptions, status) {
			return fmt.Errorf("status '%s' is not valid, available values: %v", status, uri.StatusOptions)
		}
	}
	return nil
}

// movieStatus builds the status filter of a movies list route from the -status flag.
func movieStatus(options *str.Options) []string {
	if len(options.MovieStatus) == consts.ZeroValue {
		return nil
	}
	return strings.Split(options.MovieStatus, consts.SeparatorString)
}

// checkMovieFilters reports a -watchnow or -status value the API does not know.
func checkMovieFilters(options *str.Options) error {
	if err := checkMediaFilters(options); err != nil {
		return err
	}
	for _, status := range movieStatus(options) {
		if !slices.Contains(cfg.MovieStatusFilters, status) {
			return fmt.Errorf("status '%s' is not valid, available values: %v", status, cfg.MovieStatusFilters)
		}
	}
	return nil
}

// pageFetcher fetches one page of a paginated list.
type pageFetcher[T any] func(opts *uri.ListOptions) ([]T, *str.Response, error)

// fetchAllPages fetches a list starting at page, following pages while client.HavePages allows.
func fetchAllPages[T any](client *trakt.Client, options *str.Options, page int, fetch pageFetcher[T]) ([]T, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := fetch(&opts)
	if err != nil {
		return nil, err
	}

	if client.HavePages(page, resp, options.PagesLimit) {
		waitPageDelay()
		nextPageItems, err := fetchAllPages(client, options, page+consts.NextPageStep, fetch)
		if err != nil {
			return nil, err
		}
		list = append(list, nextPageItems...)
	}

	return list, nil
}

// reportMessage returns the API message of a failed report, or the error text when the response has no message.
func reportMessage(message *string, err error) string {
	if message != nil && len(*message) > consts.ZeroValue {
		return *message
	}
	if err != nil {
		return err.Error()
	}
	return consts.EmptyString
}

// validIDCountryOptions checks the id and country needed by the watch now routes.
func validIDCountryOptions(options *str.Options, emptyIDMsg string) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(emptyIDMsg)
	}

	if len(options.Country) == consts.ZeroValue {
		return errors.New(consts.EmptyCountryMsg)
	}

	return nil
}

// watchNowError maps a watch now response to a readable error: 404, VIP limits and limited access (403).
func watchNowError(action string, kind string, id string, resp *str.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found %s for:%s", kind, id)
	}

	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}

	var forbidden *trakt.ForbiddenError
	if errors.As(err, &forbidden) {
		return fmt.Errorf(consts.LimitedAccessMsg, action, err)
	}

	if err != nil {
		return fmt.Errorf("fetch %s error: %w", action, err)
	}

	return nil
}

// isEmptySentiments reports a sentiments response with no data; the API answers an unknown id with {} instead of 404.
func isEmptySentiments(s *str.Sentiments) bool {
	return s == nil || (len(s.Good) == consts.ZeroValue && len(s.Bad) == consts.ZeroValue && s.CommentCount == nil && s.AnalyzedAt == nil)
}

// parseItemTraktID parses the numeric Trakt ID of a season or episode; those have no slugs.
func parseItemTraktID(id string) (int64, error) {
	traktID, err := strconv.ParseInt(id, consts.BaseInt, consts.BitSize)
	if err != nil || traktID <= consts.ZeroValueInt64 {
		return consts.ZeroValueInt64, fmt.Errorf("trakt id must be a positive number, got %q", id)
	}

	return traktID, nil
}

// writeResult marshals data and writes it to the output file.
func writeResult(options *str.Options, data any) error {
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(data, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal %s error: %w", options.Action, err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}

// notOpenToAPIApps turns a 401 on a route that other OAuth calls pass into a readable error; nil for any other error.
func notOpenToAPIApps(action string, err error) error {
	var invalidUser *trakt.InvalidUserError
	if errors.As(err, &invalidUser) {
		return fmt.Errorf(consts.NotOpenToAPIAppsMsg, action, err)
	}
	return nil
}

// sortRouteType maps -t to the {type} of the users watchlist / favorites /{type}/{sort} routes and checks -sort.
func sortRouteType(section string, options *str.Options) (string, error) {
	if !slices.Contains(cfg.UsersSortPathValues, options.SortPath) {
		return consts.EmptyString, fmt.Errorf("sort '%s' is not valid, avaliable values: %v", options.SortPath, cfg.UsersSortPathValues)
	}
	switch options.Type {
	case consts.Movies, consts.Shows:
		return options.Type, nil
	case consts.ActionTypeAll:
		if section == consts.Favorites {
			return consts.Media, nil
		}
		return consts.MovieShow, nil
	default:
		return consts.EmptyString, fmt.Errorf("-sort works with -t all, movies or shows, not '%s'", options.Type)
	}
}

// readStrictInput decodes the -items file or stdin into v; unknown keys are rejected so a typo fails before the request.
func readStrictInput(common *CommonLogic, options *str.Options, v any) error {
	data, err := common.ReadInputBytes(*options)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("invalid %s JSON: %w", options.Action, err)
	}

	return nil
}

// enumCheck is one optional input value that must be one of valid when it is sent.
type enumCheck struct {
	name  string
	value *string
	valid []string
}

// checkEnums returns an error for the first sent value that is not a contract value; "" is not valid.
func checkEnums(checks []enumCheck) error {
	for _, c := range checks {
		if c.value != nil && !slices.Contains(c.valid, *c.value) {
			return fmt.Errorf("%s '%s' is not valid, avaliable values: %v", c.name, *c.value, c.valid)
		}
	}
	return nil
}
