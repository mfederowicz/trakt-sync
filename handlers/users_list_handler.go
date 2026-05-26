// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/uri"
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersListHandler struct for handler
type UsersListHandler struct{ common CommonLogic }

// Handle to handle users: list action
func (u UsersListHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	result, resp, _ := u.fetchSingleList(client, options)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found list for:%s", options.InternalID)
	}

	printer.Printf("Found list for traktId:%s and name:%s \n", options.InternalID, *result.Name)

	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}

func (UsersListHandler) fetchSingleList(client *internal.Client, options *str.Options) (*str.PersonalList, *str.Response, error) {
	opts := uri.ListOptions{Extended: options.ExtendedInfo}
	username := options.UserName
	listID := options.ID

	result, resp, err := client.Users.GetList(
		client.BuildCtxFromOptions(options),
		&username,
		&listID,
		&opts,
	)

	return result, resp, err
}
