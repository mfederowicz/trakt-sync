// Package str used for structs
package str

// SmartListWrite represents JSON smart list create / update request object; on update every field is optional
type SmartListWrite struct {
	Name      *string           `json:"name,omitempty"`
	Source    *string           `json:"source,omitempty"`
	MediaType *string           `json:"media_type,omitempty"`
	Filters   *SmartListFilters `json:"filters,omitempty"`
	Privacy   *string           `json:"privacy,omitempty"`
}

func (s SmartListWrite) String() string {
	return Stringify(s)
}
