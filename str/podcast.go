package str

type Podcast struct {
	Title *string `json:"title,omitempty"`
	Year  *int    `json:"year,omitempty"`
	Ids   *Ids    `json:"ids,omitempty"`
}

func (p Podcast) String() string {
	return Stringify(p)
}
