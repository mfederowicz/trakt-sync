// Package str used for structs
package str

// CommentReaction represents JSON comment reaction object
type CommentReaction struct {
	ReactedAt *Timestamp   `json:"reacted_at,omitempty"`
	Reaction  *Reaction    `json:"reaction,omitempty"`
	User      *UserProfile `json:"user,omitempty"`
}

func (c CommentReaction) String() string {
	return Stringify(c)
}
