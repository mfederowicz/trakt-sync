package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Show keeps the contract field ending when decoded and encoded again.
func TestShowKeepsEnding(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{name: "concluded", in: `{"title":"The Sopranos","ids":{"trakt":1390},"status":"ended","ending":"concluded"}`},
		{name: "open", in: `{"title":"Firefly","ids":{"trakt":1427},"status":"canceled","ending":"open"}`},
		{name: "airing", in: `{"title":"Reacher","ids":{"trakt":160737},"status":"returning series"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			show := &Show{}
			assert.NoError(t, json.Unmarshal([]byte(tt.in), show))
			out, err := json.Marshal(show)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.in, string(out))
		})
	}
}
