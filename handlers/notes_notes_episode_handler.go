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

// NotesNotesEpisodeHandler struct for handler
type NotesNotesEpisodeHandler struct{ common CommonLogic }

// Handle to handle notes: episode type
func (h NotesNotesEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}
	episode, _ := h.common.FetchEpisode(client, options)
	n := new(str.Notes)
	n.Episode = episode
	n.Notes = &options.Notes

	result, resp, err := h.common.Notes(client, n, options)
	if err != nil {
		return fmt.Errorf("notes error:%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, episode notes number:%d \n", result.ID)
	}

	return nil
}
