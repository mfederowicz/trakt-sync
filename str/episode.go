// Package str used for structs
package str

// Episode represents JSON response for media object
type Episode struct {
	Season                *int       `json:"season,omitempty"`
	Number                *int       `json:"number,omitempty"`
	Plays                 *int       `json:"plays,omitempty"`
	Title                 *string    `json:"title,omitempty"`
	IDs                   *IDs       `json:"ids,omitempty"`
	NumberAbs             *int       `json:"number_abs,omitempty"`
	Overview              *string    `json:"overview,omitempty"`
	Rating                *float32   `json:"rating,omitempty"`
	Votes                 *int       `json:"votes,omitempty"`
	CommentCount          *int       `json:"comment_count,omitempty"`
	FirstAired            *Timestamp `json:"first_aired,omitempty"`
	LastWatchedAt         *Timestamp `json:"last_watched_at,omitempty"`
	UpdatedAt             *Timestamp `json:"updated_at,omitempty"`
	CompletedAt           *Timestamp `json:"completed_at,omitempty"`
	CollectedAt           *Timestamp `json:"collected_at,omitempty"`
	Metadata              *Metadata  `json:"metadata,omitempty"`
	AvailableTranslations *[]string  `json:"available_translations,omitempty"`
	Runtime               *int       `json:"runtime,omitempty"`
	EpisodeType           *string    `json:"episode_type,omitempty"`
	Completed             *bool      `json:"completed,omitempty"`
	MediaType             *string    `json:"media_type,omitempty"`
	Resolution            *string    `json:"resolution,omitempty"`
	HDR                   *string    `json:"hdr,omitempty"`
	Audio                 *string    `json:"audio,omitempty"`
	AudioChannels         *string    `json:"audio_channels,omitempty"`
	ThreeD                *bool      `json:"3d,omitempty"`
}

func (e Episode) String() string {
	return Stringify(e)
}

// UpdateCollectedData update meta data of object
func (e *Episode) UpdateCollectedData(item *ExportlistItem) {
	if item.Metadata != nil {
		e.MediaType = item.Metadata.MediaType
		e.Resolution = item.Metadata.Resolution
		e.Audio = item.Metadata.Audio
		e.AudioChannels = item.Metadata.AudioChannels
		e.ThreeD = item.Metadata.ThreeD
	}

	e.CollectedAt = item.CollectedAt.UTC()
}
