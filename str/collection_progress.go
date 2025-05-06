// Package str used for structs
package str

// CollectionProgress represents JSON show_collection_progress object
type CollectionProgress struct {
	Aired           *int       `json:"aired,omitempty"`
	Completed       *int       `json:"completed,omitempty"`
	LastCollectedAt *Timestamp `json:"last_collected_at,omitempty"`
	Seasons         []*Season  `json:"seasons,omitempty"`
	HiddenSeasons   []*Season  `json:"hidden_seasons,omitempty"`
	NextEpisode     *Episode   `json:"next_episode,omitempty"`
	LastEpisode     *Episode   `json:"last_episode,omitempty"`
}

func (s CollectionProgress) String() string {
	return Stringify(s)
}
