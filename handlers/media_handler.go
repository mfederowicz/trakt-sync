// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// MediaHandler interface to handle media module action
type MediaHandler interface {
	Handle(options *str.Options, client *internal.Client) error
}

// mediaPageFetcher fetches one page of a media list.
type mediaPageFetcher[T any] func(opts *uri.ListOptions) ([]T, *str.Response, error)

// fetchMediaPages fetches a media list starting at page, following pages while client.HavePages allows.
func fetchMediaPages[T any](client *internal.Client, options *str.Options, page int, fetch mediaPageFetcher[T]) ([]T, error) {
	opts := uri.ListOptions{Page: page, Limit: options.PerPage, Extended: options.ExtendedInfo}
	list, resp, err := fetch(&opts)
	if err != nil {
		return nil, err
	}

	if client.HavePages(page, resp, options.PagesLimit) {
		time.Sleep(time.Duration(consts.SleepNumberOfSeconds) * time.Second)
		nextPageItems, err := fetchMediaPages(client, options, page+consts.NextPageStep, fetch)
		if err != nil {
			return nil, err
		}
		list = append(list, nextPageItems...)
	}

	return list, nil
}

// exportMedia fetches all pages of a media list and writes them to the output file.
func exportMedia[T any](client *internal.Client, options *str.Options, fetch mediaPageFetcher[T]) error {
	result, err := fetchMediaPages(client, options, consts.DefaultPage, fetch)
	if err != nil {
		return fmt.Errorf("fetch media %s error: %w", options.Action, err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal media %s error: %w", options.Action, err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
