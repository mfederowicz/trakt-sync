package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AttachedTo keeps the date and rating of the history, collection or rating entry a note is attached to (contract userNoteResponseSchema, 2026-11-19).
func TestAttachedToKeepsEntryFields(t *testing.T) {
	tests := map[string]string{
		"history":    `{"type":"history","id":1982347,"watched_at":"2026-01-02T20:30:00Z"}`,
		"collection": `{"type":"collection","collected_at":"2026-01-03T10:00:00Z"}`,
		"rating":     `{"type":"rating","rated_at":"2026-01-04T08:15:00Z","rating":8}`,
	}
	for name, in := range tests {
		in := in
		t.Run(name, func(t *testing.T) {
			attached := &AttachedTo{}
			assert.NoError(t, json.Unmarshal([]byte(in), attached))
			out, err := json.Marshal(attached)
			assert.NoError(t, err)
			assert.JSONEq(t, in, string(out))
		})
	}
}
