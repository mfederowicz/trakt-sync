// Package str used for structs
package str

// AddResult represents JSON history add result object
type AddResult struct {
	Added    *ResultCounters `json:"added,omitempty"`
	Updated  *ResultCounters `json:"updated,omitempty"`
	Existing *ResultCounters `json:"existing,omitempty"`
	NotFound *ResultNotFound `json:"not_found,omitempty"`
	List     *PersonalList   `json:"list,omitempty"`
}

func (a AddResult) String() string {
	return Stringify(a)
}
