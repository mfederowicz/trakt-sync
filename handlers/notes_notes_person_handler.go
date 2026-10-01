// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/trakt"
)

// NotesNotesPersonHandler struct for handler
type NotesNotesPersonHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesPersonHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}
	person, err := h.common.FetchPerson(client, options)
	if err != nil {
		return fmt.Errorf("fetch person error:%w", err)
	}
	n := new(str.Notes)
	n.Person = person
	n.Notes = &options.Notes

	result, resp, err := h.common.Notes(client, n, options)
	if err != nil {
		return fmt.Errorf("notes error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, person notes number:%d \n", result.ID)
	}

	return nil
}
