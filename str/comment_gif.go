package str

// CommentGif represents JSON gif object attached to a comment
type CommentGif struct {
	URL  *string `json:"url,omitempty"`
	Slug *string `json:"slug,omitempty"`
}

func (c CommentGif) String() string {
	return Stringify(c)
}
