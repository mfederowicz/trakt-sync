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

// NotesNotesCollectionHandler struct for handler
type NotesNotesCollectionHandler struct{ common CommonLogic }

// Handle to handle comments: movie type
func (h NotesNotesCollectionHandler) Handle(options *str.Options, client *trakt.Client) error {
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
	case consts.Movie:
		movie, _, err := h.common.FetchMovie(client, options)
		if err != nil {
			return fmt.Errorf("fetch movie error:%w", err)
		}
		n.Movie = movie
	case consts.Episode:
		episode, err := h.common.EpisodeFromTraktID(options)
		if err != nil {
			return err
		}
		n.Episode = episode
	default:
		h.common.GenActionTypeItemUsage(options, []string{consts.Movie, consts.Episode})
		return unknownValueError("item", options.Item)
	}
	p := "private"
	n.Privacy = &p
	result, resp, err := h.common.Notes(client, n, options)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if resp.StatusCode == http.StatusCreated {
		printer.Printf("result: success, collection notes number:%d \n", result.ID)
	}

	return nil
}
