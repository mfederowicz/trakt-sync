package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// HiddenItem keeps a hidden user of the comments section when decoded and encoded again (shape seen live on 2026-11-14).
func TestHiddenItemKeepsUser(t *testing.T) {
	const in = `{"hidden_at":"2026-10-02T12:18:09Z","type":"user","user":{"username":"sean","private":false,"name":"Sean Rudford","vip":true,"ids":{"slug":"sean"}}}`
	item := &HiddenItem{}
	assert.NoError(t, json.Unmarshal([]byte(in), item))
	out, err := json.Marshal(item)
	assert.NoError(t, err)
	assert.JSONEq(t, in, string(out))
}
