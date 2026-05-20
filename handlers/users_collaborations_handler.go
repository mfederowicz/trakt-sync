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
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersCollaborationsHandler struct for handler
type UsersCollaborationsHandler struct{ common CommonLogic }

// Handle to handle users: collaborations action
func (u UsersCollaborationsHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("Get all lists a user can collaborate on")
	items, err := u.fetchCollaborations(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get collaborations error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersCollaborationsHandler) fetchCollaborations(client *internal.Client, options *str.Options, page int) ([]*str.PersonalList, error) {
	collaborations, err := u.common.FetchUsersCollaborations(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch collaborations error:%w", err)
	}

	if len(collaborations) == consts.ZeroValue {
		return nil, errors.New("empty collaborations")
	}

	return collaborations, nil
}
