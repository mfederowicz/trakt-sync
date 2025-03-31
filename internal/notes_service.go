// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"net/http"

	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesService  handles communication with the notes related
// methods of the Trakt API.
type NotesService Service

// AddNotes Add a new notes to a movie, show, season, episode, or person.
//
// API docs:https://trakt.docs.apiary.io/#reference/notes/notes/add-notes
func (n *NotesService) AddNotes(ctx context.Context, notes *str.Notes) (*str.Notes, *str.Response, error) {
	var url = "notes"
	printer.Println("create new notes")
	req, err := n.client.NewRequest(http.MethodPost, url, notes)
	if err != nil {
		return nil, nil, err
	}

	note := new(str.Notes)
	resp, err := n.client.Do(ctx, req, note)

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, nil, errors.New("500 Internal server error")
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		return nil, nil, errors.New("422 validation error")
	}

	if err != nil {
		return nil, nil, errors.Join(resp.Errors.GetComments())
	}

	return note, resp, nil
}
