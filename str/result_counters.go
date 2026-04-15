// Package str used for structs
package str

// ResultCounters represents JSON counters object
type ResultCounters struct {
	Movies   *int `json:"movies,omitempty"`
	Episodes *int `json:"episodes,omitempty"`
	Shows    *int `json:"shows,omitempty"`
	Seasons  *int `json:"seasons,omitempty"`
}

func (c ResultCounters) String() string {
	return Stringify(c)
}
