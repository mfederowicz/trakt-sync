// Package cmds used for commands modules
package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

var (
	_action    = PeopleCmd.Flag.String("a", cfg.DefaultConfig().Action, ActionUsage)
	_startDate = PeopleCmd.Flag.String("start_date", "", StartDateUsage)
	_personID  = PeopleCmd.Flag.String("i", cfg.DefaultConfig().ID, UserlistUsage)
)

// PeopleCmd returns all data for selected person.
var PeopleCmd = &Command{
	Name:    "people",
	Usage:   "",
	Summary: "Returns all data for selected person.",
	Help:    `people command`,
}

func peopleFunc(cmd *Command, _ ...string) {

	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	switch options.Action {
	case "updates":
		fmt.Println("Get recently updated people")
		date := time.Now().Format("2006-01-02T15:00Z")
		updates, err := fetchPeoplesUpdates(client, options, date, 1)
		if err != nil {
			fmt.Printf("fetch peoples updates error:%v", err)
			os.Exit(0)
		}

		if len(updates) == 0 {
			fmt.Print("empty updates lists")
			os.Exit(0)
		}

		if err == nil {
			if len(updates) > 0 {
				fmt.Printf("Found %d items \n", len(updates))
				exportJSON := []*str.PersonItem{}
				exportJSON = append(exportJSON, updates...)
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")
				writer.WriteJSON(options, jsonData)
			} else {
				fmt.Print("No update items to fetch\n")
			}

		}

	case "updated_ids":
		fmt.Println("Get recently updated people Trakt IDs")
		date := time.Now().Format("2006-01-02T15:00Z")
		updates, err := fetchPeoplesUpdatedIDs(client, options, date, 1)
		if err != nil {
			fmt.Printf("fetch peoples updated ids error:%v", err)
			os.Exit(0)
		}

		if len(updates) == 0 {
			fmt.Print("empty updates lists")
			os.Exit(0)
		}

		if err == nil {
			if len(updates) > 0 {
				fmt.Printf("Found %d items \n", len(updates))
				exportJSON := []*int{}
				exportJSON = append(exportJSON, updates...)
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")

				writer.WriteJSON(options, jsonData)
			} else {
				fmt.Print("No update items to fetch\n")
			}

		}

	case "summary":
		if len(*_personID) == 0 {
			fmt.Print("set personId ie: -i john-wayne")
			os.Exit(0)
		}
		fmt.Println("Get a single person")
		result, err := fetchSinglePerson(client, options)
		if err != nil {
			fmt.Printf("fetch single person error:%v", err)
			os.Exit(0)
		}

		if result == nil {
			fmt.Print("empty result")
			os.Exit(0)
		}

		if err == nil {
			if result != nil {
				fmt.Print("Found person \n")
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(result, "", "  ")

				writer.WriteJSON(options, jsonData)
			} else {
				fmt.Print("No person to fetch\n")
			}

		}

	case "movies":
		if len(*_personID) == 0 {
			fmt.Print("set personId ie: -i john-wayne")
			os.Exit(0)
		}
		fmt.Println("Get movie credits")
		result, err := fetchMovieCredits(client, options)
		if err != nil {
			fmt.Printf("fetch movie credits error:%v", err)
			os.Exit(0)
		}

		if result == nil {
			fmt.Print("empty result")
			os.Exit(0)
		}

		if err == nil {
			if result != nil {
				fmt.Print("Found movie credits data \n")
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(result, "", "  ")

				writer.WriteJSON(options, jsonData)
			} else {
				fmt.Print("No movie credits to fetch\n")
			}

		}

	case "shows":
		if len(*_personID) == 0 {
			fmt.Print("set personId ie: -i john-wayne")
			os.Exit(0)
		}
		fmt.Println("Get show credits")
		result, err := fetchShowCredits(client, options)
		if err != nil {
			fmt.Printf("fetch show credits error:%v", err)
			os.Exit(0)
		}

		if result == nil {
			fmt.Print("empty result")
			os.Exit(0)
		}

		if err == nil {
			if result != nil {
				fmt.Print("Found show credits data \n")
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(result, "", "  ")

				writer.WriteJSON(options, jsonData)
			} else {
				fmt.Print("No show credits to fetch\n")
			}

		}

	case "lists":
		if len(*_personID) == 0 {
			fmt.Print("set personId ie: -i john-wayne")
			os.Exit(0)
		}
		fmt.Println("Get lists containing this person")
		result, err := fetchListsContainingThisPerson(client, options, 1)
		if err != nil {
			fmt.Printf("fetch lists error:%v", err)
			os.Exit(0)
		}

		if len(result) == 0 {
			fmt.Print("empty lists")
			os.Exit(0)
		}

		if err == nil {
			if len(result) > 0 {
				fmt.Printf("Found %d result \n", len(result))
				exportJSON := []*str.PersonalList{}
				exportJSON = append(exportJSON, result...)
				print("write data to:" + options.Output)
				jsonData, _ := json.MarshalIndent(exportJSON, "", "  ")

				writer.WriteJSON(options, jsonData)
			} else {
				fmt.Print("No lists to fetch\n")
			}

		}

	default:
		fmt.Println("possible actions: updates, updated_ids, summary, movies, shows, lists")
	}

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
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != "" {

		pagesInt, _ := strconv.Atoi(pages)

		if page != pagesInt && pagesInt > 0 {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + 1
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
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != "" {

		pagesInt, _ := strconv.Atoi(pages)

		if page != pagesInt && pagesInt > 0 {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + 1
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
	if pages := resp.Header.Get(internal.HeaderPaginationPageCount); pages != "" {

		pagesInt, _ := strconv.Atoi(pages)

		if page != pagesInt && pagesInt > 0 {

			time.Sleep(time.Duration(2) * time.Second)

			// Fetch items from the next page
			nextPage := page + 1
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
