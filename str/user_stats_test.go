package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// UserStats keeps the list comments count next to the lists count (contract userStatsResponseSchema, trakt/trakt-api#1096).
func TestUserStatsKeepsListComments(t *testing.T) {
	const in = `{"lists":4,"list_comments":7,"total_minutes":1200,"total_plays":30}`
	stats := &UserStats{}
	assert.NoError(t, json.Unmarshal([]byte(in), stats))
	out, err := json.Marshal(stats)
	assert.NoError(t, err)
	assert.JSONEq(t, in, string(out))
}
