// Package str used for structs
package str

// MovieRatings represents JSON movie ratings object
type MovieRatings struct {
	Rating       *float32        `json:"rating,omitempty"`
	Votes        *int            `json:"votes,omitempty"`
	Distribution *map[string]int `json:"distribution,omitempty"`
}

func (m MovieRatings) String() string {
	return Stringify(m)
}
