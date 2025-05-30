// Package str used for structs
package str

// Season represents JSON season object
type Season struct {
	Number        *int       `json:"number,omitempty"`
	Aired         *int       `json:"aired,omitempty"`
	Completed     *int       `json:"completed,omitempty"`
	Title         *string    `json:"title,omitempty"`
	OriginalTitle *string    `json:"original_title,omitempty"`
	IDs           *IDs       `json:"ids,omitempty"`
	Rating        *float32   `json:"rating,omitempty"`
	Votes         *int       `json:"votes,omitempty"`
	EpisodeCount  *int       `json:"episode_count,omitempty"`
	AiredEpisodes *int       `json:"aired_episodes,omitempty"`
	Overview      *string    `json:"overview,omitempty"`
	FirstAired    *Timestamp `json:"first_aired,omitempty"`
	UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
	Network       *string    `json:"network,omitempty"`
}

func (s Season) String() string {
	return Stringify(s)
}
