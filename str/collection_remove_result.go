// Package str used for structs
package str

// CollectionRemoveResult represents JSON collection remove result object
type CollectionRemoveResult struct {
	Deleted  *CollectionResultCounters `json:"deleted,omitempty"`
	NotFound *CollectionResultNotFound `json:"not_found,omitempty"`
}

func (c CollectionRemoveResult) String() string {
	return Stringify(c)
}
