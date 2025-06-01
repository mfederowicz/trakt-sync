// Package str used for structs
package str

// EpisodePeople represents JSON people connected with episode object
type EpisodePeople struct {
	Cast       *[]Character `json:"cast,omitempty"`
	GuestStars *[]Character `json:"guest_stars,omitempty"`
	Crew       *Crew        `json:"crew,omitempty"`
}

func (m EpisodePeople) String() string {
	return Stringify(m)
}
