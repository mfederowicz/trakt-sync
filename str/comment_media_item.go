package str

// CommentMediaItem represents JSON comment media item object
type CommentMediaItem struct {
	Type    *string       `json:"type,omitempty"`
	Movie   *Movie        `json:"movie,omitempty"`
	Show    *Show         `json:"show,omitempty"`
	Season  *Season       `json:"season,omitempty"`
	Episode *Episode      `json:"episode,omitempty"`
	Person  *Person       `json:"person,omitempty"`
	List    *PersonalList `json:"list,omitempty"`
	Comment *Comment      `json:"comment,omitempty"`
	User    *UserProfile  `json:"user,omitempty"`
}

func (c CommentMediaItem) String() string {
	return Stringify(c)
}
