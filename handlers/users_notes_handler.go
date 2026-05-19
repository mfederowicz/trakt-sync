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

// UsersNotesHandler struct for handler
type UsersNotesHandler struct{ common CommonLogic }

// Handle to handle users: notes action
func (u UsersNotesHandler) Handle(options *str.Options, client *internal.Client) error {
	err := u.common.CheckTypes(options)
	if err != nil {
		return err
	}

	if options.Type != "" {
		printer.Println("Returns all items in a user's notes filtered by type:", options.Type)
	} else {
		printer.Println("Returns all items in a user's notes")
	}

	items, err := u.fetchNotes(client, options, consts.DefaultPage)
	if err != nil {
		return fmt.Errorf("get notes error:%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(items, "", "  ")
	writer.WriteJSON(options, jsonData)

	return nil
}

func (u UsersNotesHandler) fetchNotes(client *internal.Client, options *str.Options, page int) ([]*str.NotesItem, error) {
	notes, err := u.common.FetchUsersNotes(client, options, page)
	if err != nil {
		return nil, fmt.Errorf("fetch notes error:%w", err)
	}

	if len(notes) == consts.ZeroValue {
		return nil, errors.New("empty notes")
	}

	return notes, nil
}
