// Package str used for structs
package str

// CommentMediaItem represents JSON comment media item object
type CommentMediaItem struct {
	Type    *string       `json:"type,omitempty"`
	Movie   *Movie        `json:"movie,omitempty"`
	Show    *Show         `json:"show,omitempty"`
	Season  *Season       `json:"season,omitempty"`
	Episode *Episode      `json:"episode,omitempty"`
	List    *PersonalList `json:"list,omitempty"`
}

func (c CommentMediaItem) String() string {
	return Stringify(c)
}
