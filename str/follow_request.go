// Package str used for structs
package str

// FollowRequest represents JSON follow request object
type FollowRequest struct {
	ID          *int64       `json:"id,omitempty"`
	RequestedAt *Timestamp   `json:"requested_at,omitempty"`
	User        *UserProfile `json:"user,omitempty"`
}

func (f FollowRequest) String() string {
	return Stringify(f)
}
