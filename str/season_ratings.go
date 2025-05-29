// Package str used for structs
package str

// SeasonRatings represents JSON season ratings object
type SeasonRatings struct {
	Rating       *float32        `json:"rating,omitempty"`
	Votes        *int            `json:"votes,omitempty"`
	Distribution *map[string]int `json:"distribution,omitempty"`
}

func (s SeasonRatings) String() string {
	return Stringify(s)
}
