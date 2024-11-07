// Package str used for structs
package str

// UserStats represents JSON user stats object
type UserStats struct {
	Movies   *Movies   `json:"movies,omitempty"`
	Shows    *Shows    `json:"shows,omitempty"`
	Seasons  *Seasons  `json:"seasons,omitempty"`
	Episodes *Episodes `json:"episodes,omitempty"`
	Network  *Network  `json:"network,omitempty"`
	Ratings  *Ratings  `json:"ratings,omitempty"`
}

func (u UserStats) String() string {
	return Stringify(u)
}
