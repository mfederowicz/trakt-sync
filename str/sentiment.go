// Package str used for structs
package str

// Sentiment represents JSON sentiment object
type Sentiment struct {
	Sentiment  *string `json:"sentiment,omitempty"`
	CommentIDs *[]int  `json:"comment_ids,omitempty"`
}

func (s Sentiment) String() string {
	return Stringify(s)
}
