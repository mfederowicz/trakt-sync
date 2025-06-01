// Package str used for structs
package str

// EpisodeRatings represents JSON episode ratings object
type EpisodeRatings struct {
	Rating       *float32        `json:"rating,omitempty"`
	Votes        *int            `json:"votes,omitempty"`
	Distribution *map[string]int `json:"distribution,omitempty"`
}

func (s EpisodeRatings) String() string {
	return Stringify(s)
}
