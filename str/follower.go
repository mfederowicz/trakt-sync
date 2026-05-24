// Package str used for structs
package str

// Follower represents JSON follower object
type Follower struct {
	FollowedAt *Timestamp   `json:"followed_at,omitempty"`
	User       *UserProfile `json:"user,omitempty"`
}

func (f Follower) String() string {
	return Stringify(f)
}
