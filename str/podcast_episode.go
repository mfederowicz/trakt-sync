// Package str used for structs
package str

// PodcastEpisode represents JSON podcast episode object
type PodcastEpisode struct {
	Season *int    `json:"season,omitempty"`
	Number *int    `json:"number,omitempty"`
	Title  *string `json:"title,omitempty"`
	IDs    *IDs    `json:"ids,omitempty"`
}

func (p PodcastEpisode) String() string {
	return Stringify(p)
}
