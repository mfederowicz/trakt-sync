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

// UsersCollectionHandler struct for handler
type UsersCollectionHandler struct{ common CommonLogic }

// Handle to handle users: collection action
func (u UsersCollectionHandler) Handle(options *str.Options, client *internal.Client) error {
	err := u.common.CheckTypes(options)
	if err != nil {
		return err
	}

	if options.Type != "" {
		printer.Println("Returns all items in a user's collection filtered by type:", options.Type)
	} else {
		printer.Println("Returns all items in a user's collection")
	}

	items, err := u.fetchCollection(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get collection error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersCollectionHandler) fetchCollection(client *internal.Client, options *str.Options, page int) ([]*str.ExportlistItem, error) {
	items, err := u.common.FetchUsersCollection(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch collection error:%w", err)
	}

	if len(items) == consts.ZeroValue {
		return nil, errors.New("empty collection")
	}

	return items, nil
}
