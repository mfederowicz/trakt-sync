// Package str used for structs
package str

// CollectedItem represents JSON movies item object
type CollectedItem struct {
	CollectedAt     *Timestamp `json:"collected_at,omitempty"`
	LastCollectedAt *Timestamp `json:"last_collected_at,omitempty"`
	UpdatedAt       *Timestamp `json:"updated_at,omitempty"`
	LastUpdatedAt   *Timestamp `json:"last_updated_at,omitempty"`
	Movie           *Movie     `json:"movie,omitempty"`
	Show            *Show      `json:"show,omitempty"`
	Seasons         *[]Season  `json:"seasons,omitempty"`
	Metadata        *Metadata  `json:"metadata,omitempty"`
}

func (c CollectedItem) String() string {
	return Stringify(c)
}
