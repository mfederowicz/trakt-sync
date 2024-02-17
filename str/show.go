package str

type Show struct {
	Title                 *string    `json:"title,omitempty"`
	Year                  *int       `json:"year,omitempty"`
	Ids                   *Ids       `json:"ids,omitempty"`
	Tagline               *string    `json:"tagline,omitempty"`
	Overview              *string    `json:"overview,omitempty"`
	FirstAired            *Timestamp `json:"first_aired,omitempty"`
	Airs                  *Air       `json:"airs,omitempty"`
	Runtime               *int       `json:"runtime,omitempty"`
	Certification         *string    `json:"certification,omitempty"`
	Network               *string    `json:"network,omitempty"`
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
	AiredEpisodes         *int       `json:"aired_episodes,omitempty"`
}

func (s Show) String() string {
	return Stringify(s)
}
