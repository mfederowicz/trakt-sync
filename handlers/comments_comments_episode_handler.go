// Package handlers used to handle module actions
package handlers

import (
	"errors"

	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/internal"
	"github.com/mfederowicz/trakt-sync/printer"
	"github.com/mfederowicz/trakt-sync/str"
)

// CommentsCommentsEpisodeHandler struct for handler
type CommentsCommentsEpisodeHandler struct{ common CommonLogic }

// Handle to handle comments: episode type
func (h CommentsCommentsEpisodeHandler) Handle(options *str.Options, client *internal.Client) error {
	printer.Println("generate comment:", options.Type)

	if options.TraktID == consts.ZeroValue {
		return errors.New(consts.EmptyTraktIDMsg)
	}	

	return nil
}
