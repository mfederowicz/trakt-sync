// Package str used for structs
package str

// HistoryResultCounters represents JSON counters object
type HistoryResultCounters struct {
	Movies   *int `json:"movies,omitempty"`
	Episodes *int `json:"episodes,omitempty"`
	Shows    *int `json:"shows,omitempty"`
	Seasons  *int `json:"seasons,omitempty"`
}

func (c HistoryResultCounters) String() string {
	return Stringify(c)
}
