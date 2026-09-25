// Package str used for structs
package str

// SmartList represents JSON smart list definition object
type SmartList struct {
	Name      *string           `json:"name,omitempty"`
	Privacy   *string           `json:"privacy,omitempty"`
	CreatedAt *Timestamp        `json:"created_at,omitempty"`
	UpdatedAt *Timestamp        `json:"updated_at,omitempty"`
	IDs       *IDs              `json:"ids,omitempty"`
	Images    *SmartListImages  `json:"images,omitempty"`
	Source    *string           `json:"source,omitempty"`
	MediaType *string           `json:"media_type,omitempty"`
	Filters   *SmartListFilters `json:"filters,omitempty"`
}

func (s SmartList) String() string {
	return Stringify(s)
}
