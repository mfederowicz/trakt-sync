// Package handlers used to handle module actions
package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// PeopleMoviesHandler struct for handler
type PeopleMoviesHandler struct{}

// Handle to handle people: movies action
func (p PeopleMoviesHandler) Handle(options *str.Options, client *internal.Client) error {
		if len(options.ID) == consts.ZeroValue {
			return fmt.Errorf(consts.EmptyPersonIDMsg)
		}
		printer.Println("Get movie credits")
		result, err := p.fetchMovieCredits(client, options)
		if err != nil {
			return fmt.Errorf("fetch movie credits error:%v", err)
		}

		if result == nil {
			return fmt.Errorf("empty result")
		}

		printer.Print("Found movie credits data \n")
		print("write data to:" + options.Output)
		jsonData, _ := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
		writer.WriteJSON(options, jsonData)


	return nil
}

func (PeopleMoviesHandler) fetchMovieCredits(client *internal.Client, options *str.Options) (*str.PersonMovies, error) {
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

