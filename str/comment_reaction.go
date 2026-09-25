// Package str used for structs
package str

// CommentReaction represents JSON comment reaction object; users/reactions/comments sends type and comment instead of user
type CommentReaction struct {
	ReactedAt *Timestamp   `json:"reacted_at,omitempty"`
	Reaction  *Reaction    `json:"reaction,omitempty"`
	User      *UserProfile `json:"user,omitempty"`
	Type      *string      `json:"type,omitempty"`
	Comment   *Comment     `json:"comment,omitempty"`
}

func (c CommentReaction) String() string {
	return Stringify(c)
}
