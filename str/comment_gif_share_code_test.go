package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Comment types keep the gif and the settings account keeps share_code and display_ads when decoded and encoded again.
func TestCommentGifAndAccountShareCodeKeepContractFields(t *testing.T) {
	const comment = `{"id":417,"comment":"Great movie!","gif":{"url":"https://static.klipy.com/ii/minions.gif","slug":"minions-cheer"},"spoiler":false}`
	const gifOnly = `{"id":418,"gif":{"url":"https://static.klipy.com/ii/minions.gif"}}`

	tests := []struct {
		name string
		in   string
		v    any
	}{
		{name: "comment", in: comment, v: &Comment{}},
		{name: "comment gif without slug", in: gifOnly, v: &Comment{}},
		{name: "list comment", in: comment, v: &ListComment{}},
		{name: "comment item", in: `{"type":"movie","comment":` + comment + `}`, v: &CommentItem{}},
		{name: "settings account", in: `{"account":{"timezone":"Europe/Warsaw","time_24hr":true,"share_code":"abc123","display_ads":false}}`, v: &UserSettings{}},
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
