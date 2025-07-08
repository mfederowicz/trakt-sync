// Package str used for structs
package str

// ItemsToRemove represents JSON items object
type ItemsToRemove struct {
	Movies   *[]Movie   `json:"movies,omitempty"`
	Shows    *[]Show    `json:"shows,omitempty"`
	Seasons  *[]Season  `json:"seasons,omitempty"`
	Episodes *[]Episode `json:"episodes,omitempty"`
	IDs      *[]int64   `json:"ids,omitempty"`
}

func (i ItemsToRemove) String() string {
	return Stringify(i)
}
