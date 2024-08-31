// Package str used for structs
package str

// Movie represents JSON movie object
type Movie struct {
	Title                 *string    `json:"title,omitempty"`
	Year                  *int       `json:"year,omitempty"`
	IDs                   *IDs       `json:"ids,omitempty"`
	Tagline               *string    `json:"tagline,omitempty"`
	Overview              *string    `json:"overview,omitempty"`
	Released              *string    `json:"released,omitempty"`
	Runtime               *int       `json:"runtime,omitempty"`
	Country               *string    `json:"country,omitempty"`
	Trailer               *string    `json:"trailer,omitempty"`
	Homepage              *string    `json:"homepage,omitempty"`
	Status                *string    `json:"status,omitempty"`
	Rating                *float32   `json:"rating,omitempty"`
	Votes                 *int       `json:"votes,omitempty"`
	CommentCount          *int       `json:"comment_count,omitempty"`
	UpdatedAt             *Timestamp `json:"updated_at,omitempty"`
	Language              *string    `json:"language,omitempty"`
	Languages             *[]string  `json:"languages,omitempty"`
	AvailableTranslations *[]string  `json:"available_translations,omitempty"`
	Genres                *[]string  `json:"genres,omitempty"`
	Certification         *string    `json:"certification,omitempty"`
}

func (m Movie) String() string {
	return Stringify(m)
}
