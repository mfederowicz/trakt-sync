// Package str used for structs
package str

// ReorderResults represents JSON items object
type ReorderResults struct {
	Updated    *int          `json:"updated,omitempty"`
	SkippedIDs *[]int64      `json:"skipped_ids,omitempty"`
	List       *PersonalList `json:"list,omitempty"`
}

func (r ReorderResults) String() string {
	return Stringify(r)
}
