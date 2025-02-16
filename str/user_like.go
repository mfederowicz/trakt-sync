// Package str used for structs
package str

// UserLike represents JSON user like object
type UserLike struct {
	LikedAt *Timestamp `json:"liked_at,omitempty"`
	User    *UserProfile    `json:"user,omitempty"`
}

func (u UserLike) String() string {
	return Stringify(u)
}
