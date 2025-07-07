// Package str used for structs
package str

// HistoryAddResult represents JSON history add result object
type HistoryAddResult struct {
	Added    *HistoryResultCounters `json:"added,omitempty"`
	Updated  *HistoryResultCounters `json:"updated,omitempty"`
	Existing *HistoryResultCounters `json:"existing,omitempty"`
	NotFound *HistoryResultNotFound `json:"not_found,omitempty"`
}

func (c HistoryAddResult) String() string {
	return Stringify(c)
}
