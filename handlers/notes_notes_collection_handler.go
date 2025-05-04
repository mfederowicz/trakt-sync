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

// NotesNotesCollectionHandler struct for handler
type NotesNotesCollectionHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesCollectionHandler) Handle(options *str.Options, client *internal.Client) error {
	if len(options.InternalID) == consts.ZeroValue {
		return errors.New(consts.EmptyMovieIDMsg)
	}
	n := new(str.Notes)
	n.Notes = &options.Notes
	a := new(str.AttachedTo)
	t := "collection"
	a.Type = &t
	n.AttachedTo = a
	switch options.Item {
	case "movie":
		movie, _, _ := h.common.FetchMovie(client, options)
		n.Movie = movie
	case "episode":
		episode, _ := h.common.FetchEpisode(client, options)
		n.Episode = episode
	default:
		h.common.GenActionTypeItemUsage(options, []string{"movie", "episode"})
		return nil
	}
	p := "private"
	n.Privacy = &p
	result, resp, err := h.common.Notes(client, n)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, collection notes number:%d \n", result.ID)
	}

	return nil
}
