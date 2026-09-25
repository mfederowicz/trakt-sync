// Package str used for structs
package str

// SyncTraktItem represents JSON resolved Trakt movie, show or episode of a sync item; show is set for an episode
type SyncTraktItem struct {
	Type   *string        `json:"type,omitempty"`
	Title  *string        `json:"title,omitempty"`
	Year   *int           `json:"year,omitempty"`
	Season *int           `json:"season,omitempty"`
	Number *int           `json:"number,omitempty"`
	IDs    *IDs           `json:"ids,omitempty"`
	Show   *SyncTraktItem `json:"show,omitempty"`
}

func (s SyncTraktItem) String() string {
	return Stringify(s)
}
