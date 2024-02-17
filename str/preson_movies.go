package str

type PersonMovies struct {
	Cast *[]Character `json:"cast,omitempty"`
	Crew *Crew        `json:"crew,omitempty"`
}

func (p PersonMovies) String() string {
	return Stringify(p)
}
