// Package str used for structs
package str

// RatingItems represents JSON list object
type RatingItems struct {
	Movies   *[]Movie   `json:"movies,omitempty"`
	Shows    *[]Show    `json:"shows,omitempty"`
	Seasons  *[]Season  `json:"seasons,omitempty"`
	Episodes *[]Episode `json:"episodes,omitempty"`
}

func (r RatingItems) String() string {
	return Stringify(r)
}
