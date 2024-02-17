package str

type ExportlistItemJson struct {
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

func (i ExportlistItemJson) String() string {
	return Stringify(i)
}

type ExportlistItem struct {
	Rank            *int       `json:"rank,omitempty"`
	Id              *int64     `json:"id,omitempty"`
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
