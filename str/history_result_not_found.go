// Package str used for structs
package str

// HistoryResultNotFound represents JSON not found object
type HistoryResultNotFound struct {
	Movies   *[]ExportlistItem `json:"movies,omitempty"`
	Shows    *[]Show           `json:"shows,omitempty"`
	Seasons  *[]Season         `json:"seasons,omitempty"`
	Episodes *[]Episodes       `json:"episodes,omitempty"`
}

func (c HistoryResultNotFound) String() string {
	return Stringify(c)
}
