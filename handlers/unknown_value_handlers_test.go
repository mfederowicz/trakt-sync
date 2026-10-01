// Package handlers used to handle module actions
package handlers

import (
	"testing"

	"github.com/mfederowicz/trakt-sync/str"
	"github.com/stretchr/testify/assert"
)

func TestUnknownValueError(t *testing.T) {
	assert.EqualError(t, unknownValueError("type", "bogus"), `unknown type "bogus"`)
	assert.EqualError(t, unknownValueError("item", ""), "no item given")
}

// TestUnknownTypeOrItemFails checks the notes and scrobble handlers return an error (exit status 1)
// after the usage for a wrong -t or -item, before any request is made.
func TestUnknownTypeOrItemFails(t *testing.T) {
	tests := []struct {
		name    string
		handler Handler
		options *str.Options
		want    string
	}{
		{name: "notes notes", handler: NotesNotesHandler{}, options: &str.Options{Module: "notes", Action: "notes", Type: "bogus"}, want: `unknown type "bogus"`},
		{name: "notes collection", handler: NotesNotesCollectionHandler{}, options: &str.Options{Module: "notes", Action: "notes", Type: "collection", InternalID: "1", Item: "bogus"}, want: `unknown item "bogus"`},
		{name: "notes rating", handler: NotesNotesRatingHandler{}, options: &str.Options{Module: "notes", Action: "notes", Type: "rating", InternalID: "1", Item: ""}, want: "no item given"},
		{name: "scrobble start", handler: ScrobbleStartHandler{}, options: &str.Options{Module: "scrobble", Action: "start", Type: "bogus"}, want: `unknown type "bogus"`},
		{name: "scrobble pause", handler: ScrobblePauseHandler{}, options: &str.Options{Module: "scrobble", Action: "pause", Type: "bogus"}, want: `unknown type "bogus"`},
		{name: "scrobble stop", handler: ScrobbleStopHandler{}, options: &str.Options{Module: "scrobble", Action: "stop", Type: "bogus"}, want: `unknown type "bogus"`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// a nil client proves no request is made before the check
			assert.EqualError(t, tt.handler.Handle(tt.options, nil), tt.want)
		})
	}
}
