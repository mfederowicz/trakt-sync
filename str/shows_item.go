// Package str used for structs
package str

// ShowsItem represents JSON movies item object
type ShowsItem struct {
	Watchers       *int  `json:"watchers,omitempty"`
	Revenue        *int  `json:"revenue,omitempty"`
	UserCount      *int  `json:"user_count,omitempty"`
	WatcherCount   *int  `json:"watcher_count,omitempty"`
	PlayCount      *int  `json:"play_count,omitempty"`
	CollectedCount *int  `json:"collected_count,omitempty"`
	CollectorCount *int  `json:"collector_count,omitempty"`
	ListCount      *int  `json:"list_count,omitempty"`
	Show           *Show `json:"show,omitempty"`
}

func (m ShowsItem) String() string {
	return Stringify(m)
}
