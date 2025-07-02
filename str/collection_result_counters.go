// Package str used for structs
package str

// CollectionResultCounters represents JSON counters object
type CollectionResultCounters struct {
	Movies   *int `json:"movies,omitempty"`
	Episodes *int `json:"episodes,omitempty"`
	Shows    *int `json:"shows,omitempty"`
	Seasons  *int `json:"seasons,omitempty"`
}

func (c CollectionResultCounters) String() string {
	return Stringify(c)
}
