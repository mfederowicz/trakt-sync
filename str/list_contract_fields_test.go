package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PersonalList and WatchNowSource keep the fields added to the contract on 2026-10-13 when decoded and encoded again.
func TestListAndWatchNowSourceKeepContractFields(t *testing.T) {
	tests := []struct {
		name string
		in   string
		v    any
	}{
		{name: "list stale origin", in: `{"name":"Imported","stale_origin":{"name":"Some Sync Tool","url":"https://example.com"}}`, v: &PersonalList{}},
		{name: "list stale origin without url", in: `{"name":"Imported","stale_origin":{"name":"Some Sync Tool"}}`, v: &PersonalList{}},
		{name: "list converted to smart list", in: `{"name":"Imported","converted_to":{"ids":{"trakt":42,"slug":"imported"}}}`, v: &PersonalList{}},
		{name: "watch now source type", in: `{"source":"netflix","name":"Netflix","type":"subscription"}`, v: &WatchNowSource{}},
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

// The API sends null for a list that is not stale and was not converted; the fields stay out of the export.
func TestPersonalListNullStaleOriginAndConvertedTo(t *testing.T) {
	list := &PersonalList{}
	assert.NoError(t, json.Unmarshal([]byte(`{"name":"Mine","stale_origin":null,"converted_to":null}`), list))
	assert.Nil(t, list.StaleOrigin)
	assert.Nil(t, list.ConvertedTo)
	out, err := json.Marshal(list)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"Mine"}`, string(out))
}

func TestListStaleOriginString(t *testing.T) {
	assert.Equal(t, `str.ListStaleOrigin{Name:"Some Sync Tool"}`, ListStaleOrigin{Name: String("Some Sync Tool")}.String())
}
