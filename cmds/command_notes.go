// Package cmds used for commands modules
package cmds

import (
	"fmt"

	"github.com/mfederowicz/trakt-sync/cfg"
	"github.com/mfederowicz/trakt-sync/consts"
	"github.com/mfederowicz/trakt-sync/handlers"
)

var (
	_notesAction     = NotesCmd.Flag.String("a", cfg.DefaultConfig().Action, consts.ActionUsage)
	_notesInternalID = NotesCmd.Flag.String("i", cfg.DefaultConfig().InternalID, consts.TraktIDUsage)
	_notesItem       = NotesCmd.Flag.String("item", cfg.DefaultConfig().Item, consts.ItemUsage)
	_notesNotes      = NotesCmd.Flag.String("notes", cfg.DefaultConfig().Notes, consts.NotesUsage)
	_notesDelete     = NotesCmd.Flag.Bool("delete", cfg.DefaultConfig().Delete, consts.DeleteUsage)
	_notesSpoiler    = NotesCmd.Flag.Bool("spoiler", cfg.DefaultConfig().Spoiler, consts.SpoilerUsage)
	_notesPrivacy    = NotesCmd.Flag.String("privacy", cfg.DefaultConfig().Privacy, consts.PrivacyUsage)
)

// NotesCmd manage notes.
var NotesCmd = &Command{
	Name:    "notes",
	Usage:   "",
	Summary: "Manage notes created by user",
	Help:    `notes command`,
}

func notesFunc(cmd *Command, _ ...string) error {
	options := cmd.Options
	client := cmd.Client
	options = cmd.UpdateOptionsWithCommandFlags(options)

	var handler handlers.NotesHandler
	var notesHandlers = map[string]handlers.Handler{
		consts.Notes: handlers.NotesNotesHandler{},
		consts.Note:  handlers.NotesNoteHandler{},
		consts.Item:  handlers.NotesItemHandler{},
	}
	handler, err := cmd.common.GetHandlerForMap(options.Action, notesHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{consts.Notes, consts.Note, consts.Item})
		return unknownActionError(cmd.Name, options.Action)
	}

	// privacy is checked per action, so only after the action is known
	err = cmd.common.ValidPrivacy(options)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf("%s/%s: %w", cmd.Name, options.Action, err)
	}

	return nil
}

func init() {
	NotesCmd.Run = notesFunc
}
