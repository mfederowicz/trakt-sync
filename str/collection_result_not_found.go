// Package str used for structs
package str

// CollectionResultNotFound represents JSON not found object
type CollectionResultNotFound struct {
	Movies   *[]ExportlistItem `json:"movies,omitempty"`
	Shows    *[]Show           `json:"shows,omitempty"`
	Seasons  *[]Season         `json:"seasons,omitempty"`
	Episodes *[]Episodes       `json:"episodes,omitempty"`
}

func (c CollectionResultNotFound) String() string {
	return Stringify(c)
}
