// Package str used for structs
package str

// CollectionItems represents JSON collection items object
type CollectionItems struct {
	Movies   *[]ExportlistItem `json:"movies,omitempty"`
	Shows    *[]ExportlistItem `json:"shows,omitempty"`
	Seasons  *[]ExportlistItem `json:"seasons,omitempty"`
	Episodes *[]ExportlistItem `json:"episodes,omitempty"`
}

func (c CollectionItems) String() string {
	return Stringify(c)
}
