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

// NotesNotesPersonHandler struct for handler
type NotesNotesPersonHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesPersonHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}
	person, _ := h.common.FetchPerson(client, options)
	n := new(str.Notes)
	n.Person = person
	n.Notes = &options.Notes

	result, resp, err := h.common.Notes(client, n)
	if err != nil {
		return fmt.Errorf("notes error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, person notes number:%d \n", result.ID)
	}

	return nil
}
