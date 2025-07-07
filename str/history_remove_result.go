// Package str used for structs
package str

// HistoryRemoveResult represents JSON history remove result object
type HistoryRemoveResult struct {
	Deleted  *HistoryResultCounters `json:"deleted,omitempty"`
	NotFound *HistoryResultNotFound `json:"not_found,omitempty"`
}

func (h HistoryRemoveResult) String() string {
	return Stringify(h)
}
