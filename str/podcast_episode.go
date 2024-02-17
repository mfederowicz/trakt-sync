package str

type PodcastEpisode struct {
	Season *int     `json:"season,omitempty"`
	Number *int     `json:"number,omitempty"`
	Title  *string  `json:"title,omitempty"`
	Ids    *Ids `json:"ids,omitempty"`
}

func (p PodcastEpisode) String() string {
	return Stringify(p)
}
