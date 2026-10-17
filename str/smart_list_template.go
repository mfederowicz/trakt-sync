package str

// SmartListTemplate represents JSON official smart list template; personal templates have no total and no items
type SmartListTemplate struct {
	ID        *string           `json:"id,omitempty"`
	Source    *string           `json:"source,omitempty"`
	MediaType *string           `json:"media_type,omitempty"`
	Filters   *SmartListFilters `json:"filters,omitempty"`
	Personal  *bool             `json:"personal,omitempty"`
	Total     *int              `json:"total,omitempty"`
	Items     []*UserListItem   `json:"items,omitempty"`
}

func (s SmartListTemplate) String() string {
	return Stringify(s)
}
