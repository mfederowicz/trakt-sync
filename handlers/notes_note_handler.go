// Package handlers used to handle module actions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
	"github.com/mfederowicz/trakt-sync/writer"
)

// NotesNoteHandler struct for handler
type NotesNoteHandler struct{ common CommonLogic }

// Handle to handle checkin: checkin action
func (n NotesNoteHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyNotesIDMsg)
	}

	if options.Delete {
		resp, err := n.common.DeleteNotes(client, options)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if resp.StatusCode == http.StatusNoContent {
			printer.Printf("result: success, remove notes:%s \n", options.InternalID)
		}
		return nil
	}

	if len(options.Notes) > consts.ZeroValue {
		note := new(str.Notes)
		note.Notes = &options.Notes
		note.Spoiler = &options.Spoiler
		note.Privacy = &options.Privacy
		result, resp, err := n.common.UpdateNotes(client, options, note)
		if err != nil {
			return fmt.Errorf("update notes error:%w", err)
		}

		if resp.StatusCode == http.StatusOK {
			printer.Printf("result: success, update notes:%d \n", result.ID)
		}

		return nil
	}

	result, err := n.common.FetchNotes(client, options)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	print("write data to:" + options.Output)
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	writer.WriteJSON(options, jsonData)
	return nil
}
