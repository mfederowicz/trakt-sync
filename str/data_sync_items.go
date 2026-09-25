// Package str used for structs
package str

// DataSyncItems represents JSON added counts per section of a data sync
type DataSyncItems struct {
	History   *DataSyncCounts `json:"history,omitempty"`
	Library   *DataSyncCounts `json:"library,omitempty"`
	Ratings   *DataSyncCounts `json:"ratings,omitempty"`
	Watchlist *DataSyncCounts `json:"watchlist,omitempty"`
}

func (d DataSyncItems) String() string {
	return Stringify(d)
}
