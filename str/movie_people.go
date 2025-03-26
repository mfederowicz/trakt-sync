// Package str used for structs
package str

// MoviePeople represents JSON people connected with movie object
type MoviePeople struct {
	Cast *[]Character `json:"cast,omitempty"`
	Crew *Crew        `json:"crew,omitempty"`
}

func (m MoviePeople) String() string {
	return Stringify(m)
}
