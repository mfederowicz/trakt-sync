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
	_notesNotesID    = NotesCmd.Flag.Int("notes_id", cfg.DefaultConfig().NotesID, consts.NotesIDUsage)
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

	err := cmd.common.ValidPrivacy(options)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	var handler handlers.NotesHandler
	var notesHandlers = map[string]handlers.Handler{
		"notes": handlers.NotesNotesHandler{},
		"note":  handlers.NotesNoteHandler{},
		"item":  handlers.NotesItemHandler{},
	}
	handler, err = cmd.common.GetHandlerForMap(options.Action, notesHandlers)

	if err != nil {
		cmd.common.GenActionsUsage(cmd.Name, []string{"notes", "note", "item"})
		return nil
	}

	err = handler.Handle(options, client)
	if err != nil {
		return fmt.Errorf(cmd.Name+"/"+options.Action+":%s", err)
	}

	return nil
}

var (
	notesDumpTemplate = ``
)

func init() {
	NotesCmd.Run = notesFunc
}
