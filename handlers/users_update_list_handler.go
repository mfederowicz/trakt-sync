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
	"github.com/mfederowicz/trakt-sync/writer"
)

// UsersUpdateListHandler struct for handler
type UsersUpdateListHandler struct{ common CommonLogic }

// Handle to handle sync: update_list action
func (m UsersUpdateListHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	list, resp, _ := m.common.FetchUsersList(client, options)

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found list for:%s", options.ID)
	}

	printer.Println("Update personal list by sending 1 or more parameters.")

	result, _, err := m.usersUpdateList(client, list, options)
	if err != nil {
		return fmt.Errorf("update personal list error:%w", err)
	}
	options.Output = "export_users_update_list_results.json"
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (UsersUpdateListHandler) usersUpdateList(client *internal.Client, list *str.PersonalList, options *str.Options) (*str.PersonalList, *str.Response, error) {
	if len(options.Description) > consts.ZeroValue {
		list.Description = &options.Description
	}
	list.Privacy = &options.Privacy
	list.DisplayNumbers = &options.DisplayNumbers
	list.AllowComments = &options.AllowComments
	list.SortBy = &options.SortBy
	list.SortHow = &options.SortHow
	result, resp, err := client.Users.UpdateList(
		client.BuildCtxFromOptions(options),
		&options.UserName,
		&options.ID,
		list)
	if err != nil {
		return nil, nil, fmt.Errorf("update personal list error:%w", err)
	}

	return result, resp, nil
}
