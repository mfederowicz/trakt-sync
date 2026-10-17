package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Collected items keep the contract field available_on (extended available_on) when decoded and encoded again.
func TestCollectionKeepsAvailableOn(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{name: "movie", in: `{"type":"movie","collected_at":"2026-09-01T10:20:30Z","updated_at":"2026-09-02T10:20:30Z",` +
			`"movie":{"title":"Arrival","ids":{"trakt":12601}},"available_on":[{"name":"plex"}]}`},
		{name: "episode", in: `{"type":"episode","collected_at":"2026-09-01T10:20:30Z","updated_at":"2026-09-02T10:20:30Z",` +
			`"episode":{"season":1,"number":2,"ids":{"trakt":73640}},"show":{"title":"The Sopranos","ids":{"trakt":1390}},"available_on":[{"name":"digital"}]}`},
		{name: "show season episodes", in: `{"type":"show","last_collected_at":"2026-09-01T10:20:30Z","show":{"title":"The Sopranos","ids":{"trakt":1390}},` +
			`"seasons":[{"number":1,"episodes":[{"number":1,"collected_at":"2026-09-01T10:20:30Z","available_on":[{"name":"plex"},{"name":"digital"}]},` +
			`{"number":2,"collected_at":"2026-09-01T10:20:30Z"}]}]}`},
		{name: "without extended", in: `{"type":"movie","collected_at":"2026-09-01T10:20:30Z","movie":{"title":"Arrival","ids":{"trakt":12601}}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &ExportlistItem{}
			assert.NoError(t, json.Unmarshal([]byte(tt.in), item))
			out, err := json.Marshal(item)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.in, string(out))
		})
	}
}

// The API sends an empty available_on on every collected item, also without extended available_on; exports leave it out.
func TestCollectionDropsEmptyAvailableOn(t *testing.T) {
	in := `{"type":"episode","episode":{"number":2,"available_on":[]},"available_on":[]}`
	item := &ExportlistItem{}
	assert.NoError(t, json.Unmarshal([]byte(in), item))
	out, err := json.Marshal(item)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"type":"episode","episode":{"number":2}}`, string(out))
}
