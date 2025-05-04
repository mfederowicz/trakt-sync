// Package str used for structs
package str

// PersonShows represents JSON cast and crew object for person
type PersonShows struct {
	Cast *[]Character `json:"cast,omitempty"`
	Crew *Crew        `json:"crew,omitempty"`
}

func (p PersonShows) String() string {
	return Stringify(p)
}
