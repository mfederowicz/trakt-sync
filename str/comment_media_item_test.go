package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CommentMediaItem keeps a person, comment or user the comment is attached to (contract commentItemResponseSchema, 2026-11-19).
func TestCommentMediaItemKeepsPersonCommentUser(t *testing.T) {
	tests := map[string]string{
		"person":  `{"type":"person","person":{"name":"Bryan Cranston","ids":{"trakt":297737,"slug":"bryan-cranston"}}}`,
		"comment": `{"type":"comment","comment":{"id":267,"comment":"Great show!","spoiler":false,"replies":1,"likes":2}}`,
		"user":    `{"type":"user","user":{"username":"sean","private":false,"name":"Sean Rudford","vip":true,"ids":{"slug":"sean"}}}`,
	}
	for name, in := range tests {
		in := in
		t.Run(name, func(t *testing.T) {
			item := &CommentMediaItem{}
			assert.NoError(t, json.Unmarshal([]byte(in), item))
			out, err := json.Marshal(item)
			assert.NoError(t, err)
			assert.JSONEq(t, in, string(out))
		})
	}
}
