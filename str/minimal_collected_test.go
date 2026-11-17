package str

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// A minimal collection value is collected_at, or collected_at|source when available_on was sent (trakt/trakt-api#1089);
// it is decoded into its parts and encoded back in the same form.
func TestMinimalCollectedJSON(t *testing.T) {
	collectedAt := &Timestamp{Time: time.Date(2002, time.July, 19, 21, 0, 0, 0, time.UTC)}
	tests := []struct {
		name string
		in   string
		want MinimalCollected
		out  string
	}{
		{name: "plain", in: `"2002-07-19T21:00:00.000Z"`, want: MinimalCollected{CollectedAt: collectedAt}, out: `"2002-07-19T21:00:00Z"`},
		{name: "empty source", in: `"2002-07-19T21:00:00.000Z|"`, want: MinimalCollected{CollectedAt: collectedAt, Source: String("")}, out: `"2002-07-19T21:00:00Z|"`},
		{name: "plex", in: `"2002-07-19T21:00:00.000Z|plex"`, want: MinimalCollected{CollectedAt: collectedAt, Source: String("plex")}, out: `"2002-07-19T21:00:00Z|plex"`},
		{name: "jellyfin", in: `"2002-07-19T21:00:00.000Z|jellyfin"`, want: MinimalCollected{CollectedAt: collectedAt, Source: String("jellyfin")}, out: `"2002-07-19T21:00:00Z|jellyfin"`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := MinimalCollected{}
			assert.NoError(t, json.Unmarshal([]byte(tt.in), &got))
			assert.Equal(t, tt.want, got)
			out, err := json.Marshal(got)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.out, string(out))
		})
	}
}

func TestMinimalCollectedJSONErrors(t *testing.T) {
	for _, in := range []string{`"not a date|plex"`, `"|plex"`, `12`} {
		got := MinimalCollected{}
		assert.Error(t, json.Unmarshal([]byte(in), &got), in)
	}
}

func TestMinimalCollectedJSONNull(t *testing.T) {
	got := MinimalCollection{}
	assert.NoError(t, json.Unmarshal([]byte(`{"28":null}`), &got))
	out, err := json.Marshal(got)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"28":null}`, string(out))

	out, err = json.Marshal(MinimalCollected{})
	assert.NoError(t, err)
	assert.Equal(t, "null", string(out))
}

func TestMinimalShowCollectionWithSource(t *testing.T) {
	in := `{"1390":{"1":{"1":"2002-07-19T21:00:00.000Z|plex","2":"2002-07-19T21:00:00.000Z|"}}}`
	got := MinimalShowCollection{}
	assert.NoError(t, json.Unmarshal([]byte(in), &got))
	assert.Equal(t, "plex", *got["1390"]["1"]["1"].Source)
	assert.Equal(t, "", *got["1390"]["1"]["2"].Source)
	out, err := json.Marshal(got)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"1390":{"1":{"1":"2002-07-19T21:00:00Z|plex","2":"2002-07-19T21:00:00Z|"}}}`, string(out))
}
