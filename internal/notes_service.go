// Package internal used for client and services
package internal

import (
	"context"
	"errors"
	"fmt"
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
		return nil, nil, errors.New("internal server error")
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		return nil, nil, errors.New("validation error")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil, errors.New("item not found or doesn't allow notes")
	}

	if err != nil {
		return nil, nil, errors.Join(resp.Errors.GetComments())
	}

	return note, resp, nil
}

// DeleteNotes Delete a single note.
//
// API docs: https://trakt.docs.apiary.io/#reference/notes/note/delete-a-note
func (n *NotesService) DeleteNotes(ctx context.Context, id *string) (*str.Response, error) {
	var url = fmt.Sprintf("notes/%s", *id)
	printer.Println("delete notes")
	req, err := n.client.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := n.client.Do(ctx, req, nil)
	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("notes not found with Id:%s", *id)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateNotes Update a single note (500 maximum characters).
//
// API docs:https://trakt.docs.apiary.io/#reference/notes/note/update-a-note 
func (n *NotesService) UpdateNotes(ctx context.Context, id *string, notes *str.Notes) (*str.Notes, *str.Response, error) {
	var url = fmt.Sprintf("notes/%s", *id)
	printer.Println("update notes")
	req, err := n.client.NewRequest(http.MethodPut, url, notes)
	if err != nil {
		return nil, nil, err
	}

	note := new(str.Notes)
	resp, err := n.client.Do(ctx, req, note)
	if err != nil {
		return nil, nil, errors.Join(resp.Errors.GetComments())
	}

	return note, resp, nil
}

// GetNotes Return a single note.
//
// API docs:https://trakt.docs.apiary.io/#reference/notes/note/get-a-note 
func (n *NotesService) GetNotes(ctx context.Context, id *string) (*str.Notes, *str.Response, error) {
	var url = fmt.Sprintf("notes/%s", *id)
	printer.Println("fetch notes url:" + url)
	req, err := n.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(str.Notes)
	resp, err := n.client.Do(ctx, req, &result)

	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("notes not found with Id:%s", *id)
	}

	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
