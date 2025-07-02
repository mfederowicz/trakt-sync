// Package str used for structs
package str

// ExportlistItemJSON represents JSON for list item
type ExportlistItemJSON struct {
	Title           *string    `json:"title,omitempty"`
	Trakt           *int64     `json:"trakt,omitempty"`
	Imdb            *string    `json:"imdb,omitempty"`
	Tmdb            *int       `json:"tmdb,omitempty"`
	Tvdb            *int       `json:"tvdb,omitempty"`
	WatchedAt       *Timestamp `json:"watched_at,omitempty"`
	ListedAt        *Timestamp `json:"listed_at,omitempty"`
	CollectedAt     *Timestamp `json:"collected_at,omitempty"`
	LastCollectedAt *Timestamp `json:"last_collected_at,omitempty"`
	UpdatedAt       *Timestamp `json:"updated_at,omitempty"`
	LastUpdatedAt   *Timestamp `json:"last_updated_at,omitempty"`
	Movie           *Movie     `json:"movie,omitempty"`
	Show            *Show      `json:"show,omitempty"`
	Season          *Season    `json:"season,omitempty"`
	Episode         *Episode   `json:"episode,omitempty"`
	Year            *int       `json:"year,omitempty"`
	Metadata        *Metadata  `json:"metadata,omitempty"`
}

func (i ExportlistItemJSON) String() string {
	return Stringify(i)
}

// Uptime update item time fields
func (i *ExportlistItemJSON) Uptime(options *Options, data *ExportlistItem) {
	switch options.Time {
	case "watched_at":
		i.WatchedAt = data.WatchedAt
	case "listed_at":
		i.ListedAt = data.ListedAt
	case "collected_at":
		i.CollectedAt = data.CollectedAt
	case "last_collected_at":
		i.LastCollectedAt = data.LastCollectedAt
	case "updated_at":
		i.UpdatedAt = data.UpdatedAt
	case "last_updated_at":
		i.LastUpdatedAt = data.LastUpdatedAt
	}
}

// ExportlistItem represents JSON for list item
type ExportlistItem struct {
	Title           *string    `json:"title,omitempty"`
	Year            *int       `json:"year,omitempty"`
	Rank            *int       `json:"rank,omitempty"`
	ID              *int64     `json:"id,omitempty"`
	IDs             *IDs       `json:"ids,omitempty"`
	WatchedAt       *Timestamp `json:"watched_at,omitempty"`
	ListedAt        *Timestamp `json:"listed_at,omitempty"`
	CollectedAt     *Timestamp `json:"collected_at,omitempty"`
	LastCollectedAt *Timestamp `json:"last_collected_at,omitempty"`
	UpdatedAt       *Timestamp `json:"updated_at,omitempty"`
	LastUpdatedAt   *Timestamp `json:"last_updated_at,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	Type            *string    `json:"type,omitempty"`
	Movie           *Movie     `json:"movie,omitempty"`
	Show            *Show      `json:"show,omitempty"`
	Season          *Season    `json:"season,omitempty"`
	Seasons         *[]Season  `json:"seasons,omitempty"`
	Episode         *Episode   `json:"episode,omitempty"`
	Metadata        *Metadata  `json:"metadata,omitempty"`
	MediaType       *string    `json:"media_type,omitempty"`
	Resolution      *string    `json:"resolution,omitempty"`
	Hdr             *string    `json:"hdr,omitempty"`
	Audio           *string    `json:"audio,omitempty"`
	AudioChannels   *string    `json:"audio_channels,omitempty"`
	ThreeD          *bool      `json:"3d,omitempty"`
}

func (i ExportlistItem) String() string {
	return Stringify(i)
}

// GetTime return Timestamp from item
func (i ExportlistItem) GetTime() *Timestamp {
	if i.WatchedAt != nil {
		return i.WatchedAt
	}

	if i.ListedAt != nil {
		return i.ListedAt
	}

	if i.UpdatedAt != nil {
		return i.UpdatedAt
	}

	if i.LastUpdatedAt != nil {
		return i.LastUpdatedAt
	}

	if i.CollectedAt != nil {
		return i.CollectedAt
	}

	if i.LastCollectedAt != nil {
		return i.LastCollectedAt
	}

	return nil
}

// UpdateCollectedData update meta data of object
func (i *ExportlistItem) UpdateCollectedData(item *ExportlistItem) {
	if item.Metadata != nil {
		i.MediaType = item.Metadata.MediaType
		i.Resolution = item.Metadata.Resolution
		i.Audio = item.Metadata.Audio
		i.AudioChannels = item.Metadata.AudioChannels
		i.ThreeD = item.Metadata.ThreeD
	}

	if item.Seasons != nil {
		i.Seasons = &[]Season{}
		for _, season := range *item.Seasons {
			s := Season{}
			s.Number = season.Number
			if season.Episodes != nil {
				s.Episodes = &[]Episode{}
				for _, ep := range *season.Episodes {
					e := Episode{}
					e.Number = ep.Number
					e.MediaType = ep.Metadata.MediaType
					e.Resolution = ep.Metadata.Resolution
					e.Audio = ep.Metadata.Audio
					e.AudioChannels = ep.Metadata.AudioChannels
					e.ThreeD = ep.Metadata.ThreeD

					*s.Episodes = append(*s.Episodes, e)
				}
			}
			*i.Seasons = append(*i.Seasons, s)
		}
	}
}
