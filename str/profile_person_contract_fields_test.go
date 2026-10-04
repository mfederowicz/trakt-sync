package str

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// UserProfile and Person keep the newer contract fields when decoded and encoded again.
func TestProfileAndPersonKeepContractFields(t *testing.T) {
	const profile = `{"email":"sean@example.com","username":"sean","vip_veteran_since":"2019-03-01T14:30:00Z","vip_veteran_years":7,"vip_veteran_tier":7,` +
		`"vip_veteran_title":"veteran","vip_grace_ends_at":"2027-04-01T14:30:00Z"}`
	const person = `{"name":"Agata Kulesza","ids":{"trakt":1},"height":179.07}`

	tests := []struct {
		name string
		in   string
		v    any
	}{
		{name: "profile", in: profile, v: &UserProfile{}},
		{name: "settings", in: `{"user":` + profile + `}`, v: &UserSettings{}},
		{name: "person", in: person, v: &Person{}},
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
