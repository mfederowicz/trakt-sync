// Package str used for structs
package str

// CommentItem represents JSON comment item object
type CommentItem struct {
	Type    *string  `json:"type,omitempty"`
	Movie   *Movie   `json:"movie,omitempty"`
	Season  *Season  `json:"season,omitempty"`
	Episode *Episode `json:"episode,omitempty"`
	Show    *Show    `json:"show,omitempty"`
	List    *List    `json:"list,omitempty"`
	Comment *Comment `json:"comment,omitempty"`
}

func (c CommentItem) String() string {
	return Stringify(c)
}
