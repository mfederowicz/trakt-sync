// Package str used for structs
package str

// ItemsToReorder represents JSON items object
type ItemsToReorder struct {
	Rank *[]int64 `json:"rank,omitempty"`
}

func (i ItemsToReorder) String() string {
	return Stringify(i)
}
