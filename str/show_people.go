// Package str used for structs
package str

// ShowPeople represents JSON people connected with show object
type ShowPeople struct {
	Cast       *[]Character `json:"cast,omitempty"`
	GuestStars *[]Character `json:"guest_stars,omitempty"`
	Crew       *Crew        `json:"crew,omitempty"`
}

func (s ShowPeople) String() string {
	return Stringify(s)
}
