package str

type HistoryListItem struct {
	Id        *int64     `json:"id,omitempty"`
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
