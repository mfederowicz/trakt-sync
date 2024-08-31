// Package str used for structs
package str

// HistoryListItem represents JSON list object
type HistoryListItem struct {
	ID        *int64     `json:"id,omitempty"`
	WatchedAt *Timestamp `json:"watched_at,omitempty"`
	Action    *string    `json:"action,omitempty"`
	Type      *string    `json:"type,omitempty"`
	Movie     *Movie     `json:"movie,omitempty"`
	Episode   *Episode   `json:"episode,omitempty"`
	Show      *Show      `json:"show,omitempty"`
}

func (h HistoryListItem) String() string {
	return Stringify(h)
}
