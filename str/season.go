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
	Episodes      *[]Episode `json:"episodes,omitempty"`
	AiredEpisodes *int       `json:"aired_episodes,omitempty"`
	Overview      *string    `json:"overview,omitempty"`
	FirstAired    *Timestamp `json:"first_aired,omitempty"`
	CollectedAt   *Timestamp `json:"collected_at,omitempty"`
	UpdatedAt     *Timestamp `json:"updated_at,omitempty"`
	Network       *string    `json:"network,omitempty"`
	MediaType     *string    `json:"media_type,omitempty"`
	Resolution    *string    `json:"resolution,omitempty"`
	HDR           *string    `json:"hdr,omitempty"`
	Audio         *string    `json:"audio,omitempty"`
	AudioChannels *string    `json:"audio_channels,omitempty"`
	ThreeD        *bool      `json:"3d,omitempty"`
}

func (s Season) String() string {
	return Stringify(s)
}

// UpdateCollectedData update meta data of object
func (s *Season) UpdateCollectedData(item *ExportlistItem) {
	if item.Metadata != nil {
		s.MediaType = item.Metadata.MediaType
		s.Resolution = item.Metadata.Resolution
		s.Audio = item.Metadata.Audio
		s.AudioChannels = item.Metadata.AudioChannels
		s.ThreeD = item.Metadata.ThreeD
	}
}
