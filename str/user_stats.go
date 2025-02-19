// Package str used for structs
package str

// UserStats represents JSON user stats object
type UserStats struct {
	Movies         *Movies   `json:"movies,omitempty"`
	Shows          *Shows    `json:"shows,omitempty"`
	Seasons        *Seasons  `json:"seasons,omitempty"`
	Episodes       *Episodes `json:"episodes,omitempty"`
	Network        *Network  `json:"network,omitempty"`
	Ratings        *Ratings  `json:"ratings,omitempty"`
	Rating         *int      `json:"rating,omitempty"`
	PlayCount      *int      `json:"play_count,omitempty"`
	CompletedCount *int      `json:"completed_count,omitempty"`
}

func (u UserStats) String() string {
	return Stringify(u)
}
