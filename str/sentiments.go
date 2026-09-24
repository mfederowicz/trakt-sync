// Package str used for structs
package str

// Sentiments represents JSON sentiments object: sentiment counts for comments and reactions
type Sentiments struct {
	Good         []*Sentiment `json:"good,omitempty"`
	Bad          []*Sentiment `json:"bad,omitempty"`
	AnalyzedAt   *Timestamp   `json:"analyzed_at,omitempty"`
	CommentCount *int         `json:"comment_count,omitempty"`
}

func (s Sentiments) String() string {
	return Stringify(s)
}
