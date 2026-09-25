// Package str used for structs
package str

// DataSyncCounts represents JSON added counts per media type
type DataSyncCounts struct {
	Movies   *int `json:"movies,omitempty"`
	Episodes *int `json:"episodes,omitempty"`
	Shows    *int `json:"shows,omitempty"`
	Seasons  *int `json:"seasons,omitempty"`
}

func (d DataSyncCounts) String() string {
	return Stringify(d)
}
