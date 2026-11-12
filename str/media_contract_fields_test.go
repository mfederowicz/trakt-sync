package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Movie, Show, Recommendation and Media keep the extended fields of the contract when decoded and encoded again.
func TestMediaTypesKeepContractFields(t *testing.T) {
	const common = `"title":"Ida","ids":{"trakt":1},"subgenres":["road-trip"],"original_title":"Ida",` +
		`"images":{"fanart":["f.webp"],"poster":["p.webp"],"logo":["l.webp"],"clearart":["c.webp"],"banner":["b.webp"],"thumb":["t.webp"]},` +
		`"colors":{"poster":["#111111","#222222"]},"social_ids":{"twitter":"ida","facebook":"idafilm","instagram":"ida","wikipedia":"Ida_(film)"}`
	const movie = `{` + common + `,"after_credits":true,"during_credits":false}`
	const show = `{` + common + `,"last_aired":"2026-09-30T20:00:00Z","total_runtime":480}`
	const both = `{` + common + `,"after_credits":true,"during_credits":false,"last_aired":"2026-09-30T20:00:00Z","total_runtime":480}`

	tests := []struct {
		name string
		in   string
		v    any
	}{
		{name: "movie", in: movie, v: &Movie{}},
		{name: "show", in: show, v: &Show{}},
		{name: "recommendation", in: both, v: &Recommendation{}},
		{name: "media", in: both, v: &Media{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, json.Unmarshal([]byte(tt.in), tt.v))
			out, err := json.Marshal(tt.v)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.in, string(out))
		})
	}
}

func TestMediaImagesAndColorsString(t *testing.T) {
	images := MediaImages{Poster: &[]string{"p.webp"}}
	colors := MediaColors{Poster: &[]string{"#111111"}}
	assert.Contains(t, images.String(), "p.webp")
	assert.Contains(t, colors.String(), "#111111")
}

// Movie, Show, Season and Episode send their notes under the contract key "notes".
func TestMediaTypesSendNotes(t *testing.T) {
	notes := "rewatch"
	tests := []struct {
		name string
		v    any
	}{
		{name: "movie", v: Movie{Notes: &notes}},
		{name: "show", v: Show{Notes: &notes}},
		{name: "season", v: Season{Notes: &notes}},
		{name: "episode", v: Episode{Notes: &notes}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := json.Marshal(tt.v)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"notes":"rewatch"}`, string(out))
		})
	}
}
