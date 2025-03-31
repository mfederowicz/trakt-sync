// Package handlers used to handle module actions
package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// NotesNotesHistoryHandler struct for handler
type NotesNotesHistoryHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesHistoryHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyHistoryIDMsg)
	}
	n := new(str.Notes)
	n.Notes = &options.Notes
	a := new(str.AttachedTo)
	t := "history"
	a.Type = &t
	intID, _ := strconv.Atoi(options.InternalID)
	a.ID = &intID
	n.AttachedTo = a

	result, resp, err := h.common.Notes(client, n)
	if err != nil {
		return fmt.Errorf("notes error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, history notes number:%d \n", result.ID)
	}

	return nil
}
