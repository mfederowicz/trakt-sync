package str

type PersonShows struct {
	Cast *[]Character `json:"cast,omitempty"`
	Crew *Crew        `json:"crew,omitempty"`
}

func (p PersonShows) String() string {
	return Stringify(p)
}

