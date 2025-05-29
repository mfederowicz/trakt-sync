// Package str used for structs
package str

// SeasonPeople represents JSON people connected with season object
type SeasonPeople struct {
	Cast       *[]Character `json:"cast,omitempty"`
	GuestStars *[]Character `json:"guest_stars,omitempty"`
	Crew       *Crew        `json:"crew,omitempty"`
}

func (s SeasonPeople) String() string {
	return Stringify(s)
}
