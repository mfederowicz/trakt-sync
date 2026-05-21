// Package handlers used to handle module actions
package handlers

import (
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// UsersListLikeHandler struct for handler
type UsersListLikeHandler struct{ common CommonLogic }

// Handle to handle users: list_like action
func (u UsersListLikeHandler) Handle(options *str.Options, client *internal.Client) error {
	err := u.common.CheckTypes(options)
	if err != nil {
		return err
	}

	if options.Delete {
		return u.removeListLike(client, options)
	}

	return u.sendListLike(client, options)
}

func (u UsersListLikeHandler) removeListLike(client *internal.Client, options *str.Options) error {
	printer.Println("Remove like on a list")
	resp, err := u.common.UsersRemoveListLike(client, options)
	if err != nil {
		return fmt.Errorf("result: error, remove like on list:%s, %w", options.ID, err)
	}
	if resp.StatusCode == http.StatusNoContent {
		printer.Printf("result: success, remove like on list:%s \n", options.ID)
	}
	return nil
}

func (u UsersListLikeHandler) sendListLike(client *internal.Client, options *str.Options) error {
	printer.Println("Like a list")
	resp, err := u.common.UsersListLike(client, options)
	if err != nil {
		return fmt.Errorf("result: error, like on list:%s %w", options.ID, err)
	}
	if resp.StatusCode == http.StatusNoContent {
		printer.Printf("result: success, like on list:%s \n", options.ID)
	}

	return nil
}
