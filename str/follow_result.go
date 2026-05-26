// Package str used for structs
package str

// FollowResult represents JSON follow result object
type FollowResult struct {
	ApprovedAt *Timestamp   `json:"approved_at,omitempty"`
	User       *UserProfile `json:"user,omitempty"`
}

func (f FollowResult) String() string {
	return Stringify(f)
}
