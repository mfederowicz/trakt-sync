// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mfederowicz/trakt-sync/cli"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
)

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

// pageFetcher fetches one page of a paginated list.
type pageFetcher[T any] func(opts *uri.ListOptions) ([]T, *str.Response, error)

// fetchAllPages fetches a list starting at page, following pages while client.HavePages allows.
func fetchAllPages[T any](client *internal.Client, options *str.Options, page int, fetch pageFetcher[T]) ([]T, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := fetch(&opts)
	if err != nil {
		return nil, err
	}

	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
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
func watchNowError(action string, kind string, options *str.Options, resp *str.Response, err error) error {
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found %s for:%s", kind, options.InternalID)
	}

	if vipErr := cli.HandleVIPResponse(resp, err); vipErr != nil {
		return vipErr
	}

	var forbidden *internal.ForbiddenError
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
