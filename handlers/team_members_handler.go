// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// TeamMembersHandler struct for handler
type TeamMembersHandler struct{}

// Handle to handle team: members action
func (TeamMembersHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Returns Trakt team members.")
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	result, _, err := client.Team.GetTeamMembers(client.BuildCtxFromOptions(options), &opts)
	if err != nil {
		return fmt.Errorf("fetch team %s error: %w", options.Action, err)
	}

	if len(result) == consts.ZeroValue {
		return errors.New(consts.EmptyResult)
	}

	printer.Printf("Found %d result \n", len(result))
	printer.Println("write data to:" + options.Output)
	jsonData, err := json.MarshalIndent(result, consts.EmptyString, consts.JSONDataFormat)
	if err != nil {
		return fmt.Errorf("marshal team %s error: %w", options.Action, err)
	}

	writer.WriteJSON(options, jsonData)
	return nil
}
