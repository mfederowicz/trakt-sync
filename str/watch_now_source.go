// Package str used for structs
package str

// WatchNowSource represents JSON watch now source (provider) object
type WatchNowSource struct {
	Source    *string               `json:"source,omitempty"`
	Name      *string               `json:"name,omitempty"`
	Free      *bool                 `json:"free,omitempty"`
	Cinema    *bool                 `json:"cinema,omitempty"`
	Amazon    *bool                 `json:"amazon,omitempty"`
	Color     *string               `json:"color,omitempty"`
	LinkCount *int                  `json:"link_count,omitempty"`
	Images    *WatchNowSourceImages `json:"images,omitempty"`
}

func (w WatchNowSource) String() string {
	return Stringify(w)
}
