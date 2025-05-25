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

// NotesNotesMovieHandler struct for handler
type NotesNotesMovieHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesMovieHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	movie, _, _ := h.common.FetchMovie(client, options)
	n := new(str.Notes)
	n.Movie = movie
	n.Notes = &options.Notes

	result, resp, err := h.common.Notes(client, n, options)
	if err != nil {
		return fmt.Errorf("notes error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, movie notes number:%d \n", result.ID)
	}

	return nil
}
