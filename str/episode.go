// Package str used for structs
package str

// Episode represents JSON response for media object
type Episode struct {
	Season                *int       `json:"season,omitempty"`
	Number                *int       `json:"number,omitempty"`
	Title                 *string    `json:"title,omitempty"`
	IDs                   *IDs       `json:"ids,omitempty"`
	NumberAbs             *int       `json:"number_abs,omitempty"`
	Overview              *string    `json:"overview,omitempty"`
	Rating                *float32   `json:"rating,omitempty"`
	Votes                 *int       `json:"votes,omitempty"`
	CommentCount          *int       `json:"comment_count,omitempty"`
	FirstAired            *Timestamp `json:"first_aired,omitempty"`
	UpdatedAt             *Timestamp `json:"updated_at,omitempty"`
	AvailableTranslations *[]string  `json:"available_translations,omitempty"`
	Runtime               *int       `json:"runtime,omitempty"`
	EpisodeType           *string    `json:"episode_type,omitempty"`
}

func (s Episode) String() string {
	return Stringify(s)
}
