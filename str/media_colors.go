package str

// MediaColors represents JSON colors object of a movie or a show (extended info "colors")
type MediaColors struct {
	Poster *[]string `json:"poster,omitempty"`
}

func (m MediaColors) String() string {
	return Stringify(m)
}
