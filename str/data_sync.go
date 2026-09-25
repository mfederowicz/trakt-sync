// Package str used for structs
package str

// DataSync represents JSON data sync object (younify, plex or import) with its added counts
type DataSync struct {
	ID           *int64         `json:"id,omitempty"`
	CreatedAt    *Timestamp     `json:"created_at,omitempty"`
	Kind         *string        `json:"kind,omitempty"`
	Source       *string        `json:"source,omitempty"`
	Application  *string        `json:"application,omitempty"`
	Undone       *bool          `json:"undone,omitempty"`
	UndoneAt     *Timestamp     `json:"undone_at,omitempty"`
	Items        *DataSyncItems `json:"items,omitempty"`
	PausedCount  *int           `json:"paused_count,omitempty"`
	SkippedCount *int           `json:"skipped_count,omitempty"`
}

func (d DataSync) String() string {
	return Stringify(d)
}
