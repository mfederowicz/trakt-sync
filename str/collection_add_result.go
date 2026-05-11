// Package str used for structs
package str

// CollectionAddResult represents JSON collection add result object
type CollectionAddResult struct {
	Added    *CollectionResultCounters `json:"added,omitempty"`
	Updated  *CollectionResultCounters `json:"updated,omitempty"`
	Existing *CollectionResultCounters `json:"existing,omitempty"`
	NotFound *CollectionResultNotFound `json:"not_found,omitempty"`
}

func (c CollectionAddResult) String() string {
	return Stringify(c)
}
