package str

// Recommendation represents JSON recommendation object: a movie or a show.
// The fields after IDs are set with extended info "full" (Images with "images", Colors with "colors");
// Released, AfterCredits and DuringCredits are movie only, FirstAired, LastAired, Airs, TotalRuntime,
// Network and AiredEpisodes are show only.
type Recommendation struct {
	Title                 *string      `json:"title,omitempty"`
	Year                  *int         `json:"year,omitempty"`
	IDs                   *IDs         `json:"ids,omitempty"`
	Tagline               *string      `json:"tagline,omitempty"`
	Overview              *string      `json:"overview,omitempty"`
	Released              *string      `json:"released,omitempty"`
	FirstAired            *Timestamp   `json:"first_aired,omitempty"`
	LastAired             *Timestamp   `json:"last_aired,omitempty"`
	Airs                  *Air         `json:"airs,omitempty"`
	Runtime               *int         `json:"runtime,omitempty"`
	TotalRuntime          *int         `json:"total_runtime,omitempty"`
	Certification         *string      `json:"certification,omitempty"`
	Network               *string      `json:"network,omitempty"`
	Country               *string      `json:"country,omitempty"`
	Trailer               *string      `json:"trailer,omitempty"`
	Homepage              *string      `json:"homepage,omitempty"`
	Status                *string      `json:"status,omitempty"`
	Rating                *float32     `json:"rating,omitempty"`
	Votes                 *int         `json:"votes,omitempty"`
	CommentCount          *int         `json:"comment_count,omitempty"`
	UpdatedAt             *Timestamp   `json:"updated_at,omitempty"`
	Language              *string      `json:"language,omitempty"`
	Languages             *[]string    `json:"languages,omitempty"`
	AvailableTranslations *[]string    `json:"available_translations,omitempty"`
	Genres                *[]string    `json:"genres,omitempty"`
	Subgenres             *[]string    `json:"subgenres,omitempty"`
	OriginalTitle         *string      `json:"original_title,omitempty"`
	AfterCredits          *bool        `json:"after_credits,omitempty"`
	DuringCredits         *bool        `json:"during_credits,omitempty"`
	Images                *MediaImages `json:"images,omitempty"`
	Colors                *MediaColors `json:"colors,omitempty"`
	SocialIDs             *SocialIDs   `json:"social_ids,omitempty"`
	AiredEpisodes         *int         `json:"aired_episodes,omitempty"`
	FavoritedBy           *[]UserNotes `json:"favorited_by,omitempty"`
	RecommendedBy         *[]UserNotes `json:"recommended_by,omitempty"`
}

func (r Recommendation) String() string {
	return Stringify(r)
}
