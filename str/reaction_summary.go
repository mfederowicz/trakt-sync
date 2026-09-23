// Package str used for structs
package str

// ReactionSummary represents JSON reaction summary object
type ReactionSummary struct {
	ReactionCount *int           `json:"reaction_count,omitempty"`
	UserCount     *int           `json:"user_count,omitempty"`
	Distribution  map[string]int `json:"distribution,omitempty"`
}

func (r ReactionSummary) String() string {
	return Stringify(r)
}
