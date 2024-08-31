// Package str used for structs
package str

// Character represents JSON character object
type Character struct {
	Character     *string   `json:"character,omitempty"`
	Characters    *[]string `json:"characters,omitempty"`
	EpisodeCount  *int      `json:"episode_count,omitempty"`
	SeriesRegular *bool     `json:"series_regular,omitempty"`
	Movie         *Movie    `json:"movie,omitempty"`
	Show          *Show     `json:"show,omitempty"`
}

func (c Character) String() string {
	return Stringify(c)
}
