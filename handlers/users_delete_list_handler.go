// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersDeleteListHandler struct for handler
type UsersDeleteListHandler struct{ common CommonLogic }

// Handle to handle sync: delete_list action
func (m UsersDeleteListHandler) Handle(options *str.Options, client *internal.Client) error {
	err := m.common.CheckTypes(options)
	if err != nil {
		return err
	}

	if len(options.ID) == consts.ZeroValue {
		return errors.New(consts.EmptyInternalIDMsg)
	}

	printer.Println("Delete personal list:", options.ID)

	resp, err := m.usersDeleteList(client, options)

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("personal list not found")
	}

	if err != nil {
		return fmt.Errorf("delete personal list error:%w", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		printer.Printf("result: success, deleted personal list:%s \n", options.ID)
	}

	return nil
}

func (UsersDeleteListHandler) usersDeleteList(client *internal.Client, options *str.Options) (*str.Response, error) {
	resp, err := client.Users.DeleteList(client.BuildCtxFromOptions(options), &options.UserName, &options.ID)

	if resp.StatusCode == http.StatusNotFound {
		return resp, nil
	}

	if err != nil {
		return nil, fmt.Errorf("delete personal list error:%w", err)
	}

	return resp, nil
}
