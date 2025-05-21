// Package str used for structs
package str

// ShowRatings represents JSON show ratings object
type ShowRatings struct {
	Rating       *float32        `json:"rating,omitempty"`
	Votes        *int            `json:"votes,omitempty"`
	Distribution *map[string]int `json:"distribution,omitempty"`
}

func (s ShowRatings) String() string {
	return Stringify(s)
}
