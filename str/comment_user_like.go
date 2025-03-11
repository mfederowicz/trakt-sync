// Package str used for structs
package str

// CommentUserLike represents JSON user comment like object
type CommentUserLike struct {
	LikedAt *Timestamp   `json:"liked_at,omitempty"`
	User    *UserProfile `json:"user,omitempty"`
}

func (c CommentUserLike) String() string {
	return Stringify(c)
}
