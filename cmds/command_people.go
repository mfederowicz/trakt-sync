// Package cmds used for commands modules
package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

var (
	_action    = PeopleCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_startDate = PeopleCmd.Flag.String("start_date", "", consts.StartDateUsage)
	_personID  = PeopleCmd.Flag.String("i", cfg.DefaultConfig().ID, consts.UserlistUsage)
)

// PeopleCmd returns all data for selected person.
var PeopleCmd = &Command{
	Name:    "people",
	Usage:   "",
	Summary: "Returns all data for selected person.",
	Help:    `people command`,
}

func peopleFunc(cmd *Command, _ ...string) error {

	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	switch options.Action {
	case "updates":
		fmt.Println("Get recently updated people")
		date := time.Now().Format("2006-01-02T15:00Z")
		updates, err := fetchPeoplesUpdates(client, options, date, consts.DefaultPage)
		if err != nil {
			return fmt.Errorf("fetch peoples updates error:%w", err)
		}

		if len(updates) == consts.ZeroValue {
			return fmt.Errorf("empty updates lists")
		}

		if len(updates) > consts.ZeroValue {
			fmt.Printf("Found %d items \n", len(updates))
			exportJSON := []*str.PersonItem{}
			exportJSON = append(exportJSON, updates...)
			print("write data to:" + options.Output)
			jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")
			writer.WriteJSON(options, jsonData)
		} else {
			fmt.Print("No update items to fetch\n")
		}

	case "updated_ids":
		fmt.Println("Get recently updated people Trakt IDs")
		date := time.Now().Format("2006-01-02T15:00Z")
		updates, err := fetchPeoplesUpdatedIDs(client, options, date, consts.DefaultPage)
		if err != nil {
			return fmt.Errorf("fetch peoples updated ids error:%w", err)
		}

		if len(updates) == consts.ZeroValue {
			return fmt.Errorf("empty updates lists")
		}

		if len(updates) > consts.ZeroValue {
			fmt.Printf("Found %d items \n", len(updates))
			exportJSON := []*int{}
			exportJSON = append(exportJSON, updates...)
			print("write data to:" + options.Output)
			jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")

			writer.WriteJSON(options, jsonData)
		} else {
			fmt.Print("No update items to fetch\n")
		}

	case "summary":
		if len(*_personID) == consts.ZeroValue {
			return fmt.Errorf("set personId ie: -i john-wayne")
		}
		fmt.Println("Get a single person")
		result, err := fetchSinglePerson(client, options)
		if err != nil {
			return fmt.Errorf("fetch single person error:%w", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}

		fmt.Print("Found person \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JsonDataFormat)

		writer.WriteJSON(options, jsonData)

	case "movies":
		if len(*_personID) == consts.ZeroValue {
			return fmt.Errorf("set personId ie: -i john-wayne")
		}
		fmt.Println("Get movie credits")
		result, err := fetchMovieCredits(client, options)
		if err != nil {
			return fmt.Errorf("fetch movie credits error:%v", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}

		fmt.Print("Found movie credits data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JsonDataFormat)
		writer.WriteJSON(options, jsonData)

	case "shows":
		if len(*_personID) == consts.ZeroValue {
			return fmt.Errorf(consts.EmptyPersonIdMsg)
		}
		fmt.Println("Get show credits")
		result, err := fetchShowCredits(client, options)
		if err != nil {
			return fmt.Errorf("fetch show credits error:%w", err)
		}

		if result == nil {
			return fmt.Errorf(consts.EmptyResult)
		}

		fmt.Print("Found show credits data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JsonDataFormat)

		writer.WriteJSON(options, jsonData)

	case "lists":
		if len(*_personID) == consts.ZeroValue {
			return fmt.Errorf(consts.EmptyPersonIdMsg)
		}
		fmt.Println("Get lists containing this person")
		result, err := fetchListsContainingThisPerson(client, options, consts.DefaultPage)
		if err != nil {
			return fmt.Errorf("fetch lists error:%v", err)
		}

		if len(result) == consts.ZeroValue {
			return fmt.Errorf("empty lists")
		}

		fmt.Printf("Found %d result \n", len(result))
		exportJSON := []*str.PersonalList{}
		exportJSON = append(exportJSON, result...)
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(exportJSON, consts.EmptyString, consts.JsonDataFormat)

		writer.WriteJSON(options, jsonData)

	default:
		fmt.Println("possible actions: updates, updated_ids, summary, movies, shows, lists")
	}
	return nil
}

var (
	peopleDumpTemplate = `{{.Head}} {{.Pattern}}{{end}}`
)

func init() {
	PeopleCmd.Run = peopleFunc
}

func fetchPeoplesUpdates(client *internal.Client, options *str.Options, startDate string, page int) ([]*str.PersonItem, error) {

	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.People.GetRecentlyUpdatedPeople(
		context.Background(),
		&startDate,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != consts.EmptyString {

		pagesInt, _ := strconv.Atoi(pages)

		if page != pagesInt && pagesInt > consts.ZeroValue {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + consts.NextPageStep
			nextPageItems, err := fetchPeoplesUpdates(client, options, startDate, nextPage)
			if err != nil {
				return nil, err
			}

			// Append items from the next page to the current page
			list = append(list, nextPageItems...)

		}

	}

	return list, nil

}

func fetchPeoplesUpdatedIDs(client *internal.Client, options *str.Options, startDate string, page int) ([]*int, error) {

	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.People.GetRecentlyUpdatedPeopleTraktIDs(
		context.Background(),
		&startDate,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != consts.EmptyString {

		pagesInt, _ := strconv.Atoi(pages)

		if page != pagesInt && pagesInt > consts.ZeroValue {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + consts.NextPageStep
			nextPageItems, err := fetchPeoplesUpdatedIDs(client, options, startDate, nextPage)
			if err != nil {
				return nil, err
			}

			// Append items from the next page to the current page
			list = append(list, nextPageItems...)

		}

	}

	return list, nil

}

func fetchSinglePerson(client *internal.Client, options *str.Options) (*str.Person, error) {

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, _, err := client.People.GetSinglePerson(
		context.Background(),
		&options.ID,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func fetchMovieCredits(client *internal.Client, options *str.Options) (*str.PersonMovies, error) {

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, _, err := client.People.GetMovieCredits(
		context.Background(),
		&options.ID,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func fetchShowCredits(client *internal.Client, options *str.Options) (*str.PersonShows, error) {

	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, _, err := client.People.GetShowCredits(
		context.Background(),
		&options.ID,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func fetchListsContainingThisPerson(client *internal.Client, options *str.Options, page int) ([]*str.PersonalList, error) {

	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := client.People.GetListsContainingThisPerson(
		context.Background(),
		&options.ID,
		&options.Type,
		&options.Sort,
		&opts,
	)

	if err != nil {
		return nil, err
	}

	// Check if there are more pages
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != consts.EmptyString {

		pagesInt, _ := strconv.Atoi(pages)

		if page != pagesInt && pagesInt > consts.ZeroValue {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + consts.NextPageStep
			nextPageItems, err := fetchListsContainingThisPerson(client, options, nextPage)
			if err != nil {
				return nil, err
			}

			// Append items from the next page to the current page
			list = append(list, nextPageItems...)

		}

	}

	return list, nil

}
