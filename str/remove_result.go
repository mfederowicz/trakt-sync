// Package str used for structs
package str

// RemoveResult represents JSON history remove result object
type RemoveResult struct {
	Deleted  *ResultCounters `json:"deleted,omitempty"`
	NotFound *ResultNotFound `json:"not_found,omitempty"`
}

func (h RemoveResult) String() string {
	return Stringify(h)
}
