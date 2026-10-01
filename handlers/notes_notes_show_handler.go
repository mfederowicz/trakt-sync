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

// NotesNotesShowHandler struct for handler
type NotesNotesShowHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesShowHandler) Handle(options *str.Options, client *trakt.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}
	show, err := h.common.FetchShow(client, options)
	if err != nil {
		return fmt.Errorf("fetch show error:%w", err)
	}
	n := new(str.Notes)
	n.Show = show
	n.Notes = &options.Notes

	result, resp, err := h.common.Notes(client, n, options)
	if err != nil {
		return fmt.Errorf("notes error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, show notes number:%d \n", result.ID)
	}

	return nil
}
