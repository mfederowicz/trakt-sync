// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
	"github.com/mfederowicz/trakt-sync/writer"
)

// RecommendationsMoviesHandler struct for handler
type RecommendationsMoviesHandler struct{ common CommonLogic }

// Handle to handle recommendations: movies action
func (r RecommendationsMoviesHandler) Handle(options *str.Options, client *trakt.Client) error {
	if options.Hide {
		return r.HandleHide(client, options)
	}

	if err := checkMediaFilters(options); err != nil {
		return err
	}

	result, err := r.common.FetchMovieRecommendations(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("fetch movie recommendations error:%w", err)
	}
	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("encode %s result: %w", options.Action, err)
	}
	printer.Println("write data to:" + options.Output)
	writer.WriteJSON(options, jsonData)

	return nil
}

// HandleHide hide recommendation for movie
func (r RecommendationsMoviesHandler) HandleHide(client *trakt.Client, options *str.Options) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}
	resp, err := r.common.HideMovieRecommendation(client, options)
	if err != nil {
		return fmt.Errorf("hide recommendation error:%w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		printer.Printf("result: success, hide recommendation for:%s \n", options.InternalID)
	}
	return nil
}
