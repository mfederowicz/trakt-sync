// Package str used for structs
package str

// Friend represents JSON friend object
type Friend struct {
	FriendsAt *Timestamp   `json:"friends_at,omitempty"`
	User      *UserProfile `json:"user,omitempty"`
}

func (f Friend) String() string {
	return Stringify(f)
}
