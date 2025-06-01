// Package str used for structs
package str

// Account represents JSON account object
type Account struct {
	SettingsAt  *Timestamp `json:"settings_at,omitempty"`
	FollowedAt  *Timestamp `json:"followed_at,omitempty"`
	FollowingAt *Timestamp `json:"following_at,omitempty"`
	PenndingAt  *Timestamp `json:"pending_at,omitempty"`
	RequestedAt *Timestamp `json:"requested_at,omitempty"`
}

func (f Account) String() string {
	return Stringify(f)
}
