// Package str used for structs
package str

// Movie represents JSON movie object
type Movie struct {
	Title                 *string      `json:"title,omitempty"`
	Year                  *int         `json:"year,omitempty"`
	IDs                   *IDs         `json:"ids,omitempty"`
	Tagline               *string      `json:"tagline,omitempty"`
	Overview              *string      `json:"overview,omitempty"`
	Released              *string      `json:"released,omitempty"`
	Runtime               *int         `json:"runtime,omitempty"`
	Country               *string      `json:"country,omitempty"`
	Trailer               *string      `json:"trailer,omitempty"`
	Homepage              *string      `json:"homepage,omitempty"`
	Status                *string      `json:"status,omitempty"`
	Rating                *float32     `json:"rating,omitempty"`
	Votes                 *int         `json:"votes,omitempty"`
	CommentCount          *int         `json:"comment_count,omitempty"`
	UpdatedAt             *Timestamp   `json:"updated_at,omitempty"`
	CollectedAt           *Timestamp   `json:"collected_at,omitempty"`
	Language              *string      `json:"language,omitempty"`
	Languages             *[]string    `json:"languages,omitempty"`
	AvailableTranslations *[]string    `json:"available_translations,omitempty"`
	Genres                *[]string    `json:"genres,omitempty"`
	Certification         *string      `json:"certification,omitempty"`
	User                  *UserProfile `json:"user,omitempty"`
	MediaType             *string      `json:"media_type,omitempty"`
	Resolution            *string      `json:"resolution,omitempty"`
	HDR                   *string      `json:"hdr,omitempty"`
	Audio                 *string      `json:"audio,omitempty"`
	AudioChannels         *string      `json:"audio_channels,omitempty"`
	ThreeD                *bool        `json:"3d,omitempty"`
}

// UpdateCollectedData update meta data of object
func (m *Movie) UpdateCollectedData(item *ExportlistItem) {
	if item.Metadata != nil {
		m.MediaType = item.Metadata.MediaType
		m.Resolution = item.Metadata.Resolution
		m.Audio = item.Metadata.Audio
		m.AudioChannels = item.Metadata.AudioChannels
		m.ThreeD = item.Metadata.ThreeD
	}

	m.CollectedAt = item.CollectedAt.UTC()
}

func (m Movie) String() string {
	return Stringify(m)
}
