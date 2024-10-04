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
	Rank            *int       `json:"rank,omitempty"`
	ID              *int64     `json:"id,omitempty"`
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
	Episode         *Episode   `json:"episode,omitempty"`
	Metadata        *Metadata  `json:"metadata,omitempty"`
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
