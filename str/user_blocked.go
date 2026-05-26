// Package str used for structs
package str

// UserBlocked represents JSON user blocked object
type UserBlocked struct {
	BlockedAt *Timestamp   `json:"blocked_at,omitempty"`
	User      *UserProfile `json:"user,omitempty"`
}

func (u UserBlocked) String() string {
	return Stringify(u)
}
